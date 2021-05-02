package wallet

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	EnvServerAdress   = "GODANO_WALLET_CLIENT_SERVER_ADDRESS"
	EnvTLSSkipVerify  = "GODANO_WALLET_CLIENT_TLS_SKIP_VERIFY"
	EnvServerCAFile   = "GODANO_WALLET_CLIENT_SERVER_CA"
	EnvClientCertFile = "GODANO_WALLET_CLIENT_CLIENT_CERT"
	EnvClientKeyFile  = "GODANO_WALLET_CLIENT_CLIENT_KEY"

	DefaulWalletServer = "https://127.0.0.1:44107/v2"
)

func NewWalletClient() (ClientWithResponsesInterface, error) {
	tlsConfig := new(tls.Config)

	if skipVerifyStr := os.Getenv(EnvTLSSkipVerify); skipVerifyStr != "" {
		skipVerify, err := strconv.ParseBool(skipVerifyStr)
		if err != nil {
			Log.Warnf("Failed to parse env-var %v=%v as bool: %v. Assuming false.",
				EnvTLSSkipVerify, skipVerifyStr, err)
		} else {
			tlsConfig.InsecureSkipVerify = skipVerify
		}
	}

	// Load server CA
	if serverCAFile := os.Getenv(EnvServerCAFile); serverCAFile != "" {
		caRootPool, err := loadCACert(serverCAFile)
		if err != nil {
			return nil, fmt.Errorf("Failed to load server CA file '%v': %v", serverCAFile, err)
		}
		tlsConfig.RootCAs = caRootPool
	}

	// Load client certificate
	clientCertFile := os.Getenv(EnvClientCertFile)
	clientKeyFile := os.Getenv(EnvClientKeyFile)
	// Check that either both variables are defined, or none of them
	if (clientCertFile == "") != (clientKeyFile == "") {
		return nil, fmt.Errorf("Either none or both of these env-vars must be defined: %v, %v",
			EnvClientCertFile, EnvClientKeyFile)
	}
	if clientCertFile != "" {
		cert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	address := os.Getenv(EnvServerAdress)
	if address == "" {
		address = DefaulWalletServer
	}

	return NewWalletClientFor(address, tlsConfig)
}

func NewWalletClientFor(server string, tlsConfig *tls.Config) (ClientWithResponsesInterface, error) {
	// Default Transport values copied from http package, TLSClientConfig modified
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       tlsConfig,
	}
	client := &http.Client{
		Transport: transport,
	}

	return NewClientWithResponses(server, WithHTTPClient(client))
}

func loadCACert(caFileName string) (*x509.CertPool, error) {
	caCert, err := ioutil.ReadFile(caFileName)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	return caCertPool, nil
}
