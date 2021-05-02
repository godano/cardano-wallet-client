package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/godano/cardano-wallet-client/wallet"
	"github.com/spf13/cobra"
)

var objectCommands = make(map[string]*cobra.Command)

func initCommands() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "godano-wallet-cli",
		Short: "CLI for the cardano-wallet REST API",
		Long: `godano-wallet-cli connects to the REST API of a cardano-wallet process and
	translates the CLI commands and parameters to appropriate REST API calls`,
	}

	// TODO make env vars cofigurable through flags
	// Config file? Map flags to accepted env vars?

	return rootCmd
}

func (c *clientMethod) getObjectCommand(rootCmd *cobra.Command) *cobra.Command {
	objectCommand, ok := objectCommands[c.methodObject]
	if !ok {
		objectCommand = &cobra.Command{
			Use:     c.methodObject,
			Short:   fmt.Sprintf("Query and modify %v objects", c.methodObject),
			Aliases: []string{strings.ToLower(c.methodObject)},
		}
		rootCmd.AddCommand(objectCommand)
		objectCommands[c.methodObject] = objectCommand
	}
	return objectCommand
}

func (c *clientMethod) registerCommand(rootCmd *cobra.Command, clientObj wallet.ClientWithResponsesInterface, ctx context.Context) {
	objectCommand := c.getObjectCommand(rootCmd)
	var useByronVariant bool
	cmd := &cobra.Command{
		Use:     c.shortUseString(),
		Short:   fmt.Sprintf("%v operation for %v objects", c.methodVerb, c.methodObject),
		Aliases: []string{strings.ToLower(c.methodVerb)},
		Args:    cobra.ExactArgs(len(c.args)),
		Run: func(cmd *cobra.Command, args []string) {
			for i := range args {
				// TODO do not modify referenced arg objects here...
				c.args[i].value = args[i]
			}
			res, err := c.call(clientObj, ctx, useByronVariant)
			cobra.CheckErr(err)
			if err == nil {
				outputResponse(res)
			}
		},
	}
	if c.byronVariant != nil {
		cmd.Flags().BoolVar(&useByronVariant, "byron", false, "Use the Byron variant of this command")
	}

	objectCommand.AddCommand(cmd)
}

func (c *clientMethod) shortUseString() string {
	use := c.methodVerb
	for _, arg := range c.args {
		use += fmt.Sprintf(" <%v>", arg.name)
	}
	return use
}
