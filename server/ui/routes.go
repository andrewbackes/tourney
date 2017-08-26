package ui

import (
	"github.com/gorilla/mux"
	"net/http"
)

func Bind(r *mux.Router) {
	r.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("This is the UI."))
	}).Methods("GET")
}
