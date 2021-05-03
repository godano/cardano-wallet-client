package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/ghodss/yaml"
	"github.com/godano/cardano-wallet-client/wallet"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type walletCLI struct {
	log            *logrus.Logger
	ctx            context.Context
	objectCommands map[string]*cobra.Command

	serverAddress string
	logVerbose    bool
	logQuiet      bool
	logVeryQuiet  bool
	outputYAML    bool
}

func main() {
	cli := walletCLI{
		log:            logrus.StandardLogger(),
		ctx:            context.Background(),
		objectCommands: make(map[string]*cobra.Command),

		// Read default value from the environment
		serverAddress: os.Getenv(wallet.EnvVarWalletServerAddress),
	}
	cli.configureEarlyLogLevel()
	cli.log.SetFormatter(newLogFormatter())

	cobra.OnInitialize(func() {
		cli.configureLogLevel()
	})
	rootCmd := cli.rootCommand()

	// Construct commands based on methods in the client object
	// This does not yet connect to the server
	pseudoClient := new(wallet.Client)
	methods := cli.inspectClientInterface(pseudoClient)
	for _, method := range methods {
		method.mergeByronVariant(methods)
	}
	objectVerbs := make(map[string][]string)
	for _, method := range methods {
		objectVerbs[method.methodObject] = append(objectVerbs[method.methodObject], method.methodVerb)
	}
	for _, method := range methods {
		method.registerCommand(rootCmd, pseudoClient, objectVerbs)
	}

	rootCmd.Execute() // The returned error is already printed by Cobra itself
}

func (c *walletCLI) checkErr(err interface{}) {
	if err != nil {
		c.log.Fatal(err)
	}
}

func (c *walletCLI) connectClient() (*wallet.Client, error) {
	tlsConfig, err := wallet.MakeTLSConfig()
	if err != nil {
		return nil, err
	}
	return wallet.NewHTTPSClient(c.serverAddress, tlsConfig)
}

func (c *walletCLI) outputResponse(response *http.Response) {
	success := response.StatusCode >= http.StatusOK && response.StatusCode < http.StatusMultipleChoices
	if success {
		c.log.Debugf("Response status: %v", response.Status)
	} else {
		c.log.Errorf("Response status: %v", response.Status)
	}

	// 1. Fully read the response body
	var body interface{}
	bodyContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c.log.Errorf("Failed to read response body: %v", err)
		return
	}

	// 2. Parse the response body as JSON
	err = json.Unmarshal(bodyContent, &body)
	if err != nil {
		c.log.Errorf("Failed to unmarshal response body: %v", err)
		fmt.Println(string(bodyContent))
		return
	}

	// 3. Marshall the data again, this time with pretty-printing (JSON or YAML)
	marshalled, err := json.MarshalIndent(body, "", "    ")
	if err != nil {
		c.log.Errorf("Failed to JSON-marshal object: %v", err)
		return
	}
	if c.outputYAML {
		marshalled, err = yaml.JSONToYAML(marshalled)
		if err != nil {
			c.log.Errorf("Failed convert JSON to YAML: %v", err)
			return
		}
	} else {
		marshalled = append(marshalled, '\n') // Properly end JSON output
	}
	fmt.Print(string(marshalled))
}
