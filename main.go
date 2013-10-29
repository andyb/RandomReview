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
	router.HandleFunc("/", PostGitHubHookHandler).Methods("POST")
	http.Handle("/", router)
	http.ListenAndServe(ip, router)
}

func PostGitHubHookHandler(rw http.ResponseWriter, req *http.Request) {
	log.Println("Recieved POST request from GitHub. Processing....")
	if req.Method == "POST" {

	} else {
		log.Println("Received request that wasn't a POST so ignored")
		http.Error(rw, "File not found", 404)
	}

}

func logErrorAndReturnHttpError(err error, w http.ResponseWriter, statusCode int) {
	log.Println(err)
	http.Error(w, err.Error(), statusCode)
}
