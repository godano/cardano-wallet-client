package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"

	"github.com/ghodss/yaml"
	"github.com/godano/cardano-wallet-client/wallet"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type walletCLI struct {
	log *logrus.Logger
	ctx context.Context

	rootCmd             *cobra.Command
	byronCmd            *cobra.Command
	objectCommands      map[string]*cobra.Command
	byronObjectCommands map[string]*cobra.Command

	serverAddress string
	dryRun        bool
	logTrace      bool
	logVerbose    bool
	logQuiet      bool
	logVeryQuiet  bool
	outputYAML    bool
}

func main() {
	cli := walletCLI{
		log:                 logrus.StandardLogger(),
		ctx:                 context.Background(),
		objectCommands:      make(map[string]*cobra.Command),
		byronObjectCommands: make(map[string]*cobra.Command),

		// Read default value from the environment
		serverAddress: os.Getenv(wallet.EnvVarWalletServerAddress),
	}
	cli.configureEarlyLogLevel()
	cli.log.SetFormatter(newLogFormatter())
	cobra.OnInitialize(cli.configureLogLevel)

	// Inspect the client object type and find all methods that we can represent as commands
	// This does not yet connect to the server
	pseudoClient := new(wallet.Client)
	methods := cli.inspectClientInterface(pseudoClient)
	objectVerbs := map[bool]map[string][]string{
		true:  make(map[string][]string),
		false: make(map[string][]string),
	}
	for _, method := range methods {
		objectVerbs[method.isByronMethod][method.object] =
			append(objectVerbs[method.isByronMethod][method.object], method.verb)
	}
	for _, eraVerbs := range objectVerbs {
		for _, verbs := range eraVerbs {
			sort.Strings(verbs)
		}
	}

	// Construct commands based on the inspected methods
	cli.initRootCommand()
	cli.initByronCommand()
	for _, method := range methods {
		cmd := &methodCommand{
			cli:    &cli,
			method: method,
		}
		cmd.verbCommand(objectVerbs[method.isByronMethod][method.object])
	}

	cli.rootCmd.Execute() // The returned error is already printed by Cobra itself
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

	c.outputData(response.Body)
}

func (c *walletCLI) outputData(data io.ReadCloser) {
	// Fully read the request or response data
	content, err := ioutil.ReadAll(data)
	if err != nil {
		c.log.Errorf("Failed to read HTTP body data: %v", err)
		return
	}

	// Parse the request or response body as JSON
	var body interface{}
	err = json.Unmarshal(content, &body)
	if err != nil {
		c.log.Errorf("Failed to unmarshal HTTP body data: %v", err)
		fmt.Println(string(content))
		return
	}

	// Marshall the data again, this time with pretty-printing (JSON or YAML)
	marshalled, err := json.MarshalIndent(body, "", "    ")
	if err != nil {
		c.log.Errorf("Failed to JSON-marshal HTTP body object: %v", err)
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

func (c *walletCLI) outputDryRunRequest(req *http.Request) {
	c.log.Info("Dry-run mode, would have performed the following request:")
	c.log.Infof("%v request to URL: %v", req.Method, req.URL)
	if req.Body != nil {
		c.log.Info("Dumping request body...")
		c.outputData(req.Body)
	}
}
