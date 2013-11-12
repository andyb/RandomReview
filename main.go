package main

import (
	"randomreview/review"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
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
	router.HandleFunc("/h", HealthCheckHandler)
	http.Handle("/", router)
	http.ListenAndServe(ip, router)
}

func HealthCheckHandler(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "I'm good, thanks for asking")
}

func PostGitHubHookHandler(rw http.ResponseWriter, req *http.Request) {
	log.Println("Recieved POST request from GitHub. Processing....")
	if req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			logErrorAndReturnHttpError(err, rw, 400)
		}
		
		reviewers := loadReviewers()

			
		values, err := url.ParseQuery(string(body))

		var jsonBody interface{}
		json.Unmarshal([]byte(values["payload"][0]), &jsonBody)

		revreq, err := review.GenerateReviewRequest(jsonBody, reviewers)
		review.SendReviewRequestEmail(revreq)

		if err != nil {
			logErrorAndReturnHttpError(err, rw, 400)
		}

	} else {
		log.Println("Received request that wasn't a POST so ignored")
		http.Error(rw, "File not found", 404)
	}

}

func loadReviewers() []review.Reviewer {
	file, err := ioutil.ReadFile("reviewers.json")
	if err != nil {
		log.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	
	var reviewers []review.Reviewer
	err = json.Unmarshal(file, &reviewers)
	
	if err != nil {
		log.Printf("Error parsing reviewers file: %v\n", err)
		os.Exit(1)
	}
	
	return reviewers
}

func logErrorAndReturnHttpError(err error, w http.ResponseWriter, statusCode int) {
	log.Println(err)
	http.Error(w, err.Error(), statusCode)
}
