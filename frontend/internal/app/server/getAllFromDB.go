package server

import (
	"fmt"
	"log"
	"net/http"
)

func (server *server) getAllFromDB() map[string]interface{} {
	out := make(map[string]interface{})
	getURL := fmt.Sprintf("http://%s:%s/tables", server.config.Api_host, server.config.Api_port)
	getMap, err := server.getRequest(http.MethodGet, getURL, nil)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	out["tables"] = getMap["tables"]
	getURL = fmt.Sprintf("http://%s:%s/views", server.config.Api_host, server.config.Api_port)
	getMap, err = server.getRequest(http.MethodGet, getURL, nil)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	out["views"] = getMap["views"]
	getURL = fmt.Sprintf("http://%s:%s/procedures", server.config.Api_host, server.config.Api_port)
	getMap, err = server.getRequest(http.MethodGet, getURL, nil)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	out["procedures"] = getMap["procedures"]
	return out
}
