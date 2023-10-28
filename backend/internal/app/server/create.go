package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func (server *server) create(w http.ResponseWriter, r *http.Request) {
	table := strings.Split(r.URL.Path[len("/tables/"):], "/")[0]
	if !server.checkExistTable(table) {
		http.Error(w, "No such table", 403)
		return
	}
	if r.Method == http.MethodPost {
		reqMap := make(map[string]interface{})
		data, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(data, &reqMap)
		valuesMap, ok := reqMap["createRow"].(map[string]interface{})
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		pkMap := server.getPrimaryColumns(table)
		columns, _ := server.getTableColumns(table)
		row := make([]interface{}, len(columns))
		var buffer []string
		for key, values := range valuesMap {
			buffer = append(buffer, "?")
			i, err := strconv.Atoi(key[4:])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if innerValue, ok := pkMap[i]; ok && innerValue[1] == "1" {
				row[i] = 0
				continue
			}
			if values == "" {
				row[i] = sql.NullString{}
				continue
			}
			row[i] = values
		}
		query := fmt.Sprintf("INSERT INTO `%s` (`%s`) VALUES (%s)",
			table, strings.Join(columns, "`,`"), strings.Join(buffer, ","))
		_, err = server.db.Exec(query, row...)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		return
	}
	if r.Method == http.MethodGet {
		resp := make(map[string]interface{})
		resp["pkMap"] = server.getPrimaryColumns(table)
		resp["columns"], _ = server.getTableColumns(table)
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
