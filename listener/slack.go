package listener

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/KykeStack/health-check-monitor/monitor"

	"github.com/KykeStack/health-check-monitor/config"
)

// SlackListener Will send a slack notification when a monitor status changes
type SlackListener struct {
	Config config.SlackConfig
}

func (sl *SlackListener) OnRegister() error {
	if sl.Config.WebhookURL == "" {
		return fmt.Errorf("slack's webhook url is empty")
	}

	return nil
}

func (sl *SlackListener) OnStatusChange(m *monitor.URLMonitor) error {
	color := "#36a64f"
	message := sl.Config.Messages.Healthy

	status, _ := m.GetCurrentStatus()

	if !status {
		color = "#ea6767"
		message = sl.Config.Messages.Unhealthy
	}

	message, err := parseMessage(message, *m)

	slackMessage := slackMessage{
		Attachments: []slackAttachment{
			{
				Color: color,
				Title: m.Name,
				Text:  message,
			},
		},
	}

	payload, err := json.Marshal(slackMessage)

	if err != nil {
		return fmt.Errorf("could not marshall struct: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, sl.Config.WebhookURL, bytes.NewBuffer(payload))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("could not make request to slack: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("bad response from slack: %v", resp.Status)
	}

	return nil
}

type slackMessage struct {
	Attachments []slackAttachment `json:"attachments"`
}

type slackAttachment struct {
	Color string `json:"color"`
	Title string `json:"title"`
	Text  string `json:"text"`
}
