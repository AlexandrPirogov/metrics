package config

import (
	"log"

	"github.com/caarlos0/env/v7"
	"github.com/spf13/cobra"
)

type Interval string

// commands

var (
	rootServerCmd = &cobra.Command{
		Use:   "server",
		Short: "A server that working with metrics",
		Long:  `Some long decrs`,
	}

	rootClientCmd = &cobra.Command{
		Use:   "agent",
		Short: "Collects metrics and sends to server",
		Long:  `Some long decrs`,
	}
)

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

func init() {
	initFlags()
	initEnvVars()
}

func initFlags() {
	//rootClientCmd.PersistentFlags().StringVarP(&Address, "address", "a", "localhost:8080", "ADDRESS OF AGNET. Default value: localhost:8080")
	rootClientCmd.PersistentFlags().StringVarP(&ReportInterval, "report", "r", "10s", "How ofter sends metrics to server. Examples: 0s, 10s, 100s")
	rootClientCmd.PersistentFlags().StringVarP(&PollInterval, "poll", "p", "10s", "How often metrics are updates. Examples: 0s, 10s, 100s")

	rootServerCmd.PersistentFlags().StringVarP(&StoreInterval, "interval", "i", "0s", "Interval of replication")
	rootServerCmd.PersistentFlags().StringVarP(&StoreFile, "file", "f", "./logs.json", "File to replicate")
	rootServerCmd.PersistentFlags().BoolVarP(&Restore, "restore", "r", true, "Should restore DB")
	rootServerCmd.PersistentFlags().StringVarP(&Address, "address", "a", "localhost:8080", "ADDRESS OF SERVER. Default value: localhost:8080")

	ClientCfg.Address = Address
	ClientCfg.PollInterval = Interval(PollInterval)

	ServerCfg.Address = Address

	JournalCfg.ReadInterval = StoreFile
	JournalCfg.Restore = Restore
	JournalCfg.StoreFile = StoreFile

}

func initEnvVars() {
	if err := env.Parse(&ClientCfg); err != nil {
		log.Printf("error while read client env variables %v", err)
	}

	if err := env.Parse(&ServerCfg); err != nil {
		log.Printf("error while read server env variables %v", err)
	}

	if err := env.Parse(&JournalCfg); err != nil {
		log.Printf("error while read journal env variables %v", err)
	}
}
