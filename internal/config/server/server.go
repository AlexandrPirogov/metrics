package server

import (
	"log"

	"github.com/caarlos0/env/v7"
	"github.com/spf13/cobra"
)

type Interval string

// default values
var (
	DefaultFileStore     = ""
	DefaultStoreInterval = ""
	DefaultHost          = ""
	DefaultRestore       = true
)

// commands

var (
	rootServerCmd = &cobra.Command{
		Use:   "server",
		Short: "A server that working with metrics",
		Long:  `Some long decrs`,
	}
)

// flags
var (
	Address       string // agent & server addr
	Restore       bool   // Should db be restored
	StoreInterval string // period of replication
	StoreFile     string // file where replication is goint to be written
)

// Configs
var (
	ServerCfg  = ServerConfig{}  // Config for server
	JournalCfg = JournalConfig{} //Config for replication
)

type ServerConfig struct {
	Address string `env:"ADDRESS" envDefault:"localhost:8080"`
}

type JournalConfig struct {
	StoreFile    string `env:"STORE_FILE"`
	Restore      bool   `env:"RESTORE" envDefault:"true"`
	ReadInterval string `env:"STORE_INTERVAL" envDefault:"300s"`
}

func Exec() {

	if err := env.Parse(&ServerCfg); err != nil {
		log.Fatalf("error while read server env variables %v", err)
	}

	if err := env.Parse(&JournalCfg); err != nil {
		log.Fatalf("error while read journal env variables %v", err)
	}

	if err := rootServerCmd.Execute(); err != nil {
		log.Fatalf("%v", err)
	}

	if Address != DefaultHost {
		ServerCfg.Address = Address
	}

	if StoreInterval != DefaultStoreInterval {
		JournalCfg.ReadInterval = StoreInterval
	}
	if StoreFile != DefaultFileStore {
		JournalCfg.StoreFile = StoreFile
	}
}

func init() {
	rootServerCmd.PersistentFlags().StringVarP(&StoreInterval, "interval", "i", DefaultStoreInterval, "Interval of replication")
	rootServerCmd.PersistentFlags().StringVarP(&StoreFile, "file", "f", DefaultFileStore, "File to replicate")
	rootServerCmd.PersistentFlags().BoolVarP(&Restore, "restore", "r", DefaultRestore, "Should restore DB")
	rootServerCmd.PersistentFlags().StringVarP(&Address, "address", "a", DefaultHost, "ADDRESS OF SERVER. Default value: localhost:8080")

	if err := rootServerCmd.Execute(); err != nil {
		log.Fatalf("%v", err)
	}
}