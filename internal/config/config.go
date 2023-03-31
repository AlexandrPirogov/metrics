package config

type Interval string

type ClientConfig struct {
	Address        string   `env:"ADDRESS" envDefault:"localhost:8080"`
	ReportInterval Interval `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   Interval `env:"POLL_INTERVAL" envDefault:"2s"`
}

type ServerConfig struct {
	Address string `env:"ADDRESS" envDefault:"localhost:8080"`
}
