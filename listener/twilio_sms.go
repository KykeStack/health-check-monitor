package listener

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/KykeStack/health-check-monitor/config"
	"github.com/KykeStack/health-check-monitor/monitor"
)

type TwilioSMSListener struct {
	Config config.TwilioSMSConfig
}

func (tl *TwilioSMSListener) OnStatusChange(URLMonitor *monitor.URLMonitor) error {
	client := &http.Client{Timeout: time.Second * time.Duration(tl.Config.Timeout)}

	twilioURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", tl.Config.AccountSID)

	message := tl.Config.Body.Healthy

	currentStatus, _ := URLMonitor.GetCurrentStatus()

	if !currentStatus {
		message = tl.Config.Body.Unhealthy
	}

	body, err := parseMessage(message, *URLMonitor)

	if err != nil {
		return fmt.Errorf("could not create message: %v", err)
	}

	var wg sync.WaitGroup

	wg.Add(len(tl.Config.To))

	for _, number := range tl.Config.To {
		go func(n string, wg *sync.WaitGroup) {
			defer wg.Done()
			data := url.Values{}
			data.Add("From", tl.Config.From)
			data.Add("Body", body)
			data.Add("To", n)
			req, err := http.NewRequest("POST", twilioURL, strings.NewReader(data.Encode()))
			req.SetBasicAuth(tl.Config.AccountSID, tl.Config.AuthToken)
			req.Header.Add("Accept", "application/json")
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			if err != nil {
				log.Printf("Error sending sms to %s: %v", n, err)
				return
			}
			resp, err := client.Do(req)

			if err != nil {
				log.Printf("Error making request to twilio: %v", err)
				return
			}

			defer resp.Body.Close()

			if resp.StatusCode < 200 || resp.StatusCode > 299 {
				log.Printf("Twilio responded with not 2xx http code: %v", resp.Status)
				body, err := ioutil.ReadAll(resp.Body)

				if err != nil {
					log.Printf("Could not decode de body: %v", err)
					return
				}

				log.Printf("Twilio's response on nr %s: %s", data.Get("To"), body)
			}
		}(number, &wg)
	}

	wg.Wait()

	return nil
}

func (tl *TwilioSMSListener) OnRegister() error {
	if tl.Config.AccountSID == "" {
		return fmt.Errorf("no account sid set")
	}

	if tl.Config.AuthToken == "" {
		return fmt.Errorf("no auth token set")
	}

	if tl.Config.From == "" {
		return fmt.Errorf("no from number set")
	}

	if tl.Config.Body.Healthy == "" {
		return fmt.Errorf("no body set for healthy message")
	}

	if tl.Config.Body.Unhealthy == "" {
		return fmt.Errorf("no body set for unhealthy message")
	}

	if len(tl.Config.To) == 0 {
		return fmt.Errorf("no destination numbers found")
	}

	log.Println("Twilio SMS listener registered")

	return nil
}
