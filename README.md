<p align="center">
  <h1 align="center">
    Cardano-Wallet Client
    <br/>
    <a href="https://github.com/godano/cardano-wallet-client/blob/master/LICENSE" ><img alt="license" src="https://img.shields.io/badge/license-MIT%20License%202.0-E91E63.svg?style=flat-square" /></a>
    <a href="https://t.me/godano"><img src="https://img.shields.io/badge/Chat%20on-Telegram-blue.svg"/></a>
  </h1>
</p>

A Go client for the [cardano-wallet](https://github.com/input-output-hk/cardano-wallet) by IOG.

The bulk of this client code is generated using [oapi-codegen](https://github.com/deepmap/oapi-codegen), based on the [Open API definition](https://input-output-hk.github.io/cardano-wallet/api/edge/swagger.yaml) of `cardano-wallet`.

The [wallet package](wallet/) contains the generated client library, along with a few convenience functions and tests.
A `cmd/godano-wallet-cli` executable is planned as a thin wrapper for the client library.

# Using the client library

A client can be conveniently created using `NewWalletClient`. See the `ClientWithResponses` interface for a full list of supported operations.

```
client, err := wallet.NewWalletClient()
ctx := context.Background()
settings, err := client.GetSettingsWithResponse(ctx)
pools, err := client.ListStakePoolsWithResponse(ctx, &ListStakePoolsParams {
	Stake: 1000
})
```

The following environment variables control the connection to the `cardano-wallet` server.
The variables are designed to enable communication with the `cardano-wallet` process started by the Daedalus wallet. Other instances of `cardano-wallet` might require different parameters, see below.

```
DAEDALUS_DIR="$HOME/.local/share/Daedalus/mainnet"

export GODANO_WALLET_CLIENT_SERVER_ADDRESS=""
export GODANO_WALLET_CLIENT_TLS_SKIP_VERIFY="false" # Can be set to true instead of providing GODANO_WALLET_CLIENT_SERVER_CA
export GODANO_WALLET_CLIENT_SERVER_CA="$DAEDALUS_DIR/tls/server/ca.crt"
export GODANO_WALLET_CLIENT_CLIENT_CERT="$DAEDALUS_DIR/tls/client/client.crt"
export GODANO_WALLET_CLIENT_CLIENT_KEY="$DAEDALUS_DIR/tls/client/client.key"
```

Alternatively, the configuration based on environment variables can be skipped using `NewWalletClientFor`:
```
addr := "https://127.0.0.1:44107/v2"
conf := &tls.TLSConfig{/*...*/}
client, err := wallet.NewWalletClientFor(addr, conf)
```

For even more control over the HTTP client parameters, copy and edit the content of `NewWalletClientFor`.

# Updating the generated code

The `generate.sh` script updates the generated code:
```
./generate.sh
```

The script uses `oapi-codegen`, `goimports`, and `gofumpt`, as well as `mkdir`, `wget`, and `sed`.See the comments in the script on how to install these requirements.
The script only updates files named `wallet/generated-*.go`. If for example the package name is changed to something different than `wallet`, the other files must be updated manually.

Due to a (presumed) bug in `oapi-codegen`, the generated code is patched using `sed` after generation. Two missing structs (`Metadata` and `Distributions`) are implemented in `wallet/types-*.go`. Both these structs use integers as keys in JSON response objects, which is a rare use ase and could be the reason for the errorenous code generation.

# Testing

Run

```
go test ./wallet
```

to perform some basic connection tests. The tests connect to the `cardano-wallet` process as configured through the environment variables above. The tests aim to cover all simple and non-destructive read operations (`Get*`, `List*`, `Inspect*`).

Add the `-v` switch to see the received and parsed information from each request:

```
go test -v ./wallet
```
