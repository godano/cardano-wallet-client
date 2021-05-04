<p align="center">
  <h1 align="center">
    Cardano-Wallet Client
    <br/>
    <a href="https://github.com/godano/cardano-wallet-client/blob/master/LICENSE" ><img alt="license" src="https://img.shields.io/badge/license-MIT%20License%202.0-E91E63.svg?style=flat-square" /></a>
    <a href="https://t.me/godano"><img src="https://img.shields.io/badge/Chat%20on-Telegram-blue.svg"/></a>
  </h1>
</p>

A Go client for the [cardano-wallet](https://github.com/input-output-hk/cardano-wallet) by IOG.

The bulk of this client code is generated using [oapi-codegen](https://github.com/deepmap/oapi-codegen), based on the [Open API definition](https://input-output-hk.github.io/cardano-wallet/api/edge/swagger.yaml) of `cardano-wallet`, which is documented [here](https://input-output-hk.github.io/cardano-wallet/api/edge/).

The [wallet package](wallet/) contains the generated client library, along with a few convenience functions and tests.
[The cmd/godano-wallet-cli package](cmd/godano-wallet-cli/) is thin CLI wrapper for the client library.

# Using the client library

A client can be conveniently created using any of the `New[HTTPS]Client[WithResponses]` methods.
See the `Client` or `ClientWithResponses` interfaces for a full list of supported operations, which mirror the [cardano-wallet REST API](https://input-output-hk.github.io/cardano-wallet/api/edge/).

```
addr := "https://localhost:12345/v2"
client, err := wallet.NewHTTPSClientWithResponses(addr, wallet.MakeTLSConfig())
ctx := context.Background()
settings, err := client.GetSettingsWithResponse(ctx)
pools, err := client.ListStakePoolsWithResponse(ctx, &ListStakePoolsParams {
	Stake: 1000
})
```

The following environment variables control the connection to the `cardano-wallet` server.
The `wallet.MakeTLSConfig()` method creates a TLS configuration for communication with the `cardano-wallet` process started by the Daedalus wallet.
Other instances of `cardano-wallet` might require different parameters.
Using `wallet.MakeTLSConfig()`, or HTTPS for that matter, is optional.
`wallet.MakeTLSConfig()` uses the following environment variables:

```
DAEDALUS_DIR="$HOME/.local/share/Daedalus/mainnet"
export GODANO_WALLET_CLIENT_TLS_SKIP_VERIFY="false" # Can be set to true instead of providing GODANO_WALLET_CLIENT_SERVER_CA
export GODANO_WALLET_CLIENT_SERVER_CA="$DAEDALUS_DIR/tls/server/ca.crt"
export GODANO_WALLET_CLIENT_CLIENT_CERT="$DAEDALUS_DIR/tls/client/client.crt"
export GODANO_WALLET_CLIENT_CLIENT_KEY="$DAEDALUS_DIR/tls/client/client.key"
```

# Using the CLI

Run the executable for a list of available commands. The commands mirror CRUD operations of the [`cardano-wallet` REST API](https://input-output-hk.github.io/cardano-wallet/api/edge/).

```
$ go run ./cmd/godano-wallet-client
godano-wallet-cli connects to the REST API of a cardano-wallet process and
translates the CLI commands and parameters to appropriate REST API calls

Usage:
  godano-wallet-cli [command]

Available Commands:
  AccountKey               post AccountKey objects
  Address                  create, import, inspect, list, or post Address objects
  AddressBatch             import AddressBatch objects
  Asset                    get or list Asset objects
  AssetDefault             get AssetDefault objects
  Byron                    Commands for Byron-era objects
  Coins                    select Coins objects
  CurrentSmashHealth       get CurrentSmashHealth objects
  DelegationFee            get DelegationFee objects
  MaintenanceAction        get or post MaintenanceAction objects
  Metadata                 sign Metadata objects
  NetworkClock             get NetworkClock objects
  NetworkInformation       get NetworkInformation objects
  NetworkParameters        get NetworkParameters objects
  Settings                 get or put Settings objects
  SharedWallet             delete, get, or post SharedWallet objects
  SharedWalletInDelegation patch SharedWalletInDelegation objects
  SharedWalletInPayment    patch SharedWalletInPayment objects
  StakePool                join, list, or quit StakePool objects
  Transaction              delete, get, list, or post Transaction objects
  TransactionFee           post TransactionFee objects
  UTxOsStatistics          get UTxOsStatistics objects
  Wallet                   delete, get, list, migrate, post, or put Wallet objects
  WalletKey                get WalletKey objects
  WalletMigrationInfo      get WalletMigrationInfo objects
  WalletPassphrase         put WalletPassphrase objects
  help                     Help about any command

Flags:
  -n, --dry-run         Show the resulting request instead of executing it
  -h, --help            help for godano-wallet-cli
  -q, --quiet           Set the log level to Warning
  -Q, --quieter         Set the log level to Error
  -s, --server string   Endpoint of the cardano-wallet process to connect to
  -V, --trace           Set the log level to Trace
  -v, --verbose         Set the log level to Debug
  -y, --yaml            Output responses as YAML instead of JSON (more compact)

Use "godano-wallet-cli [command] --help" for more information about a command.
```

The environment variable `GODANO_WALLET_CLIENT_SERVER_ADDRESS` is the default server URL to connect to, which can be overwritten by the `-s` flag. The tests (see below), also use this environment variable.

The environment variable `GODANO_WALLET_CLIENT_VERBOSE` can be set to a non-empty value to enable early debug-level logging in the CLI.
This will show how the CLI analyses methods in the `wallet.Client` interface for dynamically generating commands and sub-commands.

# Updating the generated code

The `generate.sh` script updates the generated code:
```
./generate.sh
```

The script uses `oapi-codegen`, `goimports`, and `gofumpt`, as well as `mkdir`, `wget`, and `sed`. See the comments in the script on how to install these requirements.
The script updates `swagger.yaml` and files named `wallet/generated-*.go`. If for example the package name is changed to something different than `wallet`, the other files in [the wallet package](wallet) must be updated manually.

Due to a (presumed) bug in `oapi-codegen`, the generated code is patched using `sed` after generation. Two missing structs (`Metadata` and `Distributions`) are implemented in `wallet/types-*.go`. Both these structs use integers as keys in JSON response objects, which is a rare use case and could be the reason for the errorenous code generation.

# Testing

Run

```
go test ./wallet
```

to perform some basic request tests.
The tests connect to the `cardano-wallet` process as configured through the environment variables above.
The tests aim to cover all simple and non-destructive read operations (`Get*`, `List*`, `Inspect*`).

Add the `-v` switch to see the received and parsed responses from each request:

```
go test -v ./wallet
```
