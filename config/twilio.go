package config

type TwilioConfig struct {
	SMS TwilioSMSConfig
}

type TwilioSMSConfig struct {
	Enabled    bool
	AccountSID string
	AuthToken  string
	From       string
	To         []string
	Body       TwilioSMSBodyConfig
	Timeout    int
}

type TwilioSMSBodyConfig struct {
	Healthy   string
	Unhealthy string
}
