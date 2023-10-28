package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"text/template"

	"github.com/EwvwGeN/universalDbSite/frontend/internal/help"
)

func (server *server) update(w http.ResponseWriter, r *http.Request) {
	table := strings.Split(r.URL.Path[len("/tables/"):], "/")[0]
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", 500)
			return
		}
		preparedMap := make(map[string]interface{})
		updateRow := make(map[string]interface{})
		for key, value := range r.PostForm {
			updateRow[key] = value[0]
		}
		preparedMap["updateRow"] = updateRow
		bytesRepresentation, err := json.Marshal(preparedMap)
		if err != nil {
			log.Fatalln(err)
		}
		updateURL := fmt.Sprintf("http://%s:%s/tables/%s/update", server.config.Api_host, server.config.Api_port, table)
		req, err := http.NewRequest(http.MethodPut, updateURL, bytes.NewBuffer(bytesRepresentation))
		if err != nil {
			log.Fatal(err)
		}
		resp, err := server.client.Do(req)
		if err != nil {
			panic(err)
		}
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		resp.Body.Close()
		if resp.StatusCode == http.StatusForbidden {
			http.Redirect(w, r, "/home", 301)
			return
		}
		if resp.StatusCode == http.StatusBadRequest {
			http.Redirect(w, r, fmt.Sprintf("/tables/%s/update?%s", table, r.URL.Query().Encode()), 301)
			return
		}
		if resp.StatusCode == http.StatusOK {
			http.Redirect(w, r, fmt.Sprintf("/tables/%s", table), 301)
			return
		}
	}
	if r.Method == http.MethodGet {
		files := []string{
			"./web/html/update_page.html",
			"./web/html/layout.html",
			"./web/html/footer.html",
		}
		input := r.URL.Query().Encode()
		v, _ := url.ParseQuery(input)
		var preparedLink string
		for i, value := range v {
			preparedLink += i + "=" + value[0] + "&"
		}
		updateURL := fmt.Sprintf("http://%s:%s/tables/%s/update?%s", server.config.Api_host, server.config.Api_port, table, preparedLink)
		receivedMap, err := server.getRequest(http.MethodGet, updateURL, nil)
		if err != nil && err.Error() == "400" {
			http.Redirect(w, r, fmt.Sprintf("/%s", table), 301)
			return
		}
		if err != nil {
			http.Redirect(w, r, "/home", 301)
			return
		}
		bufferData := help.MapToIntKeys(receivedMap["recRow"].(map[string]interface{}))
		bufferPkMap := help.MapToIntKeys(receivedMap["pkMap"].(map[string]interface{}))
		navMap := server.getAllFromDB()
		type Navigation struct {
			Tables     interface{}
			Views      interface{}
			Procedures interface{}
		}
		type Record_data struct {
			PkCol   interface{}
			Columns interface{}
			Data    interface{}
		}

		type Record_page struct {
			Page_title string
			Record     Record_data
			Navigation Navigation
		}

		data := Record_page{
			Page_title: table,
			Record: Record_data{
				PkCol:   bufferPkMap,
				Columns: receivedMap["columns"],
				Data:    bufferData,
			},
			Navigation: Navigation{
				Tables:     navMap["tables"],
				Views:      navMap["views"],
				Procedures: navMap["procedures"],
			},
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}

		err = ts.Execute(w, data)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}
	}
	return
}
