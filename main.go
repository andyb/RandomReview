package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const (
	ip = "0.0.0.0:8080"
)

func main() {
	startWeb()
}

func startWeb() {
	log.Printf("Starting web server @ %v...", ip)
	router := mux.NewRouter()
	router.HandleFunc("/", PostGitHubHook).Methods("POST")
	http.Handle("/", router)
	http.ListenAndServe(ip, router)
}

func PostGitHubHook(rw http.ResponseWriter, req *http.Request) {
	log.Println("Revived POST request from GitHub. Processing....")

}
