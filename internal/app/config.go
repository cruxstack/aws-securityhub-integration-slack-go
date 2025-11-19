package app

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DebugEnabled           bool
	AwsConsoleURL          string
	AwsAccessPortalURL     string
	AwsAccessRoleName      string
	AWSSecurityHubv2Region string
	SlackToken             string
	SlackChannel           string
}

func NewConfig() (*Config, error) {
	debugEnabled, _ := strconv.ParseBool(os.Getenv("APP_DEBUG_ENABLED"))

	cfg := Config{
		DebugEnabled:           debugEnabled,
		AwsConsoleURL:          os.Getenv("APP_AWS_CONSOLE_URL"),
		AwsAccessPortalURL:     os.Getenv("APP_AWS_ACCESS_PORTAL_URL"),
		AwsAccessRoleName:      os.Getenv("APP_AWS_ACCESS_ROLE_NAME"),
		AWSSecurityHubv2Region: os.Getenv("APP_AWS_SECURITYHUBV2_REGION"),
		SlackToken:             os.Getenv("APP_SLACK_TOKEN"),
		SlackChannel:           os.Getenv("APP_SLACK_CHANNEL"),
	}

	if cfg.AwsConsoleURL == "" {
		cfg.AwsConsoleURL = "https://console.aws.amazon.com"
	}

	var missing []string
	if cfg.SlackToken == "" {
		missing = append(missing, "app_slack_token")
	}
	if cfg.SlackChannel == "" {
		missing = append(missing, "app_slack_channel")
	}
	if len(missing) > 0 {
		return &Config{}, fmt.Errorf("missing env vars: %s", strings.Join(missing, ", "))
	}

	return &cfg, nil
}
