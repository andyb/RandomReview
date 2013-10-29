package review

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

const (
	expected_review_link    = "https://github.com/sportingsolutions/SS.GTP/compare/75ee12a20fb0...ecea7cbada12"
	expected_from           = "tom-scott"
	expected_reviewer_name  = "Reviewer One"
	expected_reviewer_email = "abritcliffe@sportingindex.com"
)

func Test_Parsing_GitHub_WebHook(t *testing.T) {
	payload := loadJSONPayload("payload.json")
	req, _ := GenerateReviewRequest(payload, loadTestReviewers())
	if req.from != expected_from {
		t.Errorf("Incorrect from property: %v expecting %v", req.from, expected_from)
	}

	if req.to.name != expected_reviewer_name {
		t.Errorf("Incorrect to property: %v expecting %v", req.to, expected_reviewer_name)
	}

	if req.to.email != expected_reviewer_email {
		t.Errorf("Incorrect to property: %v expecting %v", req.to, expected_reviewer_email)
	}

	if req.review_link != expected_review_link {
		t.Errorf("Incorrect review_link property: %v expecting %v", req.review_link, expected_review_link)
	}

	SendReviewRequestEmail(req)
}

func Test_Failed_Parsing_GitHub_WebHook(t *testing.T) {
	payload := loadJSONPayload("badpayload.json")
	_, err := GenerateReviewRequest(payload, loadTestReviewers())
	if err == nil {
		t.Error("Expecting error to be returned")
	}
}

func loadJSONPayload(fileName string) (payload interface{}) {
	file, e := ioutil.ReadFile(fileName)
	if e != nil {
		log.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	json.Unmarshal(file, &payload)
	return
}

func loadTestReviewers() (reviewers []Reviewer) {
	expected_reviewer := reviewer{name: expected_reviewer_name, email: expected_reviewer_email}
	reviewers = make([]reviewer, 3)
	reviewers[0] = expected_reviewer
	reviewers[1] = expected_reviewer
	reviewers[2] = expected_reviewer
	return
}
