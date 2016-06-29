package mux

const DefaultMuxAPIURL = "https://app.mux.io/"

type Config struct {
	// Whether Mux integration is enabled.
	Enabled bool `toml:"enabled"`
	// The Mux API URL, should not need to be changed.
	URL string `toml:"url"`
	// The Mux username.
	Username string `toml:"username"`
	// The Mux password.
	Password string `toml:"password"`
	// Whether every alert should automatically go to Mux
	Global bool `toml:"global"`
}

func NewConfig() Config {
	return Config{
		URL: DefaultMuxAPIURL,
	}
}
