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
	Address        string // agent & server addr
	ReportInterval string // how often agent will sends metrics to server
	PollInterval   string // how often agent will updates metrics
)

// Configs
var (
	ClientCfg = ClientConfig{} // Config for agent
)

type ClientConfig struct {
	Address        string   `env:"ADDRESS" envDefault:"localhost:8080"`
	ReportInterval Interval `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   Interval `env:"POLL_INTERVAL" envDefault:"2s"`
}

func Exec() {

	if err := env.Parse(&ClientCfg); err != nil {
		log.Printf("error while read client env variables %v", err)
	}
	log.Printf("Env config:\nclient:%v", ClientCfg)

	rootClientCmd.PersistentFlags().StringVarP(&Address, "address", "a", "localhost:8080", "ADDRESS OF AGNET. Default value: localhost:8080")
	rootClientCmd.PersistentFlags().StringVarP(&ReportInterval, "report", "r", "10s", "How ofter sends metrics to server. Examples: 0s, 10s, 100s")
	rootClientCmd.PersistentFlags().StringVarP(&PollInterval, "poll", "p", "2s", "How often metrics are updates. Examples: 0s, 10s, 100s")

	rootClientCmd.Execute()

	if err := rootClientCmd.Execute(); err != nil {
		log.Fatalf("%v", err)
	}

	if ClientCfg.Address == "" {
		ClientCfg.Address = Address
	}

	log.Printf("Flags config:\nclient:%v", ClientCfg)
}

func init() {
	//initFlags()
	//initEnvVars()

}

func initFlags() {
	rootClientCmd.PersistentFlags().StringVarP(&Address, "address", "a", "localhost:8080", "ADDRESS OF AGNET. Default value: localhost:8080")
	rootClientCmd.PersistentFlags().StringVarP(&ReportInterval, "report", "r", "10s", "How ofter sends metrics to server. Examples: 0s, 10s, 100s")
	rootClientCmd.PersistentFlags().StringVarP(&PollInterval, "poll", "p", "2s", "How often metrics are updates. Examples: 0s, 10s, 100s")

	if err := rootClientCmd.Execute(); err != nil {
		log.Fatalf("%v", err)
	}

	ClientCfg.Address = Address

	log.Printf("Flags config:\nclient:%v", ClientCfg)
}

func initEnvVars() {
	if err := env.Parse(&ClientCfg); err != nil {
		log.Printf("error while read client env variables %v", err)
	}
	log.Printf("Env config:\nclient:%v", ClientCfg)
}
