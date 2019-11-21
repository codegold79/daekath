// Package slack is a Cloud Function which recieves a query from
// a Slack command and responds with a random phrase from a json file.
package slack

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
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
	w.Header().Set("Content-Type", "application/json")
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
	quotes := [10]string{
		"You are awesome for doing a job well done!",
		"Everyone sees you as an intelligent, kind, hard-working, ethical person.",
		"It’s only after you’ve stepped outside your comfort zone that you begin to change, grow, and transform. ― Roy T. Bennett ",
		"Keep your face to the sunshine and you cannot see a shadow. – Helen Keller",
		"Stay positive in every situation and everything you do, never stop trying, have faith don’t stop due to failure. – Anurag Prakash Ray",
		"Be the reason someone smiles. Be the reason someone feels loved and believes in the goodness in people. – Roy T. Bennett",
		"Optimism is a happiness magnet. If you stay positive, good things and good people will be drawn to you. – Mary Lou Retton",
		"Yesterday is not ours to recover, but tomorrow is ours to win or lose. – Lyndon B. Johnson",
		"Dwell on the beauty of life. Watch the stars, and see yourself running with them. – Marcus Aurelius",
		"I always like to look on the optimistic side of life, but I am realistic enough to know that life is a complex matter. – Walt Disney",
	}

	quote := quotes[rand.Intn(10)]

	attach := attachment{
		Color: "#02e58a",
		Text:  quote,
	}
	message := &Message{
		Text:        "Here is your goodminder: ",
		Attachments: []attachment{attach},
	}
	return message, nil
}
