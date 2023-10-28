package server

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-sql-driver/mysql"
)

type server struct {
	config *Config
	db     *sql.DB
}

func CreateServer(config *Config) *server {
	return &server{
		config: config,
		db: func() *sql.DB {
			cfg := mysql.Config{
				User: config.Db_user,
				Passwd: config.Db_pass,
				Net: "tcp",
				Addr:   fmt.Sprintf("%s:%s", config.Db_host, config.Db_port),
      	  		DBName: config.Db_name,
				ParseTime: true,
			}
			conn, _ := sql.Open("mysql", cfg.FormatDSN())
			return conn
		}(),
	}
}

func (server *server) Start() {
	server.configureRouter()
	workURL := fmt.Sprintf("%s:%s", server.config.Api_host, server.config.Api_port)
	fmt.Println("server started on: ", workURL)
	http.ListenAndServe(workURL, nil)
}

func (server *server) configureRouter() {
	http.HandleFunc("/tables/", server.tables())
	http.HandleFunc("/views/", server.views())
	http.HandleFunc("/procedures/", server.procedures())
}
