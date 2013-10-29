package main

import (
	"log"
	"net/http"
	"net/http/httptest"
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
	req, err := http.NewRequest("POST", "http://api.sportingsolutions.com/", nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	PostGitHubHookHandler(w, req)

	if w.Code != 200 {
		t.Errorf("Expecting a 200 on POST")
	}
}
