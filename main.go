package main

import (
	"./review"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
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
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			logErrorAndReturnHttpError(err, rw, 400)
		}

		//need to tidy this up
		reviewers := make([]review.Reviewer, 1)
		reviewers[0] = review.Reviewer{Name: "Andy Britcliffe", Email: "abritcliffe@sportingsolutions.com"}

		var jsonBody interface{}
		json.Unmarshal(body, &jsonBody)
		revreq, err := review.GenerateReviewRequest(jsonBody, reviewers)
		review.SendReviewRequestEmail(revreq)

		if err != nil {
			logErrorAndReturnHttpError(err, rw, 400)
		}

		log.Printf("%s", revreq)

	} else {
		log.Println("Received request that wasn't a POST so ignored")
		http.Error(rw, "File not found", 404)
	}

}

func logErrorAndReturnHttpError(err error, w http.ResponseWriter, statusCode int) {
	log.Println(err)
	http.Error(w, err.Error(), statusCode)
}
