package main

import (
	"fmt"
	"net/http"
)

type WebHandler struct {
	*http.ServeMux
}

func NewWebHandler() *WebHandler {
	handler := &WebHandler{
		ServeMux: http.NewServeMux(),
	}
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintln(w, "It works.")
	})
	return handler
}
