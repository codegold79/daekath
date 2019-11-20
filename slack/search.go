// Package slack is a Cloud Function which recieves a query from
// a Slack command and responds with a random phrase from a json file.
package slack

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type attachment struct {
	Color     string `json:"color"`
	Title     string `json:"title"`
	TitleLink string `json:"title_link"`
	Text      string `json:"text"`
	ImageURL  string `json:"image_url"`
}

// Message is the a Slack message event.
// see https://api.slack.com/docs/message-formatting
type Message struct {
	Text        string       `json:"text"`
	Attachments []attachment `json:"attachments"`
}

func Goodminder(w http.ResponseWriter, r *http.Request) {
	setup(r.Context())
	if r.Method != "POST" {
		http.Error(w, "Only POST requests are accepted", 405)
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Couldn't parse form", 400)
		log.Fatalf("ParseForm: %v", err)
	}
	if err := verifyWebHook(r.Form); err != nil {
		log.Fatalf("veryifyWebHook: %v", err)
	}
	searchResponse, err := doSearch()
	if err != nil {
		log.Fatalf("doSearch: %v", err)
	}
	w.Header().Set("Content-Type:", "application/json")
	if err = json.NewEncoder(w).Encode(searchResponse); err != nil {
		log.Fatalf("json.Marshal: %v", err)
	}
}

func verifyWebHook(form url.Values) error {
	t := form.Get("token")
	if len(t) == 0 {
		return fmt.Errorf("empty form token")
	}
	if t != config.Token {
		return fmt.Errorf("invalid request/credentials: %q", t[0])
	}
	return nil
}

func doSearch() (*Message, error) {
	attach := attachment{
		Color:    "#3367d6",
		Text:     "A test attachment",
		Title:    "A test title",
		ImageURL: "A test ImageURL",
	}
	message := &Message{
		Text:        fmt.Sprint("This is a test message"),
		Attachments: []attachment{attach},
	}
	return message, nil
}
