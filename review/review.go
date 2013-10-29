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

type review_request struct {
	from        string
	to          string
	message     string
	review_link string
}

//Takes a json payload request from GitHub and attempts to generate a Review request from it.
func GenerateReviewRequest(payload interface{}, reviewers []string) (rr review_request, err error) {
	log.Println("GenerateReviewRequest")

	defer func() {
		if r := recover(); r != nil {
			err = errors.New("Unable to parse properties from github request.")
		}
	}()

	rr = parsePropertiesAndRandomGenReviewer(payload, reviewers)
	return
}

func SendReviewRequestEmail(request review_request) {
	toAddresses := []string{"abritcliffe@sportingsolutions.com"}
	var output bytes.Buffer
	subject := "A code review request"
	output.WriteString(fmt.Sprintf("Hi %v: \n\n", request.to))
	output.WriteString("Congratulations, it's you're luck day.\n\n")
	output.WriteString(fmt.Sprintf("You can review the commits here: %v \n\n", request.review_link))
	output.WriteString("Happy reviewing!!!\n\n")
	sendMail(output.Bytes(), subject, "abritcliffe@sportingindex.com", toAddresses)
}

func parsePropertiesAndRandomGenReviewer(payload interface{}, reviewers []string) review_request {
	pusher := payload.(map[string]interface{})["pusher"]
	review_link := payload.(map[string]interface{})["compare"].(string)
	return review_request{from: pusher.(map[string]interface{})["name"].(string), to: generateReviewer(reviewers), message: "Please review", review_link: review_link}
}

func generateReviewer(reviewers []string) string {
	count := len(reviewers)
	if count == 0 {
		log.Println("No reviewers provided. Exiting...")
		os.Exit(1)
	}

	rand.Seed(time.Now().UTC().UnixNano())
	r := rand.Intn(count)
	log.Printf("%v index randomly selected", r)

	return reviewers[r]
}

func sendMail(output []byte, subject string, fromAddress string, toAddresses []string) {
	smtpServer := "smtp.mailgun.org"
	from := mail.Address{"Random Review Request", "pmaalert-noreply@sportingsolutions.com"}
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
