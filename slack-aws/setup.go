package slack

import (
	"context"
	"os"
)

type configuration struct {
	ProjectID string `json:"PROJECT_ID"`
	Token     string `json:"SLACK_TOKEN"`
}

var config *configuration

func setup(ctx context.Context) {
	if config == nil {
		config.Token = os.Getenv("SLACK_TOKEN")
	}
}
