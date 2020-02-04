package server

import (
	"github.com/gorilla/mux"
)

// CreateRouter creates the routes for the server
func CreateRouter() {
	router := mux.NewRouter()
	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/welcome.html", welcomeHandler)
	router.HandleFunc("/headers", headersHandler)
}
