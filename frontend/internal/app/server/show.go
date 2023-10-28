package server

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"
)

func (server *server) showTable(w http.ResponseWriter, r *http.Request) {
	table := strings.Split(r.URL.Path[len("/tables/"):], "/")[0]
	getURL := fmt.Sprintf("http://%s:%s/tables/%s", server.config.Api_host, server.config.Api_port, table)
	receivedMap, err := server.getRequest(http.MethodGet, getURL, nil)
	if err != nil {
		fmt.Fprintf(w, "error: %s", err)
		return
	}
	navMap := server.getAllFromDB()
	type Navigation struct {
		Tables     interface{}
		Views      interface{}
		Procedures interface{}
	}
	type Table_data struct {
		Columns interface{}
		PkCol   interface{}
		Data    interface{}
	}

	type Table_page struct {
		Page_title string
		Table      Table_data
		Navigation Navigation
	}

	data := Table_page{
		Page_title: table,
		Table: Table_data{
			Columns: receivedMap["columns"],
			PkCol:   receivedMap["pkMap"],
			Data:    receivedMap["rows"],
		},
		Navigation: Navigation{
			Tables:     navMap["tables"],
			Views:      navMap["views"],
			Procedures: navMap["procedures"],
		},
	}
	files := []string{
		"./web/html/table_page.html",
		"./web/html/layout.html",
		"./web/html/footer.html",
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

func (server *server) showView(w http.ResponseWriter, r *http.Request) {
	view := strings.Split(r.URL.Path[len("/views/"):], "/")[0]

	receivedMap := make(map[string]interface{})
	getURL := fmt.Sprintf("http://%s:%s/views/%s", server.config.Api_host, server.config.Api_port, view)
	receivedMap, err := server.getRequest(http.MethodGet, getURL, nil)
	if err != nil {
		fmt.Fprintf(w, "error: %s", err)
		return
	}
	navMap := server.getAllFromDB()
	type Navigation struct {
		Tables     interface{}
		Views      interface{}
		Procedures interface{}
	}
	type View_data struct {
		Columns interface{}
		Data    interface{}
	}

	type View_page struct {
		Page_title string
		View       View_data
		Navigation Navigation
	}

	data := View_page{
		Page_title: view,
		View: View_data{
			Columns: receivedMap["columns"],
			Data:    receivedMap["rows"],
		},
		Navigation: Navigation{
			Tables:     navMap["tables"],
			Views:      navMap["views"],
			Procedures: navMap["procedures"],
		},
	}
	files := []string{
		"./web/html/view_page.html",
		"./web/html/layout.html",
		"./web/html/footer.html",
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

func (server *server) showProcedure(w http.ResponseWriter, r *http.Request) {
	procedure := strings.Split(r.URL.Path[len("/procedures/"):], "/")[0]

	receivedMap := make(map[string]interface{})
	getURL := fmt.Sprintf("http://%s:%s/procedures/%s", server.config.Api_host, server.config.Api_port, procedure)
	receivedMap, err := server.getRequest(http.MethodGet, getURL, nil)
	if err != nil {
		fmt.Fprintf(w, "error: %s", err)
		return
	}
	navMap := server.getAllFromDB()
	type Navigation struct {
		Tables     interface{}
		Views      interface{}
		Procedures interface{}
	}
	type Proc_data struct {
		Param interface{}
	}

	type Proc_page struct {
		Page_title string
		Procedure  Proc_data
		Navigation Navigation
	}

	data := Proc_page{
		Page_title: procedure,
		Procedure: Proc_data{
			Param: receivedMap["procParam"],
		},
		Navigation: Navigation{
			Tables:     navMap["tables"],
			Views:      navMap["views"],
			Procedures: navMap["procedures"],
		},
	}
	files := []string{
		"./web/html/proc_page.html",
		"./web/html/layout.html",
		"./web/html/footer.html",
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

func (server *server) showProcedureAns(w http.ResponseWriter, r *http.Request, input map[string]interface{}) {
	navMap := server.getAllFromDB()
	type Navigation struct {
		Tables     interface{}
		Views      interface{}
		Procedures interface{}
	}
	type Answer_data struct {
		Columns interface{}
		Data    interface{}
	}

	type Answer_page struct {
		Page_title interface{}
		Answer     Answer_data
		Navigation Navigation
	}

	data := Answer_page{
		Page_title: input["procedure"],
		Answer: Answer_data{
			Columns: input["columns"],
			Data:    input["rows"],
		},
		Navigation: Navigation{
			Tables:     navMap["tables"],
			Views:      navMap["views"],
			Procedures: navMap["procedures"],
		},
	}
	files := []string{
		"./web/html/answer_page.html",
		"./web/html/layout.html",
		"./web/html/footer.html",
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
