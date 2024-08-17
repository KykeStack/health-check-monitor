package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/KykeStack/health-check-monitor/config"
	"github.com/KykeStack/health-check-monitor/listener"
	"github.com/KykeStack/health-check-monitor/monitor"

	"github.com/joho/godotenv"
)

func main() {
	// Find .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	// Getting and using a value from .env
	greeting := os.Getenv("GREETING")

	fmt.Println(greeting)
	configFile := flag.String("config", "config.json", "Path to the configuration file")

	flag.Parse()

	config, err := config.CreateConfigurationFromFile(*configFile)

	if err != nil {
		log.Fatalf("Could not load file: %v", err)
	}

	checker := Checker{Config: config}

	// Register slack listener
	if config.Slack.Enabled {
		err = checker.RegisterListener(&listener.SlackListener{Config: config.Slack})

		if err != nil {
			// We continue to monitor, even if we can't register slack listener
			log.Printf("Warning: %v", err)
		}
	}

	// Register twilio sms listener
	if config.Twilio.SMS.Enabled {
		err = checker.RegisterListener(&listener.TwilioSMSListener{
			Config: config.Twilio.SMS,
		})

		if err != nil {
			// We continue to monitor, even if we can't register twilio sms listener
			log.Printf("Warning: %v", err)
		}
	}

	for _, provider := range config.URLMonitors {
		err = checker.RegisterProvider(&monitor.URLMonitor{
			Name: provider.Name,
			URL:  provider.URL,
		})

		if err != nil {
			log.Fatalf("Error registering monitor: %v", err)
		}

		log.Printf("Provider %s registered", provider.Name)
	}

	go checker.Run()

	statusHandler := StatusHandler{Checker: &checker}

	http.HandleFunc(config.Server.Endpoint, statusHandler.Handle)

	port := fmt.Sprintf(":%d", config.Server.Port)
	log.Printf("Listening on port %d...", config.Server.Port)
	log.Fatal(http.ListenAndServe(port, nil))
}

type Response struct {
	Providers   []*monitor.URLMonitor `json:"providers"`
	Version     string                `json:"version"`
	Environment string                `json:"environment"`
}

type StatusHandler struct {
	Checker *Checker
}

func (sh *StatusHandler) Handle(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Providers: sh.Checker.URLMonitors,
		Version:   os.Getenv("APPLICATION_VERSION"),
	}

	jsonResponse, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(jsonResponse)

	if err != nil {
		panic(err)
	}
}
