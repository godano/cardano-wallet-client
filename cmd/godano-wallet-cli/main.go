package main

import (
	"context"

	"github.com/godano/cardano-wallet-client/wallet"
	"github.com/spf13/cobra"
)

var Log = wallet.Log

func main() {
	// Log.SetLevel(logrus.DebugLevel)

	client, err := wallet.NewWalletClient()
	cobra.CheckErr(err)

	ctx := context.Background()
	rootCmd := initCommands()
	methods := inspectClientInterface(client)
	for _, method := range methods {
		method.mergeByronVariant(methods)
	}
	for _, method := range methods {
		method.registerCommand(rootCmd, client, ctx)
	}

	cobra.CheckErr(rootCmd.Execute())
}
