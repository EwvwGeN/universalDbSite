package server

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func (server *server) delete(w http.ResponseWriter, r *http.Request) {
	table := strings.Split(r.URL.Path[len("/tables/"):], "/")[0]
	input := r.URL.Query().Encode()
	parsedValues, _ := url.ParseQuery(input)
	prepared := "?"
	for i, v := range parsedValues {
		prepared += i + "=" + v[0] + "&"
	}
	
	deleteUrl := fmt.Sprintf("http://%s:%s/tables/%s/delete%s", server.config.Api_host, server.config.Api_port, table, prepared)
	req, err := http.NewRequest(http.MethodDelete, deleteUrl, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	resp, err := server.client.Do(req)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	http.Redirect(w, r, fmt.Sprintf("/tables/%s", table), 301)
}
