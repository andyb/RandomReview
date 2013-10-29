package review

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

const (
	expected_review_link = "https://github.com/sportingsolutions/SS.GTP/compare/75ee12a20fb0...ecea7cbada12"
	expected_reviewer    = "Reviewer One"
	expected_from        = "tom-scott"
)

func Test_Parsing_GitHub_WebHook(t *testing.T) {
	payload := loadJSONPayload("payload.json")
	req, _ := GenerateReviewRequest(payload, loadTestReviewers())
	if req.from != expected_from {
		t.Errorf("Incorrect from property: %v expecting %v", req.from, expected_from)
	}

	if req.to != expected_reviewer {
		t.Errorf("Incorrect to property: %v expecting %v", req.to, expected_reviewer)
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

func loadTestReviewers() (reviewers []string) {
	reviewers = make([]string, 3)
	reviewers[0] = expected_reviewer
	reviewers[1] = expected_reviewer
	reviewers[2] = expected_reviewer
	return
}
