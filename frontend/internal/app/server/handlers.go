package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"
)

type outFunc func(http.ResponseWriter, *http.Request)

func (server *server) tables() outFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(r.URL.Path[len("/tables/"):], "/")
		if len(path) == 1 {
			server.showTable(w, r)
			return
		}
		if len(path) == 2 && path[1] == "delete" {
			server.delete(w, r)
			return
		}
		if len(path) == 2 && path[1] == "update" {
			server.update(w, r)
			return
		}
		if len(path) == 2 && path[1] == "create" {
			server.create(w, r)
			return
		}
		http.Redirect(w, r, "/home", 301)
		return
	}
}

func (server *server) views() outFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(r.URL.Path[len("/views/"):], "/")
		if len(path) == 1 {
			server.showView(w, r)
			return
		}
	}
}

func (server *server) procedures() outFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(r.URL.Path[len("/procedures/"):], "/")
		procedure := path[0]
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
			preparedMap["procParam"] = updateRow
			bytesRepresentation, err := json.Marshal(preparedMap)
			if err != nil {
				log.Fatalln(err)
			}
			procURL := fmt.Sprintf("http://%s:%s/procedures/%s", server.config.Api_host, server.config.Api_port, procedure)
			req, err := http.NewRequest(http.MethodPost, procURL, bytes.NewBuffer(bytesRepresentation))
			if err != nil {
				log.Fatal(err)
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
			if resp.StatusCode == http.StatusForbidden {
				http.Redirect(w, r, "/home", 301)
				return
			}
			if resp.StatusCode == http.StatusBadRequest {
				http.Redirect(w, r, fmt.Sprintf("/procedures/%s", procedure), 301)
				return
			}
			receivedMap := make(map[string]interface{})
			if len(body) != 0 {
				err = json.Unmarshal(body, &receivedMap)
				if err != nil {
					log.Fatal(err.Error())
				}
				server.showProcedureAns(w, r, receivedMap)
				return
			}
			if resp.StatusCode == http.StatusOK {
				http.Redirect(w, r, "/home", 301)
				return
			}
		}
		if r.Method == http.MethodGet {
			if len(path) == 1 {
				server.showProcedure(w, r)
				return
			}
		}
	}
}

func (server *server) home() outFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		files := []string{
			"./web/html/home.html",
			"./web/html/layout.html",
			"./web/html/footer.html",
		}
		receivedMap := server.getAllFromDB()
		type Navigation struct {
			Tables     interface{}
			Views      interface{}
			Procedures interface{}
		}
		nav := Navigation{
			Tables:     receivedMap["tables"],
			Views:      receivedMap["views"],
			Procedures: receivedMap["procedures"],
		}
		type Page struct {
			Page_title string
			Navigation Navigation
		}

		data := Page{
			Page_title: "home",
			Navigation: nav,
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
}
