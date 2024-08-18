package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/KykeStack/health-check-monitor/config"
	"github.com/KykeStack/health-check-monitor/listener"
	"github.com/KykeStack/health-check-monitor/monitor"
)

var statusString = map[bool]string{
	true:  "up",
	false: "down",
}

var validStatusCodes = map[int]bool{
	http.StatusCreated: true,
	http.StatusOK:      true,
}

type Checker struct {
	Config      config.Config
	URLMonitors []*monitor.URLMonitor
	Listeners   []listener.StatusChangeListener
}

func (c *Checker) RegisterProvider(URLMonitor *monitor.URLMonitor) error {
	if URLMonitor.Name == "" {
		return fmt.Errorf("monitor is missing name")
	}

	if URLMonitor.URL == "" {
		return fmt.Errorf("monitor %s is missing url", URLMonitor.Name)
	}

	c.URLMonitors = append(c.URLMonitors, URLMonitor)

	return nil
}

func (c *Checker) RegisterListener(l listener.StatusChangeListener) error {
	if err := l.OnRegister(); err != nil {
		return fmt.Errorf("could not register listener: %v", err)
	}

	c.Listeners = append(c.Listeners, l)
	log.Println("Slack listener registered")

	return nil
}

func (c *Checker) Run() {
	ticker := time.NewTicker(time.Second * time.Duration(c.Config.Checker.Interval))

	for {
		select {
		case <-ticker.C:
			c.checkAll()
		}
	}
}

func (c *Checker) checkAll() {
	client := http.Client{
		Timeout: time.Second * time.Duration(c.Config.Checker.Timeout),
	}

	for _, urlMonitor := range c.URLMonitors {
		resp, err := fetchMonitorURL(&client, urlMonitor)

		status := true

		if err != nil {
			log.Printf("Error on request for monitor %s: %v", urlMonitor.Name, err)
			status = false
		}

		if bool(!validStatusCodes[resp]) {
			status = false
		}

		urlMonitor.PushStatus(status)

		if urlMonitor.StatusChanged() && urlMonitor.IsReady() {
			currentStatus, _ := urlMonitor.GetCurrentStatus()
			log.Printf("%s monitor is now %s", urlMonitor.Name, statusString[currentStatus])

			for _, l := range c.Listeners {
				go l.OnStatusChange(urlMonitor)
			}
		}
	}
}

func fetchMonitorURL(client *http.Client, URLMonitor *monitor.URLMonitor) (int, error) {
	req, err := http.NewRequest("GET", URLMonitor.URL, nil)
	if err != nil {
		log.Fatal(err)
	}

	if URLMonitor.Authetication.Header != "" && URLMonitor.Authetication.Value != "" {
		headerValue := os.Getenv(URLMonitor.Authetication.Value)
		req.Header.Add(URLMonitor.Authetication.Header, headerValue)
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("cannot fetch monitor %s: %v", URLMonitor.Name, err)
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}
