package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// Config represents configuration set for the program to run
type Config struct {
	Server struct {
		Port     int
		Endpoint string
	}
	Checker struct {
		Interval int
		Timeout  int
	}
	URLMonitors []URLMonitorConfig
	Slack       SlackConfig
	Twilio      TwilioConfig
}

// CreateConfigurationFromFile Returns a new configuration loaded from a file
func CreateConfigurationFromFile(configFile string) (Config, error) {
	config := Config{
		Server: struct {
			Port     int
			Endpoint string
		}{
			Port:     8001,
			Endpoint: "/status",
		},
		Checker: struct {
			Interval int
			Timeout  int
		}{
			Interval: 30,
			Timeout:  5,
		},
		URLMonitors: []URLMonitorConfig{},
		Slack: SlackConfig{
			Enabled:    false,
			WebhookURL: "",
			Messages: SlackMessagesConfig{
				Healthy:   "{{.Name}} is up!",
				Unhealthy: "{{.Name}} is down!",
			},
		},
		Twilio: TwilioConfig{
			SMS: TwilioSMSConfig{
				Enabled:    false,
				AccountSID: "",
				AuthToken:  "",
				Body: TwilioSMSBodyConfig{
					Healthy:   "{{.Name}} is up!",
					Unhealthy: "{{.Name}} is down!",
				},
				Timeout: 1,
			},
		},
	}

	if _, err := os.Stat(configFile); !os.IsNotExist(err) {
		log.Printf("Found configuration file: %s\n", configFile)

		configData, err := ioutil.ReadFile(configFile)
		if err != nil {
			return config, fmt.Errorf("error reading configuation file: %v", err)
		}

		err = json.Unmarshal(configData, &config)
		if err != nil {
			return config, fmt.Errorf("error decoding configuration: %v", err)
		}

		log.Print("Configuration loaded successfully\n")
	}

	if len(config.URLMonitors) == 0 {
		return config, fmt.Errorf("no providers found")
	}

	maxTotalTimeout := config.Checker.Timeout * len(config.URLMonitors)

	if maxTotalTimeout >= config.Checker.Interval {
		return config, fmt.Errorf("timeout value (times monitor count) cannot be greater or equal than interval")
	}

	return config, nil
}
