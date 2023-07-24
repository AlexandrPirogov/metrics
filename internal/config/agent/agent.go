package agent

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"

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
	CryptoKey      string   `env:"CRYPTO_KEY"`
	TransportCfg   *http.Transport
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

	if ClientCfg.CryptoKey == "" {
		ClientCfg.TransportCfg = &http.Transport{}
	} else {
		ClientCfg.TransportCfg = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				Certificates: []tls.Certificate{
					certTemplate(ClientCfg.CryptoKey),
				},
			},
		}
	}

}

func certTemplate(clientKet string) tls.Certificate {
	crt, err := tls.LoadX509KeyPair("client.pem", "client.key")
	if err != nil {
		log.Fatalf("err while loading x509 key pair: %v", err)
	}
	certPool, err := x509.SystemCertPool()
	if err != nil {
		log.Fatalf("system certpool %v", err)
	}
	caCertPem, err := ioutil.ReadFile("cert.pem")
	if err != nil {
		log.Fatalf("err while reading cert.pem %v", err)
	}
	if ok := certPool.AppendCertsFromPEM(caCertPem); !ok {
		log.Fatal("invalid cert in CA PEM")
	}
	return crt
}
