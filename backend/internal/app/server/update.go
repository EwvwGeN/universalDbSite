package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func (server *server) update(w http.ResponseWriter, r *http.Request) {
	table := strings.Split(r.URL.Path[len("/tables/"):], "/")[0]
	if !server.checkExistTable(table) {
		http.Error(w, "No such table", http.StatusForbidden)
		return
	}
	if r.Method == http.MethodPut {
		reqMap := make(map[string]interface{})
		data, _ := io.ReadAll(r.Body)
		err := json.Unmarshal(data, &reqMap)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		valuesMap, ok := reqMap["updateRow"].(map[string]interface{})
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		pkMap := server.getPrimaryColumns(table)
		colMap := server.getNonPrimaryColumns(table)
		var pkCol []string
		var pkVal []interface{}
		var updateRow []string
		var row []interface{}
		for key, values := range valuesMap {
			i, err := strconv.Atoi(key[4:])
			if err != nil {
				return
			}
			if innerValue, ok := pkMap[i]; ok {
				pkCol = append(pkCol, fmt.Sprintf("`%s`=?", innerValue[0]))
				pkVal = append(pkVal, values)
				continue
			}
			updateRow = append(updateRow, fmt.Sprintf("`%s`=?", colMap[i]))
			row = append(row, values)
		}
		query := fmt.Sprintf("UPDATE `%s` SET %s WHERE %s",
			table, strings.Join(updateRow, ","), strings.Join(pkCol, " AND "))
		row = append(row, pkVal...)
		_, err = server.db.Exec(query, row...)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method == http.MethodGet {
		input := r.URL.Query().Encode()
		parsedValues, _ := url.ParseQuery(input)
		pkMap := server.getPrimaryColumns(table)
		var pkCol []string
		var pkVal []interface{}
		var buffer []string
		for _, v := range pkMap {
			inputValue, exist := parsedValues[v[0]]
			if !exist {
				http.Error(w, fmt.Sprintf("No pk column in request. Missing pk: %s", v), 400)
				return
			}
			buffer = append(buffer, "?")
			pkCol = append(pkCol, v[0])
			pkVal = append(pkVal, inputValue[0])

		}
		resp := make(map[string]interface{})
		var recRow map[int]interface{}
		query := fmt.Sprintf("SELECT * FROM `%s` WHERE (`%s`)=(%s)", table, strings.Join(pkCol, "`,`"), strings.Join(buffer, ","))
		rows, err := server.db.Query(query, pkVal...)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		columns, _ := rows.Columns()
		rows.Next()
		recRow, err = getInterfaceMap(len(columns), rows)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		rows.Close()
		resp["pkMap"] = pkMap
		resp["columns"] = columns
		resp["recRow"] = recRow

		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
