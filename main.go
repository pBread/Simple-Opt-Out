package main

import (
	"log"
	"net/http"
	"regexp"

	"github.com/gorilla/schema"
)

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

func main() {
	http.HandleFunc("/handler", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	var decoder = schema.NewDecoder()
	var ev SMSEvent

	r.ParseForm()

	err := decoder.Decode(&ev, r.Form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatal(err)
		return
	}

	if hasKeywords(ev.Body) {
		unsubscribe(ev)
		log.Println("User opted out")

	} else {
		log.Println("User did not opt out")

	}

}

func hasKeywords(body string) bool {
	// do not use STOP or UNSUBSCRIBE
	matched, _ := regexp.Match(`(?i)tester`, []byte(body))
	return matched
}

func unsubscribe(ev SMSEvent) {
	// update neighbors database to flag phone number invalid

}
