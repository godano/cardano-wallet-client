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

// These env-var names are defined as `var` and not `const` on purpose,
// them configurabled by client code, if necessary.
var (
	EnvTLSSkipVerify  = "GODANO_WALLET_CLIENT_TLS_SKIP_VERIFY"
	EnvServerCAFile   = "GODANO_WALLET_CLIENT_SERVER_CA"
	EnvClientCertFile = "GODANO_WALLET_CLIENT_CLIENT_CERT"
	EnvClientKeyFile  = "GODANO_WALLET_CLIENT_CLIENT_KEY"

	// This env-var is not evaluated by MakeTLSConfig(), but in client_test.go and cmd/godano-wallet-cli
	EnvVarWalletServerAddress = "GODANO_WALLET_CLIENT_SERVER_ADDRESS"
)

// NewHTTPSClient returns a `Client` with the given TLS configuration.
func NewHTTPSClient(server string, tlsConfig *tls.Config) (*Client, error) {
	return NewClient(server, WithHTTPSClient(tlsConfig))
}

// NewHTTPSClient returns a `ClientWithResponse` with the given TLS configuration.
func NewHTTPSClientWithResponses(server string, tlsConfig *tls.Config) (*ClientWithResponses, error) {
	return NewClientWithResponses(server, WithHTTPSClient(tlsConfig))
}

// WithHTTPSClient returns a `ClientOption` that sets the given TLS configuration on clients.
func WithHTTPSClient(tlsConfig *tls.Config) ClientOption {
	// Default Transport values copied from http package, TLSClientConfig modified. Avoid copying http.DefaultTransport.
	return WithHTTPClient(&http.Client{
		Transport: &http.Transport{
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
		},
	})
}

// MakeTLSConfig creates a *tls.Config objects based on the GODANO_WALLET_CLIENT_* environment
// variables defined above.
func MakeTLSConfig() (*tls.Config, error) {
	tlsConfig := new(tls.Config)

	if skipVerifyStr := os.Getenv(EnvTLSSkipVerify); skipVerifyStr != "" {
		skipVerify, err := strconv.ParseBool(skipVerifyStr)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse env-var %v=%v as bool: %v",
				EnvTLSSkipVerify, skipVerifyStr, err)
		} else {
			tlsConfig.InsecureSkipVerify = skipVerify
		}
	}

	// Load server CA
	if serverCAFile := os.Getenv(EnvServerCAFile); serverCAFile != "" {
		caRootPool, err := LoadCACert(serverCAFile)
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

	return tlsConfig, nil
}

// LoadCACert loads the given server certificate into a certificate pool, which can be
// set in the `RootCAs` field of `tls.Config`.
func LoadCACert(caFileName string) (*x509.CertPool, error) {
	caCert, err := ioutil.ReadFile(caFileName)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	return caCertPool, nil
}
