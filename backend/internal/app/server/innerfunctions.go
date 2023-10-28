package server

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func getInterfaceMap(size int, rows *sql.Rows) (map[int]interface{}, error) {
	out := make(map[int]interface{})
	rawResult := make([]interface{}, size)

	dest := make([]interface{}, size)
	for i := range rawResult {
		dest[i] = &rawResult[i]
	}
	err := rows.Scan(dest...)
	if err != nil {
		return nil, err
	}
	for j, v := range rawResult {
		switch v.(type) {
		case []uint8:
			out[j] = string(v.([]uint8))
		case time.Time:
			out[j] = v.(time.Time).Format("2006-01-02 15:04:05")
		default:
			out[j] = v
		}
	}
	return out, nil
}

func (server *server) getTables() []string {
	var tables []string
	query := "SHOW FULL TABLES WHERE TABLE_TYPE LIKE 'BASE TABLE'"
	t, err := server.db.Query(query)
	if err != nil {
		panic(err)
	}
	for t.Next() {
		var tbl string
		var tblType string
		err := t.Scan(&tbl, &tblType)
		if err != nil {
			fmt.Println(err)
			continue
		}
		tables = append(tables, tbl)
	}
	t.Close()
	return tables
}

func (server *server) getViews() []string {
	var views []string
	query := "SHOW FULL TABLES WHERE TABLE_TYPE LIKE 'VIEW'"
	t, err := server.db.Query(query)
	if err != nil {
		panic(err)
	}
	for t.Next() {
		var tbl string
		var tblType string
		err := t.Scan(&tbl, &tblType)
		if err != nil {
			fmt.Println(err)
			continue
		}
		views = append(views, tbl)
	}
	t.Close()
	return views
}

func (server *server) getProcedures() []string {
	var procedures []string
	query := `SELECT routine_name
	FROM information_schema.routines
	WHERE routine_type = 'PROCEDURE'
	AND routine_schema = database();`
	t, err := server.db.Query(query)
	if err != nil {
		panic(err)
	}
	for t.Next() {
		var prc string
		err := t.Scan(&prc)
		if err != nil {
			fmt.Println(err)
			continue
		}
		procedures = append(procedures, prc)
	}
	t.Close()
	return procedures
}

func (server *server) getProcedureParam(procedure string) map[int]string {
	out := make(map[int]string)
	query := `SELECT ORDINAL_POSITION, PARAMETER_NAME 
	FROM information_schema.parameters 
	WHERE SPECIFIC_NAME = ?
	AND PARAMETER_MODE = 'IN';`
	t, err := server.db.Query(query, procedure)
	if err != nil {
		panic(err)
	}
	for t.Next() {
		var pos int
		var param string
		err := t.Scan(&pos, &param)
		if err != nil {
			fmt.Println(err)
			continue
		}
		out[pos-1] = param
	}
	t.Close()
	return out
}

func (server *server) checkExistTable(path string) bool {
	tables := server.getTables()
	for _, value := range tables {
		if value == path {
			return true
		}
	}
	return false
}

func (server *server) checkExistView(path string) bool {
	views := server.getViews()
	for _, value := range views {
		if value == path {
			return true
		}
	}
	return false
}

func (server *server) checkExistProc(path string) bool {
	proc := server.getProcedures()
	for _, value := range proc {
		if value == path {
			return true
		}
	}
	return false
}

func (server *server) getPrimaryColumns(table string) map[int][]string {
	query := fmt.Sprintf(
		`SELECT COLUMN_NAME, ORDINAL_POSITION, IF(EXTRA LIKE '%%auto_increment%%', 1,0) as 'isAI'
		FROM information_schema.columns
		WHERE table_name='%s' AND COLUMN_KEY = 'PRI' ORDER BY ORDINAL_POSITION`,
		table)
	pkC, err := server.db.Query(query)
	if err != nil {
		println(err.Error())
		return nil
	}
	pkCol := make(map[int][]string)
	for pkC.Next() {
		var col string
		var pos int
		var isAI string
		pkC.Scan(&col, &pos, &isAI)
		pkCol[pos-1] = []string{col, isAI}
	}
	pkC.Close()
	return pkCol
}

func (server *server) getNonPrimaryColumns(table string) map[int]string {
	query := fmt.Sprintf(
		"SELECT COLUMN_NAME, ORDINAL_POSITION FROM information_schema.columns WHERE table_name='%s' AND COLUMN_KEY <> 'PRI' ORDER BY ORDINAL_POSITION",
		table)
	pkC, err := server.db.Query(query)
	if err != nil {
		println(err.Error())
		return nil
	}
	pkCol := make(map[int]string)
	for pkC.Next() {
		var col string
		var pos int
		pkC.Scan(&col, &pos)
		pkCol[pos-1] = col
	}
	pkC.Close()
	return pkCol
}

func (server *server) getTableColumns(table string) ([]string, error) {
	query := fmt.Sprintf("SELECT COLUMN_NAME FROM information_schema.columns WHERE table_name='%s' AND TABLE_SCHEMA <> 'mysql'", table)
	col, err := server.db.Query(query)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	var out []string
	for col.Next() {
		var column string
		col.Scan(&column)
		out = append(out, column)
	}
	col.Close()
	return out, nil
}

func (server *server) getTableData(table string) []map[int]interface{} {
	var outRows []map[int]interface{}
	query := fmt.Sprintf("SELECT * FROM `%s`", table)
	rows, err := server.db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	columns, _ := rows.Columns()
	for rows.Next() {
		row, _ := getInterfaceMap(len(columns), rows)
		outRows = append(outRows, row)
	}
	return outRows
}
