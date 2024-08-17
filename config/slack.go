package config

type SlackConfig struct {
	Enabled    bool
	WebhookURL string
	Messages   SlackMessagesConfig
}

type SlackMessagesConfig struct {
	Healthy   string
	Unhealthy string
}
