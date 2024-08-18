package monitor

import "encoding/json"

// THRESHOLD Do keep > 2
const THRESHOLD = 5

// URLMonitor This struct handles status of each url
type URLMonitor struct {
	Name          string
	Status        []bool
	URL           string
	Authetication struct {
		Header string
		Value  string
	}
}

type URLMonitorSafe struct {
	Name   string
	Status []bool
	URL    string
}

func (monitor *URLMonitor) PushStatus(status bool) {
	monitor.Status = append(monitor.Status, status)

	if len(monitor.Status) > THRESHOLD {
		monitor.Status = monitor.Status[1:]
	}
}

func (monitor *URLMonitor) StatusChanged() bool {
	// If not enough statuses pushed, just return false
	if !monitor.IsReady() {
		return false
	}

	lastStatus := monitor.Status[0]
	otherStatuses := monitor.Status[1:]

	// Check if status is changing
	for _, status := range otherStatuses {
		if status != otherStatuses[0] {
			return false
		}
	}

	// Only changed when status stack is stable
	return otherStatuses[0] != lastStatus
}

func (monitor *URLMonitor) GetCurrentStatus() (bool, bool) {
	if !monitor.IsReady() {
		return false, false
	}

	return monitor.Status[len(monitor.Status)-1], true
}

func (monitor *URLMonitor) IsReady() bool {
	return len(monitor.Status) == THRESHOLD
}

func (monitor *URLMonitor) MarshalJSON() ([]byte, error) {
	status, ready := monitor.GetCurrentStatus()

	return json.Marshal(struct {
		Name   string `json:"name"`
		URL    string `json:"url"`
		Status bool   `json:"status"`
		Ready  bool   `json:"ready"`
	}{
		Name:   monitor.Name,
		Status: status,
		URL:    monitor.URL,
		Ready:  ready,
	})
}
