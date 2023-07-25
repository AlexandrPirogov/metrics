package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/caarlos0/env/v7"
	"github.com/spf13/cobra"

	f "memtracker/internal/function"
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
	DefaultCfgFile       = ""
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
	cfgFile       string // path to config json file
	storeInterval string // period of replication
	storeFile     string // file where replication is goint to be written
	hash          string // key for hashing
	dbURL         string // url connection for postgres

	restore bool // Should db be restored
)

// Configs
var (
	ServerCfg  = ServerConfig{}  // Config for server
	JournalCfg = JournalConfig{} //Config for replication
)

type ServerConfig struct {
	Address   string `env:"ADDRESS" envDefault:"localhost:8080" json:"address"`
	Hash      string `env:"KEY"`
	DBUrl     string `env:"DATABASE_DSN" json:"database_dsn"`
	CryptoKey string `env:"CRYPTO_KEY" json:"crypto_key"`
	Run       func(serv *http.Server) error
}

type JournalConfig struct {
	StoreFile    string `json:"store_file" env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore      bool   `json:"restore" env:"RESTORE" envDefault:"true"`
	ReadInterval string `json:"store_interval" env:"STORE_INTERVAL" envDefault:"300s"`
}

var (
	serverTLSAssign = func() {
		ServerCfg.Run = func(serv *http.Server) error {
			log.Println("Running non tls server")
			return serv.ListenAndServe()
		}
	}

	serverNonTLSAssign = func() {
		ServerCfg.Run = func(serv *http.Server) error {
			log.Println("Running tls server")
			return serv.ListenAndServeTLS("server.pem", ServerCfg.CryptoKey)
		}
	}
)

func Exec() {
	initEnv()
	initFlags()
}

func initEnv() {
	err := env.Parse(&ServerCfg)
	f.ErrFatalCheck("error while read server env variables", err)

	err = env.Parse(&JournalCfg)
	f.ErrFatalCheck("error while read journal env variables", err)
}

func initFlags() {
	rootServerCmd.PersistentFlags().StringVarP(&address, "address", "a", "", "ADDRESS OF SERVER. Default value: localhost:8080")
	rootServerCmd.PersistentFlags().StringVarP(&storeInterval, "interval", "i", DefaultStoreInterval, "Interval of replication")
	rootServerCmd.PersistentFlags().StringVarP(&storeFile, "file", "f", DefaultFileStore, "File to replicate")
	rootServerCmd.PersistentFlags().BoolVarP(&restore, "restore", "r", DefaultRestore, "Should restore DB")
	rootServerCmd.PersistentFlags().StringVarP(&hash, "key", "k", "", "key for encrypt data that's passes to agent")
	rootServerCmd.PersistentFlags().StringVarP(&dbURL, "db", "d", "", "database url connection")

	err := rootServerCmd.Execute()
	f.ErrFatalCheck("flags error", err)

	f.CompareStringsDo(cfgFile, DefaultCfgFile, func() { readConfigFile(cfgFile) })
	f.CompareStringsDo(address, DefaultHost, func() { ServerCfg.Address = address })
	f.CompareStringsDo(hash, DefaultHash, func() { ServerCfg.Hash = hash })
	f.CompareStringsDo(storeInterval, DefaultStoreInterval, func() { JournalCfg.ReadInterval = storeInterval })
	f.CompareStringsDo(storeFile, DefaultFileStore, func() { JournalCfg.StoreFile = storeFile })
	f.CompareStringsDo(dbURL, DefaultDBURL, func() { ServerCfg.DBUrl = dbURL })
	f.CompareStringsDoOthewise(ServerCfg.CryptoKey, DefaultCryptoKey, serverNonTLSAssign, serverTLSAssign)

}

func readConfigFile(path string) {
	bytes, err := os.ReadFile(path)
	f.ErrFatalCheck("", err)

	err = json.Unmarshal(bytes, &ServerCfg)
	f.ErrFatalCheck("err while reading config", err)

	err = json.Unmarshal(bytes, &JournalCfg)
	f.ErrFatalCheck("err while reading config", err)
}
