package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

func (server *server) create(w http.ResponseWriter, r *http.Request) {
	table := strings.Split(r.URL.Path[len("/tables/"):], "/")[0]
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", 500)
			return
		}
		preparedMap := make(map[string]interface{})
		for key, values := range r.PostForm {
			preparedMap[key] = values[0]
		}
		outMap := make(map[string]interface{})
		outMap["createRow"] = preparedMap
		bytesRepresentation, err := json.Marshal(outMap)
		if err != nil {
			log.Fatalln(err)
		}
		createURL := fmt.Sprintf("http://%s:%s/tables/%s/create", server.config.Api_host, server.config.Api_port, table)
		req, err := http.NewRequest(http.MethodPost, createURL, bytes.NewBuffer(bytesRepresentation))
		if err != nil {
			log.Fatal(err.Error())
		}
		resp, err := server.client.Do(req)
		if err != nil {
			panic(err)
		}
		resp.Body.Close()
		if resp.StatusCode == http.StatusCreated {
			http.Redirect(w, r, fmt.Sprintf("/tables/%s", table), 301)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/tables/%s/create", table), 301)
		return
	}
	files := []string{
		"./web/html/create_page.html",
		"./web/html/layout.html",
		"./web/html/footer.html",
	}
	createUrl := fmt.Sprintf("http://%s:%s/tables/%s/create", server.config.Api_host, server.config.Api_port, table)
	req, err := http.NewRequest(http.MethodGet, createUrl, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	resp, err := server.client.Do(req)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	receivedMap := make(map[string]interface{})
	err = json.Unmarshal(body, &receivedMap)
	if err != nil {
		log.Println(err.Error())
		return
	}
	navMap := server.getAllFromDB()
	type Navigation struct {
		Tables     interface{}
		Views      interface{}
		Procedures interface{}
	}
	type Record_data struct {
		PkCol   interface{}
		Columns interface{}
	}

	type Record_page struct {
		Page_title string
		Record     Record_data
		Navigation Navigation
	}
	bufferPkMap := make(map[int]interface{})
	for i, v := range receivedMap["pkMap"].(map[string]interface{}) {
		idx, _ := strconv.Atoi(i)
		bufferPkMap[idx] = v
	}
	data := Record_page{
		Page_title: table,
		Record: Record_data{
			PkCol:   bufferPkMap,
			Columns: receivedMap["columns"],
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
