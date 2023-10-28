package server

import (
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type outFunc func(http.ResponseWriter, *http.Request)

func (server *server) tables() outFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(r.URL.Path[len("/tables/"):], "/")
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			if len(path) == 0 {
				server.showTable(w, r, "")
				return
			}
			if len(path) == 1 {
				table := path[0]
				server.showTable(w, r, table)
				return
			}
			if len(path) == 2 && path[1] == "create" {
				server.create(w, r)
			}
			if len(path) == 2 && path[1] == "update" {
				server.update(w, r)
			}
		case http.MethodPost:
			if len(path) == 2 && path[1] == "create" {
				server.create(w, r)
			}
		case http.MethodPut:
			server.update(w, r)
		case http.MethodDelete:
			server.delete(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func (server *server) views() outFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(r.URL.Path[len("/views/"):], "/")
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			if len(path) == 0 {
				server.showView(w, r, "")
				return
			}
			if len(path) == 1 {
				view := path[0]
				server.showView(w, r, view)
				return
			}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func (server *server) procedures() outFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(r.URL.Path[len("/procedures/"):], "/")
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			if len(path) == 0 {
				server.showProcedures(w, r, "")
				return
			}
			if len(path) == 1 {
				procedure := path[0]
				server.showProcedures(w, r, procedure)
				return
			}
		case http.MethodPost:
			if len(path) == 1 {
				server.callProc(w, r)
			}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}