package server

import (
	"fmt"
	"net/http"
)

type server struct {
	config *Config
	client *http.Client
}

func NewServer(cnfg *Config) *server {
	return &server{
		config: cnfg,
		client: &http.Client{},
	}
}

func (server *server) Start() {
	server.configureRouter()
	workURL := fmt.Sprintf("%s:%s", server.config.Site_host, server.config.Site_port)
	fmt.Println("site started on: ", workURL)
	http.ListenAndServe(workURL, nil)
}

func (server *server) configureRouter() {
	http.HandleFunc("/tables/", server.tables())
	http.HandleFunc("/views/", server.views())
	http.HandleFunc("/procedures/", server.procedures())
	http.HandleFunc("/home", server.home())
	fileServer := http.FileServer(http.Dir("web/static/"))
	http.Handle("/static/", http.StripPrefix("/static", fileServer))
}
