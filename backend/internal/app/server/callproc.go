package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func (server *server) callProc(w http.ResponseWriter, r *http.Request) {
	procedure := strings.Split(r.URL.Path[len("/procedures/"):], "/")[0]
	if !server.checkExistProc(procedure) {
		http.Error(w, "No such table", 403)
		return
	}
	if r.Method == http.MethodPost {
		reqMap := make(map[string]interface{})
		data, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(data, &reqMap)
		valuesMap, ok := reqMap["procParam"].(map[string]interface{})
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		params := server.getProcedureParam(procedure)
		inputParams := make([]interface{}, len(params))
		var callRow []string
		for key, values := range valuesMap {
			i, err := strconv.Atoi(key[4:])
			if err != nil {
				return
			}
			callRow = append(callRow, "?")
			inputParams[i] = values
		}
		query := fmt.Sprintf("CALL `%s` (%s)",
			procedure, strings.Join(callRow, ","))
		rows, err := server.db.Query(query, inputParams...)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		columns, _ := rows.Columns()
		var outRows []map[int]interface{}
		for rows.Next() {
			row, _ := getInterfaceMap(len(columns), rows)
			outRows = append(outRows, row)
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		rows.Close()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		resp := make(map[string]interface{})
		if len(columns) != 0 {
			resp["procedure"] = procedure
			resp["columns"] = columns
			resp["rows"] = outRows
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			}
			w.Write(jsonResp)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
	return
}
