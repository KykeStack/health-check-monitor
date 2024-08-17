package config

// URLMonitorConfig represents each url monitor
type URLMonitorConfig struct {
	URL           string
	Name          string
	Authetication struct {
		Header string
		Value  string
	}
}
