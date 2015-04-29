package main

import (
	"fmt"
	"io"
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
	handler.HandleFunc("/pac/", handler.servePac)
	return handler
}

func (handler *WebHandler) serveMainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintln(w, "GoSogouProxy")
	serverList := handler.getServerList()
	fmt.Fprintf(w, "%d servers available:\n", len(serverList))
	for _, server := range serverList {
		fmt.Fprintln(w, server)
	}
}

func (handler *WebHandler) serveAPI_close(w http.ResponseWriter, r *http.Request) {
	log.Println("GoSogouProxy is closing.")
	os.Exit(0)
}

func (handler *WebHandler) servePac(w http.ResponseWriter, r *http.Request) {
	pac, err := os.Open("sogouproxy.pac")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	io.Copy(w, pac)
	pac.Close()
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
