package wallet

// MakeArgument returns a value of the non-string argument of the given method.
// For methods that do not have a non-string argument, or for unknown methods, MakeArgument returns nil.
// All generated client methods have at most one non-string argument, which can be a struct or a pointer to a struct.
// TODO this could and should be auto-generated from swagger.yaml or based on the source code of ClientInterface.
func MakeArgument(method string) interface{} {
	switch method {
	case "PostAnyAddress":
		return new(PostAnyAddressJSONRequestBody)
	case "PostByronWallet":
		return PostByronWalletJSONRequestBody(nilValue())
	case "PutByronWallet":
		return new(PutByronWalletJSONRequestBody)
	case "ListByronAddresses":
		return new(ListByronAddressesParams)
	case "CreateAddress":
		return new(CreateAddressJSONRequestBody)
	case "ImportAddresses":
		return new(ImportAddressesJSONRequestBody)
	case "ByronSelectCoins":
		return new(ByronSelectCoinsJSONRequestBody)
	case "MigrateByronWallet":
		return new(MigrateByronWalletJSONRequestBody)
	case "PutByronWalletPassphrase":
		return new(PutByronWalletPassphraseJSONRequestBody)
	case "PostByronTransactionFee":
		return new(PostByronTransactionFeeJSONRequestBody)
	case "ListByronTransactions":
		return new(ListByronTransactionsParams)
	case "PostByronTransaction":
		return new(PostByronTransactionJSONRequestBody)
	case "GetNetworkClock":
		return new(GetNetworkClockParams)
	case "PostExternalTransaction":
		return nil // TODO no specific body struct
	case "PutSettings":
		return new(PutSettingsJSONRequestBody)
	case "PostSharedWallet":
		return PostSharedWalletJSONRequestBody(nilValue())
	case "PatchSharedWalletInDelegation":
		return new(PatchSharedWalletInDelegationJSONRequestBody)
	case "PatchSharedWalletInPayment":
		return new(PatchSharedWalletInPaymentJSONRequestBody)
	case "GetCurrentSmashHealth":
		return new(GetCurrentSmashHealthParams)
	case "ListStakePools":
		return new(ListStakePoolsParams)
	case "QuitStakePool":
		return new(QuitStakePoolJSONRequestBody)
	case "PostMaintenanceAction":
		return new(PostMaintenanceActionJSONRequestBody)
	case "JoinStakePool":
		return new(JoinStakePoolJSONRequestBody)
	case "PostWallet":
		return PostWalletJSONRequestBody(nilValue())
	case "PutWallet":
		return new(PutWalletJSONRequestBody)
	case "ListAddresses":
		return new(ListAddressesParams)
	case "SelectCoins":
		return SelectCoinsJSONRequestBody(nilValue())
	case "PostAccountKey":
		return new(PostAccountKeyJSONRequestBody)
	case "MigrateShelleyWallet":
		return new(MigrateShelleyWalletJSONRequestBody)
	case "PutWalletPassphrase":
		return new(PutWalletPassphraseJSONRequestBody)
	case "PostTransactionFee":
		return PostTransactionFeeJSONRequestBody(nilValue())
	case "SignMetadata":
		return new(SignMetadataJSONRequestBody)
	case "ListTransactions":
		return new(ListTransactionsParams)
	case "PostTransaction":
		return PostTransactionJSONRequestBody(nilValue())
	}
	return nil
}

func nilValue() interface{} {
	// Some body structs are generated as type-aliases to interface{}
	// In these cases, use an empty map, to enable formatting and parsing with arbitrary JSON.
	m := make(map[string]interface{})
	return &m
}

// MethodHasParamsStruct contains the names of all methods that include a `params` struct
// which is not passed in the HTTP request body. This `params` struct and all contained values
// are optional for the request (i.e. the method argument can be set to nil).
var MethodHasParamsStruct = map[string]bool{
	"ListByronAddresses":    true,
	"ListByronTransactions": true,
	"GetNetworkClock":       true,
	"GetCurrentSmashHealth": true,
	"ListStakePools":        true,
	"ListAddresses":         true,
	"ListTransactions":      true,
}
