package server

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func (server *server) delete(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path[len("/tables/"):], "/")[0]
	if !server.checkExistTable(path) {
		http.Error(w, "No such table", 403)
		return
	}
	if r.Method == http.MethodDelete {
		input := r.URL.Query().Encode()
		parsedValues, _ := url.ParseQuery(input)
		pkMap := server.getPrimaryColumns(path)
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
		query := fmt.Sprintf("DELETE FROM `%s` WHERE (`%s`)=(%s)", path, strings.Join(pkCol, "`,`"), strings.Join(buffer, ","))
		_, err := server.db.Exec(query, pkVal...)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
	return
}
