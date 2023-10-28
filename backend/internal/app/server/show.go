package server

import (
	"encoding/json"
	"log"
	"net/http"
)

func (server *server) showTable(w http.ResponseWriter, r *http.Request, table string) {
	resp := make(map[string]interface{})
	if table == "" {
		resp["tables"] = server.getTables()
	}
	if !server.checkExistTable(table) && table != "" {
		http.Error(w, "No such table", http.StatusForbidden)
		return
	}
	if table != "" {
		resp["pkMap"] = server.getPrimaryColumns(table)
		resp["columns"], _ = server.getTableColumns(table)
		resp["rows"] = server.getTableData(table)
	}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func (server *server) showView(w http.ResponseWriter, r *http.Request, view string) {
	resp := make(map[string]interface{})
	if view == "" {
		resp["views"] = server.getViews()
	}
	if !server.checkExistView(view) && view != "" {
		http.Error(w, "No such view", http.StatusForbidden)
		return
	}
	if view != "" {
		resp["columns"], _ = server.getTableColumns(view)
		resp["rows"] = server.getTableData(view)
	}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func (server *server) showProcedures(w http.ResponseWriter, r *http.Request, procedure string) {
	resp := make(map[string]interface{})
	if procedure == "" {
		resp["procedures"] = server.getProcedures()
	}
	if procedure != "" && !server.checkExistProc(procedure) {
		http.Error(w, "No such procedure", http.StatusForbidden)
		return
	}
	if procedure != "" {
		resp["procParam"] = server.getProcedureParam(procedure)
	}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}
