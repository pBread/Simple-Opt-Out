package main

import (
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/gorilla/schema"
	"github.com/joho/godotenv"
	twilio "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

// payload received by the inbound SMS webook
type SMSEvent struct {
	APIVersion    string `json:"ApiVersion"`
	SMSSid        string `json:"SmsSid"`
	SMSStatus     string `json:"SmsStatus"`
	SMSMessageSid string `json:"SmsMessageSid"`
	NumSegments   string `json:"NumSegments"`
	ToState       string `json:"ToState"`
	From          string `json:"From"`
	MessageSid    string `json:"MessageSid"`
	AccountSid    string `json:"AccountSid"`
	ToCity        string `json:"ToCity"`
	FromCountry   string `json:"FromCountry"`
	ToZip         string `json:"ToZip"`
	FromCity      string `json:"FromCity"`
	To            string `json:"To"`
	FromZip       string `json:"FromZip"`
	ToCountry     string `json:"ToCountry"`
	Body          string `json:"Body"`
	NumMedia      string `json:"NumMedia"`
	FromState     string `json:"FromState"`
}

var client *twilio.RestClient

func main() {
	godotenv.Load(".env")
	creds := twilio.RestClientParams{Username: os.Getenv("ACCOUNT_SID"), Password: os.Getenv("AUTH_TOKEN")}
	client = twilio.NewRestClientWithParams(creds)

	http.HandleFunc("/opt-out/new", handleNew)
	http.HandleFunc("/opt-out/reply", handleReply)
	http.ListenAndServe(":8080", nil)
}

func handleNew(w http.ResponseWriter, r *http.Request) {
	to := r.URL.Query().Get("to")
	sendSms(to, "initial")
	w.Write([]byte("Sent opt out message to" + to))
}

func handleReply(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var ev SMSEvent
	if err := schema.NewDecoder().Decode(&ev, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatal(err)
	}

	if hasKeywords(ev.Body) {
		unsubscribe(ev)
		sendSms(ev.From, "unsubscribe")
		log.Println("User opted out")
	} else {
		sendSms(ev.From, "mistake")
		log.Println("User did not opt out")
	}
}

func hasKeywords(body string) bool {
	// do not use these keywords: https://bit.ly/3m33u6I
	matched, _ := regexp.Match(`(?i)tester`, []byte(body))
	return matched
}

func unsubscribe(ev SMSEvent) {
	// update database to flag post as invalid
}

func sendSms(to string, message string) {
	var body string
	switch message {
	case "initial":
		body = "Hello, your post was successful \n\nIf this is a mistake, respond TESTER"
	case "mistake":
		body = "I did not understand your response. If you did not post a, respond TESTER."
	case "unsubscribe":
		body = "You have been unsubscribed."
	}

	params := &openapi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(os.Getenv("OPT_OUT_PHONE"))
	params.SetBody(body)

	if _, err := client.ApiV2010.CreateMessage(params); err != nil {
		log.Println(err.Error())
	} else {
		log.Println("SMS sent successfully!")
	}
}
