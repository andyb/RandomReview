package review

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/mail"
	"net/smtp"
	"os"
	"strings"
	"time"
)

type Review_request struct {
	From        string
	From_email  string
	To          Reviewer
	Message     string
	Review_link string
}

type Reviewer struct {
	Name  string
	Email string
}

//Takes a json payload request from GitHub and attempts to generate a Review request from it.
func GenerateReviewRequest(payload interface{}, reviewers []Reviewer) (rr Review_request, err error) {
	log.Println("GenerateReviewRequest")

	defer func() {
		if r := recover(); r != nil {
			err = errors.New("Unable to parse properties from github request.")
		}
	}()

	rr = parsePropertiesAndRandomGenReviewer(payload, reviewers)
	return
}

func SendReviewRequestEmail(request Review_request) {
	log.Printf("Sending review request from %v to %v", request.From, request.To.Name)
	toAddresses := []string{request.To.Email}
	var output bytes.Buffer
	subject := "Connect: Code review request from " + request.From
	output.WriteString(fmt.Sprintf("Hi %v: \n\n", request.To.Name))
	output.WriteString(fmt.Sprintf("Congratulations, it's your lucky day. You've been randomly chosen to do a code review for %v \n\n", request.From))
	output.WriteString(fmt.Sprintf("You can review the commits here: %v \n\n", request.Review_link))
	output.WriteString("Happy reviewing!!!\n\n")
	sendMail(output.Bytes(), subject, request.From_email, toAddresses)
}

func parsePropertiesAndRandomGenReviewer(payload interface{}, reviewers []Reviewer) Review_request {
	pusher := payload.(map[string]interface{})["pusher"]
	review_link := payload.(map[string]interface{})["compare"].(string)
	rev := generateReviewer(reviewers)
	return Review_request{From: pusher.(map[string]interface{})["name"].(string), To: rev, Message: "Please review", Review_link: review_link}
}

func generateReviewer(reviewers []Reviewer) (rev Reviewer) {
	count := len(reviewers)
	if count == 0 {
		log.Println("No reviewers provided. Exiting...")
		os.Exit(1)
	}

	rand.Seed(time.Now().UTC().UnixNano())
	r := rand.Intn(count)
	log.Printf("%v index randomly selected. Reviewer is %v", r, reviewers[r].Name)

	return reviewers[r]
}

func sendMail(output []byte, subject string, fromAddress string, toAddresses []string) {
	smtpServer := "smtp.mailgun.org"
	from := mail.Address{"Code Review Request", "codereviewreq-noreply@sportingsolutions.com"}
	auth := smtp.PlainAuth(
		"",
		"postmaster@sportingsolutions.com",
		"2a312u3v1lq6",
		smtpServer,
	)
	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = toAddresses[0]
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString(output)

	err := smtp.SendMail(smtpServer+":587", auth, from.Address, toAddresses, []byte(message))
	LogError(err)

	log.Printf("Email sent. %s", strings.Join(toAddresses, ","))
}

func LogError(err error) {
	if err != nil {
		log.Printf("Error:  %s", err)
	}
}
