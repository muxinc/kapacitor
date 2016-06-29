package muxincident

const DefaultMuxIncidentAPIURL = "https://app.mux.com/"

type Config struct {
	// Whether MuxIncident integration is enabled.
	Enabled bool `toml:"enabled"`
	// The MuxIncident API URL, should not need to be changed.
	URL string `toml:"url"`
	// The MuxIncident username.
	Username string `toml:"username"`
	// The MuxIncident password.
	Password string `toml:"password"`
	// Whether every alert should automatically go to PagerDuty
	Global bool `toml:"global"`
}

func NewConfig() Config {
	return Config{
		URL: DefaultMuxIncidentAPIURL,
	}
}
