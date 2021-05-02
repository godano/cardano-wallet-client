package wallet

var ArgumentNames = map[string][]string{
	"PostAnyAddress":                {"body"},
	"InspectAddress":                {"addressId"},
	"ListByronWallets":              {},
	"PostByronWallet":               {"body"},
	"DeleteByronWallet":             {"walletId"},
	"GetByronWallet":                {"walletId"},
	"PutByronWallet":                {"walletId", "body"},
	"ListByronAddresses":            {"walletId", "params"},
	"CreateAddress":                 {"walletId", "body"},
	"ImportAddresses":               {"walletId", "body"},
	"ImportAddress":                 {"walletId", "addressId"},
	"ListByronAssets":               {"walletId"},
	"GetByronAssetDefault":          {"walletId", "policyId"},
	"GetByronAsset":                 {"walletId", "policyId", "assetName"},
	"ByronSelectCoins":              {"walletId", "body"},
	"GetByronWalletMigrationInfo":   {"walletId"},
	"MigrateByronWallet":            {"walletId", "body"},
	"PutByronWalletPassphrase":      {"walletId", "body"},
	"PostByronTransactionFee":       {"walletId", "body"},
	"GetByronUTxOsStatistics":       {"walletId"},
	"ListByronTransactions":         {"walletId", "params"},
	"PostByronTransaction":          {"walletId", "body"},
	"DeleteByronTransaction":        {"walletId", "transactionId"},
	"GetByronTransaction":           {"walletId", "transactionId"},
	"GetNetworkClock":               {"params"},
	"GetNetworkInformation":         {},
	"GetNetworkParameters":          {},
	"PostExternalTransaction":       {"body"},
	"GetSettings":                   {},
	"PutSettings":                   {"body"},
	"PostSharedWallet":              {"body"},
	"DeleteSharedWallet":            {"walletId"},
	"GetSharedWallet":               {"walletId"},
	"PatchSharedWalletInDelegation": {"walletId", "body"},
	"PatchSharedWalletInPayment":    {"walletId", "body"},
	"GetCurrentSmashHealth":         {"params"},
	"ListStakePools":                {"params"},
	"QuitStakePool":                 {"walletId", "body"},
	"GetMaintenanceActions":         {},
	"PostMaintenanceAction":         {"body"},
	"JoinStakePool":                 {"stakePoolId", "walletId", "body"},
	"ListWallets":                   {},
	"PostWallet":                    {"body"},
	"DeleteWallet":                  {"walletId"},
	"GetWallet":                     {"walletId"},
	"PutWallet":                     {"walletId", "body"},
	"ListAddresses":                 {"walletId", "params"},
	"ListAssets":                    {"walletId"},
	"GetAssetDefault":               {"walletId", "policyId"},
	"GetAsset":                      {"walletId", "policyId", "assetName"},
	"SelectCoins":                   {"walletId", "body"},
	"GetDelegationFee":              {"walletId"},
	"PostAccountKey":                {"walletId", "index", "body"},
	"GetWalletKey":                  {"walletId", "role", "index"},
	"GetShelleyWalletMigrationInfo": {"walletId"},
	"MigrateShelleyWallet":          {"walletId", "body"},
	"PutWalletPassphrase":           {"walletId", "body"},
	"PostTransactionFee":            {"walletId", "body"},
	"SignMetadata":                  {"walletId", "role", "index", "body"},
	"GetUTxOsStatistics":            {"walletId"},
	"ListTransactions":              {"walletId", "params"},
	"PostTransaction":               {"walletId", "body"},
	"DeleteTransaction":             {"walletId", "transactionId"},
	"GetTransaction":                {"walletId", "transactionId"},
}

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
	"SelectCoins":                   true, // ab hier ok
	"PostAccountKey":                true,
	"MigrateShelleyWallet":          true,
	"PutWalletPassphrase":           true,
	"PostTransactionFee":            true,
	"SignMetadata":                  true,
	"PostTransaction":               true,
}

var ArgumentNamesWithResponse = make(map[string][]string, len(ArgumentNames))

func init() {
	// Add *WithBody methods
	for name := range MethodHasBody {
		src := ArgumentNames[name]
		params := make([]string, 0, len(src))
		params = append(params, src[:len(src)-1]...)
		params = append(params, "contentType", "body") // These are always fixed
		ArgumentNames[name+"WithBody"] = params
	}

	// This method does not have the version without *WithBody
	delete(ArgumentNames, "PostExternalTransaction")

	// Add *WithResponse methods from the ClientWithResponsesInterface
	for name, params := range ArgumentNames {
		ArgumentNamesWithResponse[name+"WithResponse"] = params
	}
}
