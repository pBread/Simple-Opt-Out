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

// payload sent by the inbound SMS webook
type SMSEvent struct {
	AccountSid    string `json:"AccountSid"`
	APIVersion    string `json:"ApiVersion"`
	Body          string `json:"Body"`
	From          string `json:"From"`
	FromCity      string `json:"FromCity"`
	FromCountry   string `json:"FromCountry"`
	FromState     string `json:"FromState"`
	FromZip       string `json:"FromZip"`
	MessageSid    string `json:"MessageSid"`
	NumMedia      string `json:"NumMedia"`
	NumSegments   string `json:"NumSegments"`
	SMSMessageSid string `json:"SmsMessageSid"`
	SMSSid        string `json:"SmsSid"`
	SMSStatus     string `json:"SmsStatus"`
	To            string `json:"To"`
	ToCity        string `json:"ToCity"`
	ToCountry     string `json:"ToCountry"`
	ToState       string `json:"ToState"`
	ToZip         string `json:"ToZip"`
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
