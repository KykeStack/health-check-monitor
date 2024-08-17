package listener

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/KykeStack/health-check-monitor/monitor"
)

type StatusChangeListener interface {
	OnStatusChange(monitor *monitor.URLMonitor) error
	OnRegister() error
}

func parseMessage(msg string, mntr monitor.URLMonitor) (string, error) {
	tmpl, err := template.New("test").Parse(msg)

	if err != nil {
		return "", fmt.Errorf("could not create message. error in template: %v", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, mntr)

	return buf.String(), nil
}
