package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

func (c *walletCLI) rootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "godano-wallet-cli object operation",
		Short: "CLI for the cardano-wallet REST API",
		Long: `godano-wallet-cli connects to the REST API of a cardano-wallet process and
translates the CLI commands and parameters to appropriate REST API calls`,
	}

	flags := rootCmd.PersistentFlags()
	flags.StringVarP(&c.serverAddress, "server", "s", c.serverAddress, "Endpoint of the cardano-wallet process to connect to")
	flags.BoolVarP(&c.logQuiet, "quiet", "q", c.logQuiet, "Set the log level to Warning")
	flags.BoolVarP(&c.logVeryQuiet, "quieter", "Q", c.logVeryQuiet, "Set the log level to Error")
	flags.BoolVarP(&c.logVerbose, "verbose", "v", c.logVerbose, "Set the log level to Debug")
	flags.BoolVarP(&c.outputYAML, "yaml", "y", c.outputYAML, "Output responses as YAML instead of JSON (more compact)")
	return rootCmd
}

func (c *clientMethod) getObjectCommand(rootCmd *cobra.Command, objectVerbs map[string][]string) *cobra.Command {
	objectCommand, ok := c.objectCommands[c.methodObject]
	if !ok {
		// These messages cover the case that there are multiple sub-commands for this object
		verbs := objectVerbs[c.methodObject]
		sort.Strings(verbs)
		var joinedVerbs string
		if len(verbs) == 1 {
			joinedVerbs = verbs[0]
		} else if len(verbs) == 2 {
			joinedVerbs = fmt.Sprintf("%v or %v", verbs[0], verbs[1])
		} else {
			joinedVerbs = strings.Join(verbs[:len(verbs)-1], ", ") + ", or " + verbs[len(verbs)-1]
		}

		shortMessage := fmt.Sprintf("%v %v objects", joinedVerbs, c.methodObject)
		longMessage := shortMessage

		objectCommand = &cobra.Command{
			Use:     c.methodObject,
			Short:   shortMessage,
			Long:    longMessage,
			Aliases: []string{strings.ToLower(c.methodObject)},
		}
		rootCmd.AddCommand(objectCommand)
		c.objectCommands[c.methodObject] = objectCommand
	}
	return objectCommand
}

func (c *clientMethod) registerCommand(rootCmd *cobra.Command, clientObj interface{}, objectVerbs map[string][]string) {
	objectCommand := c.getObjectCommand(rootCmd, objectVerbs)
	isOnlyCommand := len(objectVerbs[c.methodObject]) == 1

	var cmd *cobra.Command
	if isOnlyCommand {
		// For objects with only one verb, there is no sub-command
		c.configureCommand(objectCommand)
		cmd = objectCommand
	} else {
		cmd = new(cobra.Command)
		cmd.Use = c.methodVerb
		c.configureCommand(cmd)
		objectCommand.AddCommand(cmd)
	}

	cmd.Short = fmt.Sprintf("%v %v", c.methodVerb, c.methodObject)
	cmd.Long = fmt.Sprintf("%v operation for %v objects", c.methodVerb, c.methodObject)
}

func (c *clientMethod) configureCommand(cmd *cobra.Command) {
	numArgs := 0
	for _, arg := range c.args {
		if !arg.isStructParameter() {
			// TODO support setting struct params
			numArgs++
			cmd.Use += fmt.Sprintf(" <%v>", arg.name)
		}
	}
	cmd.Args = cobra.ExactArgs(numArgs)

	var useByronVariant bool
	cmd.Run = func(cmd *cobra.Command, args []string) {
		for i := range args {
			// TODO do not modify referenced arg objects here...
			c.args[i].value = args[i]
		}
		client, err := c.connectClient()
		c.checkErr(err)
		res, err := c.call(client, c.ctx, useByronVariant)
		c.checkErr(err)
		if err == nil {
			c.outputResponse(res)
		}
	}

	if c.byronVariant != nil {
		cmd.Use += " [--byron]"
		cmd.Flags().BoolVar(&useByronVariant, "byron", false, "Use the Byron variant of this command")
	}
}
