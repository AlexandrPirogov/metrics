package config

type Interval string

// flags
var (
	Address       string // agent & server addr
	Restore       bool   // Should db be restored
	StoreInterval string // period of replication
	StoreFile     string // file where replication is goint to be written

	ReportInterval string // how often agent will sends metrics to server
	PollInterval   string // how often agent will updates metrics
)

// Configs
var (
	ClientCfg  ClientConfig  // Config for agent
	ServerCfg  ServerConfig  // Config for server
	JournalCfg JournalConfig //Config for replication
)

type ClientConfig struct {
	Address        string   `env:"ADDRESS" envDefault:"localhost:8080"`
	ReportInterval Interval `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   Interval `env:"POLL_INTERVAL" envDefault:"2s"`
}

type ServerConfig struct {
	Address string `env:"ADDRESS" envDefault:"localhost:8080"`
}

type JournalConfig struct {
	StoreFile    string `env:"STORE_FILE"`
	Restore      bool   `env:"RESTORE" envDefault:"true"`
	ReadInterval string `env:"STORE_INTERVAL" envDefault:"300"`
}
