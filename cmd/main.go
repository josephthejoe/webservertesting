package main

import (
    "fmt"
    "net/http"
    
    "github.com/jospehthejoe/webservertesting/internal/handlers/bird_handler"
    "github.com/gorilla/mux"

)

func newRouter() *mux.Router {
    r := mux.NewRouter()
    r.HandleFunc("/hello", handler).Methods("GET")

    staticFileDirectory := http.Dir("./web/")
	staticFileHandler := http.StripPrefix("/web/", http.FileServer(staticFileDirectory))
	r.PathPrefix("/web/").Handler(staticFileHandler).Methods("GET")

    r.HandleFunc("/bird", getBirdHandler).Methods("GET")
	r.HandleFunc("/bird", createBirdHandler).Methods("POST")

    return r
}
func main() {
    r := newRouter()
    http.ListenAndServe(":8080", r)
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "hello, world")
}


