package agent

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	f "memtracker/internal/function"
	"net/http"
	"os"

	"github.com/caarlos0/env/v7"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/credentials"
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
	tlscryptokey   string //key file for tls
	reportInterval string // how often agent will sends metrics to server
	pollInterval   string // how often agent will updates metrics
	hash           string // hash for metric
	limit          int    // rate limit for agent to send requests
	rpc            bool   // is client using rpc
)

// Configs
var (
	ClientCfg = &ClientConfig{} // Config for agent
)

// Functions for TLS configure
var (
	assignNonTLS = func() {
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
			return
		}
		ClientCfg.TransportCfg = &http.Transport{}
	}
)

//RPC TLS

func LoadTLSCredentials(keyFile string) (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate
	pemServerCA, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// Create the credentials and return it
	config := &tls.Config{
		RootCAs: certPool,
	}

	return credentials.NewTLS(config), nil
}

type ClientConfig struct {
	Address        string   `env:"ADDRESS" envDefault:"localhost:8080" json:"address"`
	CryptoKey      string   `env:"CRYPTO_KEY" json:"crypto_key"`
	Hash           string   `env:"KEY"`
	Limit          int      `env:"RATE_LIMIT" envDefault:"1"`
	PollInterval   Interval `env:"POLL_INTERVAL" envDefault:"2s" json:"poll_interval"`
	ReportInterval Interval `env:"REPORT_INTERVAL" envDefault:"10s" json:"report_interval"`
	RPC            bool     `env:"RPC" envDefault:"false" json:"rpc"`
	TransportCfg   *http.Transport
	Protocol       string
}

func Exec() {
	initEnv()
	initFlags()

}

func initEnv() {
	err := env.Parse(ClientCfg)
	f.ErrFatalCheck("error while read client env variables", err)
}

func initFlags() {

	rootClientCmd.PersistentFlags().StringVarP(&address, "address", "a", "", "ADDRESS OF AGNET. Default value: localhost:8080")
	rootClientCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "PATH TO CONFIG FILE")
	rootClientCmd.PersistentFlags().StringVarP(&reportInterval, "report", "r", "", "How ofter sends metrics to server. Examples: 0s, 10s, 100s")
	rootClientCmd.PersistentFlags().StringVarP(&pollInterval, "poll", "p", "", "How often metrics are updates. Examples: 0s, 10s, 100s")
	rootClientCmd.PersistentFlags().StringVarP(&hash, "key", "k", "", "key for encrypt data that's passes to server")
	rootClientCmd.PersistentFlags().IntVarP(&limit, "limit", "l", 1, "rps limit to send requests")
	rootClientCmd.PersistentFlags().StringVarP(&tlscryptokey, "tls", "t", "", "key file for tls")
	rootClientCmd.PersistentFlags().BoolVarP(&rpc, "rpc", "s", false, "set true if you want to use rpc")

	err := rootClientCmd.Execute()
	f.ErrFatalCheck("", err)

	f.CompareStringsDo(cfgFile, "", func() { readConfigFile(cfgFile) })
	f.CompareStringsDo(address, "", func() { ClientCfg.Address = address })
	f.CompareStringsDo(hash, "", func() { ClientCfg.Hash = hash })
	f.CompareStringsDo(tlscryptokey, "", func() { ClientCfg.CryptoKey = tlscryptokey })

	f.CompareStringsDoOthewise(pollInterval, "",
		func() { ClientCfg.PollInterval = Interval(pollInterval) },
		func() { ClientCfg.PollInterval = "2s" })

	f.CompareStringsDoOthewise(reportInterval, "",
		func() { ClientCfg.ReportInterval = Interval(reportInterval) },
		func() { ClientCfg.ReportInterval = "10s" })

	f.CompareIntsDo(limit, 1, func() { ClientCfg.Limit = limit })

	ClientCfg.RPC = rpc

	if !ClientCfg.RPC {
		f.CompareStringsDoOthewise(cfgFile, "", assignNonTLS, assignTLS)
		f.CompareStringsDoOthewise(cfgFile, "", func() { ClientCfg.Protocol = "http://" }, func() { ClientCfg.Protocol = "https://" })
	}
	log.Printf("Agent config: %v", ClientCfg)
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
	if err != nil {
		return
	}

	log.Printf("found agent config file %s", path)
	err = json.Unmarshal(bytes, &ClientCfg)
	f.ErrFatalCheck("", err)

	log.Printf("applied agent config file %s", path)
}
