package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
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
	handler.Handle("/pac/",
		NewFileHandlerX("web/sogouproxy.pac", "application/x-ns-proxy-autoconfig"))
	return handler
}

func (handler *WebHandler) serveMainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	NewFileHandler("web/index.html").ServeHTTP(w, r)
}

func (handler *WebHandler) serveAPI_close(w http.ResponseWriter, r *http.Request) {
	log.Println("GoSogouProxy is closing.")
	os.Exit(0)
}

func (handler *WebHandler) serveAPI_server(w http.ResponseWriter, r *http.Request) {
	log.Println("Send server list.")
	serverList := handler.getServerList()
	js, _ := json.Marshal(serverList)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
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

type FileHandler struct {
	name        string
	contentType string
}

func NewFileHandler(name string) *FileHandler {
	return NewFileHandlerX(name, "")
}

func NewFileHandlerX(name string, contentType string) *FileHandler {
	if contentType == "" {
		contentType = mime.TypeByExtension(filepath.Ext(name))
	}
	return &FileHandler{
		name:        name,
		contentType: contentType,
	}
}

func (handler *FileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filedata, err := ioutil.ReadFile(handler.name)
	if err != nil {
		if filedata, err = Asset(handler.name); err != nil {
			log.Printf("File %s not found.", handler.name)
			http.NotFound(w, r)
			return
		}
	}
	log.Printf("Serve file %s", handler.name)
	w.Header().Set("Content-Type", handler.contentType)
	w.Write(filedata)
}
