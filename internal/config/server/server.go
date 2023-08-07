package server

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/caarlos0/env/v7"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/credentials"

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
	DefaultCfgFile       = "/tmp/devops-metrics-db.json"
	DefaultRestore       = true
	DefaultSubnet        = ""
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
	tlscryptokey  string // tls file
	dbURL         string // url connection for postgres
	hash          string // key for hashing
	storeInterval string // period of replication
	storeFile     string // file where replication is goint to be written
	subnet        string //subneting

	restore bool // Should db be restored
	rpc     bool //will the server use rpc

)

// Configs
var (
	ServerCfg  = &ServerConfig{}  // Config for server
	JournalCfg = &JournalConfig{} //Config for replication
)

type ServerConfig struct {
	Address   string `env:"ADDRESS" envDefault:"localhost:8080" json:"address"`
	Hash      string `env:"KEY"`
	DBUrl     string `env:"DATABASE_DSN" json:"database_dsn"`
	CryptoKey string `env:"CRYPTO_KEY" json:"crypto_key"`
	Run       func(serv *http.Server) error
	Subnet    string `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
	RPC       bool   `env:"RPC" envDefault:"false" json:"rpc"`
}

type JournalConfig struct {
	StoreFile    string `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json" json:"store_file" `
	Restore      bool   `env:"RESTORE" envDefault:"true" json:"restore"`
	ReadInterval string `env:"STORE_INTERVAL" envDefault:"300s" json:"store_interval"`
}

// HTTP TLS
var (
	serverNonTLSAssign = func() {
		ServerCfg.Run = func(serv *http.Server) error {
			log.Println("Running non tls server")
			return serv.ListenAndServe()
		}
	}

	serverTLSAssign = func() {
		ServerCfg.Run = func(serv *http.Server) error {
			log.Println("Running tls server")
			return serv.ListenAndServeTLS("server.pem", ServerCfg.CryptoKey)
		}
	}
)

// RPC TLS

func LoadRPCTLSCredentials(key string) (credentials.TransportCredentials, error) {
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair("server.pem", key)
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}

func Exec() {
	initEnv()
	initFlags()
}

func initEnv() {
	errEnvServer := env.Parse(ServerCfg)
	f.ErrFatalCheck("error while read server env variables", errEnvServer)

	errEnvJournal := env.Parse(JournalCfg)
	f.ErrFatalCheck("error while read journal env variables", errEnvJournal)
}

func initFlags() {

	//rootServerCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "json config")
	rootServerCmd.PersistentFlags().StringVarP(&storeInterval, "interval", "i", DefaultStoreInterval, "Interval of replication")
	rootServerCmd.PersistentFlags().StringVarP(&storeFile, "file", "f", DefaultFileStore, "File to replicate")
	rootServerCmd.PersistentFlags().BoolVarP(&restore, "restore", "r", DefaultRestore, "Should restore DB")
	rootServerCmd.PersistentFlags().StringVarP(&address, "address", "a", DefaultHost, "ADDRESS OF SERVER. Default value: localhost:8080")
	rootServerCmd.PersistentFlags().StringVarP(&hash, "key", "k", "", "key for encrypt data that's passes to agent")
	rootServerCmd.PersistentFlags().StringVarP(&dbURL, "db", "d", "", "database url connection")
	rootServerCmd.PersistentFlags().StringVarP(&subnet, "subnet", "t", "", "trusted subnet")
	rootServerCmd.PersistentFlags().StringVarP(&tlscryptokey, "tls", "c", "", "key file for tls")
	rootServerCmd.PersistentFlags().BoolVarP(&rpc, "rpc", "s", false, "set true if you want to use rpc")

	if err := rootServerCmd.Execute(); err != nil {
		log.Fatalf("%v", err)
	}

	ServerCfg.RPC = rpc

	f.CompareStringsDo(cfgFile, DefaultCfgFile, func() { readConfigFile(cfgFile) })
	f.CompareStringsDo(address, DefaultHost, func() { ServerCfg.Address = address })
	f.CompareStringsDo(hash, "", func() { ServerCfg.Hash = hash })
	f.CompareStringsDo(storeInterval, DefaultStoreInterval, func() { JournalCfg.ReadInterval = storeInterval })
	f.CompareStringsDo(storeFile, DefaultFileStore, func() { JournalCfg.StoreFile = storeFile })
	f.CompareStringsDo(subnet, DefaultSubnet, func() { ServerCfg.Subnet = subnet })

	if ServerCfg.DBUrl == DefaultDBURL {
		ServerCfg.DBUrl = dbURL
	}

	if tlscryptokey == DefaultCryptoKey {
		ServerCfg.Run = func(serv *http.Server) error {
			log.Println("Running non tls server")
			return serv.ListenAndServe()
		}
	} else {
		ServerCfg.CryptoKey = tlscryptokey
		ServerCfg.Run = func(serv *http.Server) error {
			log.Println("Running tls server")
			return serv.ListenAndServeTLS("server.pem", ServerCfg.CryptoKey)
		}
	}
	log.Printf("server cfg from flags: %v", ServerCfg)
	log.Printf("jounra cfg from flasg: %v", JournalCfg)
}

func readConfigFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Println(err)
		return
	}

	err = json.Unmarshal(bytes, &ServerCfg)
	f.ErrFatalCheck("err while reading config", err)

	err = json.Unmarshal(bytes, &JournalCfg)
	f.ErrFatalCheck("err while reading config", err)
}
