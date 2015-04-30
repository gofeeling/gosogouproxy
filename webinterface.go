package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type WebHandler struct {
	*http.ServeMux
	ProxyType
	getlistReqChan chan chan []int
}

func NewWebHandler(proxyType ProxyType) *WebHandler {
	handler := &WebHandler{
		ServeMux:       http.NewServeMux(),
		ProxyType:      proxyType,
		getlistReqChan: make(chan chan []int),
	}
	handler.HandleFunc("/", handler.serveMainPage)
	handler.HandleFunc("/api/close/", handler.serveAPI_close)
	handler.HandleFunc("/api/server/", handler.serveAPI_server)
	handler.Handle("/pac/", fileHandler("sogouproxy.pac"))
	return handler
}

func (handler *WebHandler) serveMainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fileHandler("index.html").ServeHTTP(w, r)
}

func (handler *WebHandler) serveAPI_close(w http.ResponseWriter, r *http.Request) {
	log.Println("GoSogouProxy is closing.")
	os.Exit(0)
}

func (handler *WebHandler) serveAPI_server(w http.ResponseWriter, r *http.Request) {
	log.Println("Send server list.")
	serverList := handler.getServerList()
	response, _ := json.Marshal(serverList)
	w.Write(response)
}

func (handler *WebHandler) getServerList() []string {
	listChan := make(chan []int)
	handler.getlistReqChan <- listChan
	indices := <-listChan
	var serverList []string
	for _, proxyNum := range indices {
		serverList = append(serverList, fmt.Sprintf(handler.hostTemplate, proxyNum))
	}
	return serverList
}

type fileHandler string

func (f fileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filename := string(f)
	log.Printf("Serve file %s", filename)
	filedata, err := ioutil.ReadFile(filename)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Write(filedata)
}
