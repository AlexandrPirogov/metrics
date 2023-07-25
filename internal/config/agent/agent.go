package agent

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"log"
	f "memtracker/internal/function"
	"net/http"
	"os"

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
	cfgFile        string // path to config json file
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
	Address        string   `env:"ADDRESS" envDefault:"localhost:8080" json:"address"`
	ReportInterval Interval `env:"REPORT_INTERVAL" envDefault:"10s" json:"report_interval"`
	PollInterval   Interval `env:"POLL_INTERVAL" envDefault:"2s" json:"poll_interval"`
	Hash           string   `env:"KEY"`
	Limit          int      `env:"RATE_LIMIT" envDefault:"1"`
	CryptoKey      string   `env:"CRYPTO_KEY" json:"crypto_key"`
	TransportCfg   *http.Transport
}

var (
	assignNonTLS = func() {
		ClientCfg.Address = "http://" + ClientCfg.Address
		ClientCfg.TransportCfg = &http.Transport{}
	}
	assignTLS = func() {
		if crt, err := certTemplate(ClientCfg.CryptoKey); err == nil {
			ClientCfg.TransportCfg = &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
					Certificates:       []tls.Certificate{crt},
				},
			}
			ClientCfg.Address = "https://" + ClientCfg.Address
			return
		}
		ClientCfg.Address = "http://" + ClientCfg.Address
		ClientCfg.TransportCfg = &http.Transport{}
	}
)

func Exec() {
	initEnv()
	initFlags()

}

func initEnv() {
	err := env.Parse(&ClientCfg)
	f.ErrFatalCheck("error while read client env variables", err)
}

func initFlags() {

	rootClientCmd.PersistentFlags().StringVarP(&address, "address", "a", "", "ADDRESS OF AGNET. Default value: localhost:8080")
	rootClientCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "PATH TO CONFIG FILE")
	rootClientCmd.PersistentFlags().StringVarP(&reportInterval, "report", "r", "", "How ofter sends metrics to server. Examples: 0s, 10s, 100s")
	rootClientCmd.PersistentFlags().StringVarP(&pollInterval, "poll", "p", "", "How often metrics are updates. Examples: 0s, 10s, 100s")
	rootClientCmd.PersistentFlags().StringVarP(&hash, "key", "k", "", "key for encrypt data that's passes to server")
	rootClientCmd.PersistentFlags().IntVarP(&limit, "limit", "l", 1, "rps limit to send requests")

	err := rootClientCmd.Execute()
	f.ErrFatalCheck("", err)

	f.CompareStringsDo(cfgFile, "", func() { readConfigFile(cfgFile) })
	f.CompareStringsDo(address, "", func() { ClientCfg.Address = address })
	f.CompareStringsDo(hash, "", func() { ClientCfg.Hash = hash })
	f.CompareStringsDo(ClientCfg.CryptoKey, "", func() {})
	f.CompareIntsDo(limit, 1, func() { ClientCfg.Limit = limit })
	f.CompareStringsDoOthewise(cfgFile, "", assignTLS, assignNonTLS)
}

func certTemplate(clientKet string) (tls.Certificate, error) {
	crt, err := tls.LoadX509KeyPair("client.pem", clientKet)
	if err != nil {
		return crt, err
	}
	//f.ErrFatalCheck("err while loading x509 key pair", err)

	certPool, err := x509.SystemCertPool()
	f.ErrFatalCheck("system certpool", err)

	caCertPem, err := os.ReadFile("client.pem")
	f.ErrFatalCheck("err while reading client.pem", err)

	if ok := certPool.AppendCertsFromPEM(caCertPem); !ok {
		log.Fatal("invalid cert in CA PEM")
	}

	return crt, nil
}

func readConfigFile(path string) {
	bytes, err := os.ReadFile(path)
	f.ErrFatalCheck("err while reading config json file", err)

	log.Printf("found agent config file %s", path)
	err = json.Unmarshal(bytes, &ClientCfg)
	f.ErrFatalCheck("", err)

	log.Printf("applied agent config file %s", path)
}
