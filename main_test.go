package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGitHubHookHandlerNonPost(t *testing.T) {
	req, err := http.NewRequest("GET", "http://api.sportingsolutions.com/", nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	PostGitHubHookHandler(w, req)
	if w.Code != 404 {
		t.Errorf("Expecting a 404 on non POST")
	}
	log.Printf("%d - %s", w.Code, w.Body.String())
}

func TestPostGitHubHookHandlerPost(t *testing.T) {
	bodyBytes := loadRAWPayload("review/payload.raw")
	req, err := http.NewRequest("POST", "http://api.sportingsolutions.com/", bytes.NewReader(bodyBytes))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	PostGitHubHookHandler(w, req)

	if w.Code != 200 {
		t.Errorf("Expecting a 200 on POST")
	}
}

func loadRAWPayload(fileName string) (file []byte) {
	file, e := ioutil.ReadFile(fileName)
	if e != nil {
		log.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	return
}
