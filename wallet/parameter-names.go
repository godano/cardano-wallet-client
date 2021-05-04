package wallet

const (
	ParamsArgName = "params"
	BodyArgName   = "body"
)

// ArgumentNames maps method names from the `ClientInterface` to their parameter names.
// The receiver and the first parameter (ctx) are excluded.
//
// TODO this could and should be auto-generated from swagger.yaml or based on the source code of
// ClientInterface. The problem is that Go drops parameter names during compilation.
var ArgumentNames = map[string][]string{
	"PostAnyAddress":                {BodyArgName},
	"InspectAddress":                {"addressId"},
	"ListByronWallets":              {},
	"PostByronWallet":               {BodyArgName},
	"DeleteByronWallet":             {"walletId"},
	"GetByronWallet":                {"walletId"},
	"PutByronWallet":                {"walletId", BodyArgName},
	"ListByronAddresses":            {"walletId", ParamsArgName},
	"CreateAddress":                 {"walletId", BodyArgName},
	"ImportAddresses":               {"walletId", BodyArgName},
	"ImportAddress":                 {"walletId", "addressId"},
	"ListByronAssets":               {"walletId"},
	"GetByronAssetDefault":          {"walletId", "policyId"},
	"GetByronAsset":                 {"walletId", "policyId", "assetName"},
	"ByronSelectCoins":              {"walletId", BodyArgName},
	"GetByronWalletMigrationInfo":   {"walletId"},
	"MigrateByronWallet":            {"walletId", BodyArgName},
	"PutByronWalletPassphrase":      {"walletId", BodyArgName},
	"PostByronTransactionFee":       {"walletId", BodyArgName},
	"GetByronUTxOsStatistics":       {"walletId"},
	"ListByronTransactions":         {"walletId", ParamsArgName},
	"PostByronTransaction":          {"walletId", BodyArgName},
	"DeleteByronTransaction":        {"walletId", "transactionId"},
	"GetByronTransaction":           {"walletId", "transactionId"},
	"GetNetworkClock":               {ParamsArgName},
	"GetNetworkInformation":         {},
	"GetNetworkParameters":          {},
	"PostExternalTransaction":       {BodyArgName},
	"GetSettings":                   {},
	"PutSettings":                   {BodyArgName},
	"PostSharedWallet":              {BodyArgName},
	"DeleteSharedWallet":            {"walletId"},
	"GetSharedWallet":               {"walletId"},
	"PatchSharedWalletInDelegation": {"walletId", BodyArgName},
	"PatchSharedWalletInPayment":    {"walletId", BodyArgName},
	"GetCurrentSmashHealth":         {ParamsArgName},
	"ListStakePools":                {ParamsArgName},
	"QuitStakePool":                 {"walletId", BodyArgName},
	"GetMaintenanceActions":         {},
	"PostMaintenanceAction":         {BodyArgName},
	"JoinStakePool":                 {"stakePoolId", "walletId", BodyArgName},
	"ListWallets":                   {},
	"PostWallet":                    {BodyArgName},
	"DeleteWallet":                  {"walletId"},
	"GetWallet":                     {"walletId"},
	"PutWallet":                     {"walletId", BodyArgName},
	"ListAddresses":                 {"walletId", ParamsArgName},
	"ListAssets":                    {"walletId"},
	"GetAssetDefault":               {"walletId", "policyId"},
	"GetAsset":                      {"walletId", "policyId", "assetName"},
	"SelectCoins":                   {"walletId", BodyArgName},
	"GetDelegationFee":              {"walletId"},
	"PostAccountKey":                {"walletId", "index", BodyArgName},
	"GetWalletKey":                  {"walletId", "role", "index"},
	"GetShelleyWalletMigrationInfo": {"walletId"},
	"MigrateShelleyWallet":          {"walletId", BodyArgName},
	"PutWalletPassphrase":           {"walletId", BodyArgName},
	"PostTransactionFee":            {"walletId", BodyArgName},
	"SignMetadata":                  {"walletId", "role", "index", BodyArgName},
	"GetUTxOsStatistics":            {"walletId"},
	"ListTransactions":              {"walletId", ParamsArgName},
	"PostTransaction":               {"walletId", BodyArgName},
	"DeleteTransaction":             {"walletId", "transactionId"},
	"GetTransaction":                {"walletId", "transactionId"},
}

// MethodHasBody contains the names of all methods that have a *WithBody variant (e.g. PostAnyAddress -> PostAnyAddressWithBody)
// These methods have the same parameters as listed above, except that the body content is given as a
// `contentType string` and `body io.Reader`, instead of a method-specific struct.
var MethodHasBody = map[string]bool{
	"PostAnyAddress":                true,
	"PostByronWallet":               true,
	"PutByronWallet":                true,
	"CreateAddress":                 true,
	"ImportAddresses":               true,
	"ByronSelectCoins":              true,
	"MigrateByronWallet":            true,
	"PutByronWalletPassphrase":      true,
	"PostByronTransactionFee":       true,
	"PostByronTransaction":          true,
	"PutSettings":                   true,
	"PostExternalTransaction":       true,
	"PostSharedWallet":              true,
	"PatchSharedWalletInDelegation": true,
	"PatchSharedWalletInPayment":    true,
	"QuitStakePool":                 true,
	"PostMaintenanceAction":         true,
	"JoinStakePool":                 true,
	"PostWallet":                    true,
	"PutWallet":                     true,
	"SelectCoins":                   true,
	"PostAccountKey":                true,
	"MigrateShelleyWallet":          true,
	"PutWalletPassphrase":           true,
	"PostTransactionFee":            true,
	"SignMetadata":                  true,
	"PostTransaction":               true,
}

// ArgumentNamesWithResponse is a copy of ArgumentNames, with each method suffixed by `WithResponse`.
// These methods and parameter names correspond to the `ClientWithResponse` Interface.
var ArgumentNamesWithResponse = make(map[string][]string, len(ArgumentNames))

func init() {
	// Add *WithBody methods
	for name := range MethodHasBody {
		src := ArgumentNames[name]
		params := make([]string, 0, len(src))
		params = append(params, src[:len(src)-1]...)
		params = append(params, "contentType", BodyArgName) // These are always fixed
		ArgumentNames[name+"WithBody"] = params
	}

	// Exception: this method does not have the version without *WithBody
	delete(ArgumentNames, "PostExternalTransaction")

	// Add *WithResponse methods from the ClientWithResponsesInterface
	for name, params := range ArgumentNames {
		ArgumentNamesWithResponse[name+"WithResponse"] = params
	}
}
