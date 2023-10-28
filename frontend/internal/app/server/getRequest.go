package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func (server *server) getRequest(method string, url string, inputBody io.Reader) (map[string]interface{}, error) {
	receivedMap := make(map[string]interface{})
	req, err := http.NewRequest(method, url, inputBody)
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
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d", resp.StatusCode)
	}
	err = json.Unmarshal(body, &receivedMap)
	if err != nil {
		log.Fatal(err.Error())
	}
	return receivedMap, nil
}
