package server

import (
	"log"
	"net/http"

	"github.com/caarlos0/env/v7"
	"github.com/spf13/cobra"
)

type Interval string

// default values
const (
	DefaultFileStore     = ""
	DefaultStoreInterval = ""
	DefaultHost          = ""
	DefaultHash          = ""
	DefaultDBURL         = ""
	DefaultCryptoKey     = ""
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
	address       string // agent & server addr
	restore       bool   // Should db be restored
	storeInterval string // period of replication
	storeFile     string // file where replication is goint to be written
	hash          string // key for hashing
	dbURL         string // url connection for postgres
)

// Configs
var (
	ServerCfg  = ServerConfig{}  // Config for server
	JournalCfg = JournalConfig{} //Config for replication
)

type ServerConfig struct {
	Address   string `env:"ADDRESS" envDefault:"localhost:8080"`
	Hash      string `env:"KEY"`
	DBUrl     string `env:"DATABASE_DSN"`
	CryptoKey string `env:"CRYPTO_KEY"`
	Run       func(serv *http.Server) error
}

type JournalConfig struct {
	StoreFile    string `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore      bool   `env:"RESTORE" envDefault:"true"`
	ReadInterval string `env:"STORE_INTERVAL" envDefault:"300s"`
}

func Exec() {
	initEnv()
	initFlags()
}

func initEnv() {
	if err := env.Parse(&ServerCfg); err != nil {
		log.Fatalf("error while read server env variables %v", err)
	}

	if err := env.Parse(&JournalCfg); err != nil {
		log.Fatalf("error while read journal env variables %v", err)
	}
}

func initFlags() {
	rootServerCmd.PersistentFlags().StringVarP(&storeInterval, "interval", "i", DefaultStoreInterval, "Interval of replication")
	rootServerCmd.PersistentFlags().StringVarP(&storeFile, "file", "f", DefaultFileStore, "File to replicate")
	rootServerCmd.PersistentFlags().BoolVarP(&restore, "restore", "r", DefaultRestore, "Should restore DB")
	rootServerCmd.PersistentFlags().StringVarP(&address, "address", "a", DefaultHost, "ADDRESS OF SERVER. Default value: localhost:8080")
	rootServerCmd.PersistentFlags().StringVarP(&hash, "key", "k", "", "key for encrypt data that's passes to agent")
	rootServerCmd.PersistentFlags().StringVarP(&dbURL, "db", "d", "", "database url connection")

	if err := rootServerCmd.Execute(); err != nil {
		log.Fatalf("%v", err)
	}

	if address != DefaultHost {
		ServerCfg.Address = address
	}

	if hash != "" {
		ServerCfg.Hash = hash
	}

	if storeInterval != DefaultStoreInterval {
		JournalCfg.ReadInterval = storeInterval
	}

	if storeFile != DefaultFileStore {
		JournalCfg.StoreFile = storeFile
	}

	if ServerCfg.DBUrl == DefaultDBURL {
		ServerCfg.DBUrl = dbURL
	}

	if ServerCfg.CryptoKey == DefaultCryptoKey {
		ServerCfg.Run = func(serv *http.Server) error {
			log.Println("Running non tls server")
			return serv.ListenAndServe()
		}
	} else {
		ServerCfg.Run = func(serv *http.Server) error {
			log.Println("Running tls server")
			return serv.ListenAndServeTLS("server.pem", ServerCfg.CryptoKey)
		}
	}
}
