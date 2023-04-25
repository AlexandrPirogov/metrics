package agent

import (
	"log"

	"github.com/caarlos0/env/v7"
	"github.com/spf13/cobra"
)

type Interval string

// commands

var (
	rootClientCmd = &cobra.Command{
		Use:   "agent",
		Short: "Collects metrics and sends to server",
		Long:  `Some long decrs`,
	}
)

// flags
var (
	address        string // agent & server addr
	reportInterval string // how often agent will sends metrics to server
	pollInterval   string // how often agent will updates metrics
	hash           string //hash for metric
	limit          int    //rate limit for agent to send requests
)

// Configs
var (
	ClientCfg = ClientConfig{} // Config for agent
)

type ClientConfig struct {
	Address        string   `env:"ADDRESS" envDefault:"localhost:8080"`
	ReportInterval Interval `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   Interval `env:"POLL_INTERVAL" envDefault:"2s"`
	Hash           string   `env:"KEY"`
	Limit          int      `env:"RATE_LIMIT" envDefault:"1"`
}

func Exec() {
	initEnv()
	initFlags()

}

func initEnv() {
	if err := env.Parse(&ClientCfg); err != nil {
		log.Fatalf("error while read client env variables %v", err)
	}
}

func initFlags() {

	rootClientCmd.PersistentFlags().StringVarP(&address, "address", "a", "", "ADDRESS OF AGNET. Default value: localhost:8080")
	rootClientCmd.PersistentFlags().StringVarP(&reportInterval, "report", "r", "", "How ofter sends metrics to server. Examples: 0s, 10s, 100s")
	rootClientCmd.PersistentFlags().StringVarP(&pollInterval, "poll", "p", "", "How often metrics are updates. Examples: 0s, 10s, 100s")
	rootClientCmd.PersistentFlags().StringVarP(&hash, "key", "k", "", "key for encrypt data that's passes to server")
	rootClientCmd.PersistentFlags().IntVarP(&limit, "limit", "l", 1, "rps limit to send requests")

	if err := rootClientCmd.Execute(); err != nil {
		log.Fatalf("%v", err)
	}
	if address != "" {
		ClientCfg.Address = address
	}

	if hash != "" {
		ClientCfg.Hash = hash
	}

	if limit != 1 {
		ClientCfg.Limit = limit
	}

}
