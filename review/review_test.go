package review

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

const (
	expected_review_link     = "https://github.com/sportingsolutions/SS.GTP/compare/75ee12a20fb0...ecea7cbada12"
	expected_from            = "tom-scott"
	expected_reviewer_name   = "Reviewer One"
	expected_reviewer_email  = "abritcliffe@sportingindex.com"
	expected_reviewer_github = "andyb"
)

func Test_Parsing_GitHub_WebHook(t *testing.T) {
	payload := loadJSONPayload("payload.json")
	req, _ := GenerateReviewRequest(payload, loadTestReviewers())
	if req.From != expected_from {
		t.Errorf("Incorrect from property: %v expecting %v", req.From, expected_from)
	}

	if req.To.Name != expected_reviewer_name {
		t.Errorf("Incorrect to property: %v expecting %v", req.To.Name, expected_reviewer_name)
	}

	if req.To.Email != expected_reviewer_email {
		t.Errorf("Incorrect to property: %v expecting %v", req.To.Email, expected_reviewer_email)
	}

	if req.Review_link != expected_review_link {
		t.Errorf("Incorrect property: %v expecting %v", req.Review_link, expected_review_link)
	}

	if req.To.Githubusername != expected_reviewer_github {
		t.Errorf("Incorrect property: %v expecting %v", req.To.Githubusername, expected_reviewer_github)
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

func Test_selected_self_review(t *testing.T) {
	payload := loadJSONPayload("payload.json")
	req, _ := GenerateReviewRequest(payload, loadTestReviewersSelfTest())
	log.Println(req.From)
	
	if req.From != expected_from {
		t.Errorf("Incorrect from property: %v expecting %v", req.From, expected_from)
	}
	
	if req.To.Githubusername != expected_from {
		t.Errorf("Incorrect property: %v expecting %v", req.To.Githubusername, expected_reviewer_github)
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
	expected_reviewer := Reviewer{Name: expected_reviewer_name, Email: expected_reviewer_email, Githubusername: "andyb"}
	reviewers = make([]Reviewer, 3)
	reviewers[0] = expected_reviewer
	reviewers[1] = expected_reviewer
	reviewers[2] = expected_reviewer
	return
}

func loadTestReviewersSelfTest() (reviewers []Reviewer) {
	reviewers = make([]Reviewer, 2)
	reviewers[0] = Reviewer{Name: "Entry 0", Email: expected_reviewer_email, Githubusername: expected_from}
	reviewers[1] = Reviewer{Name: "Entry 1", Email: expected_reviewer_email, Githubusername: expected_from}
	return
}
