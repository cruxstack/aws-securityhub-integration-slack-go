package app

import (
	"encoding/json"
	"fmt"

	awsEvent "github.com/aws/aws-lambda-go/events"
	"github.com/cruxstack/aws-securityhub-integration-slack-go/internal/events"
	"github.com/slack-go/slack"
)

type App struct {
	Config      *Config
	SlackClient *slack.Client
}

func New(cfg *Config) *App {
	return &App{
		Config:      cfg,
		SlackClient: slack.New(cfg.SlackToken),
	}
}

type EventDetail struct {
	Findings []json.RawMessage `json:"findings"`
}

func (a *App) ParseEvent(e awsEvent.CloudWatchEvent) (events.SecurityHubEvent, error) {
	switch e.DetailType {
	case "Findings Imported V2":
		var detail EventDetail
		if err := json.Unmarshal(e.Detail, &detail); err != nil {
			return nil, fmt.Errorf("failed to unmarshal event detail: %w", err)
		}
		if len(detail.Findings) == 0 {
			return nil, fmt.Errorf("no findings in event")
		}
		return events.NewSecurityHubFinding(detail.Findings[0])
	default:
		return nil, fmt.Errorf("unknown cloudwatch event type: %s", e.DetailType)
	}
}

func (a *App) Process(evt awsEvent.CloudWatchEvent) error {
	e, err := a.ParseEvent(evt)
	if err != nil || !e.IsAlertable() {
		return err
	}
	m0, m1 := e.SlackMessage(a.Config.AwsConsoleURL, a.Config.AwsAccessPortalURL, a.Config.AwsAccessRoleName, a.Config.AWSSecurityHubv2Region)
	_, _, err = a.SlackClient.PostMessage(a.Config.SlackChannel, m0, m1)
	return err
}
