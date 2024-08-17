package monitor

import (
	"testing"
)

func TestStatusProvider_TestPushStatus(t *testing.T) {
	urlMonitor := URLMonitor{}

	urlMonitor.PushStatus(true)

	if !urlMonitor.Status[0] {
		t.Error("First value should be true!")
	}

	urlMonitor.PushStatus(true)
	urlMonitor.PushStatus(true)
	urlMonitor.PushStatus(true)
	urlMonitor.PushStatus(true)
	urlMonitor.PushStatus(true)
	urlMonitor.PushStatus(true)
	urlMonitor.PushStatus(true)
	urlMonitor.PushStatus(true)
	urlMonitor.PushStatus(true)

	if len(urlMonitor.Status) != THRESHOLD {
		t.Errorf("You should only save %d statuses! %d found.", THRESHOLD, len(urlMonitor.Status))
	}
}

func TestStatusProvider_StatusChanged(t *testing.T) {
	urlMonitor := URLMonitor{}

	if urlMonitor.StatusChanged() {
		t.Error("Newly created monitor should not have the status changed!")
	}

	// Populate status stack
	for i := 1; i <= THRESHOLD; i++ {
		urlMonitor.PushStatus(false)
	}

	if urlMonitor.StatusChanged() {
		t.Error("Status has not changed!")
	}

	// This only works if THRESHOLD > 2
	urlMonitor.PushStatus(true)

	if urlMonitor.StatusChanged() {
		t.Errorf("Status should only change if last status is different! %v", urlMonitor.Status)
	}

	// Reset to true
	for i := 1; i <= THRESHOLD; i++ {
		urlMonitor.PushStatus(true)
	}

	for i := 1; i < THRESHOLD; i++ {
		urlMonitor.PushStatus(false)
	}

	if !urlMonitor.StatusChanged() {
		t.Errorf("Status should have changed: %v", urlMonitor.Status)
	}
}
