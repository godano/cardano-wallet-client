package wallet

// MakeArgument returns a value of the non-string argument of the given method.
// For methods that do not have a non-string argument, or for unknown methods, MakeArgument returns nil.
// All generated client methods have at most one non-string argument, which can be a struct or a pointer to a struct.
// TODO this could and should be auto-generated from swagger.yaml or based on the source code of ClientInterface.
func MakeArgument(method string) interface{} {
	switch method {
	case "PostAnyAddress":
		return PostAnyAddressJSONRequestBody{}
	case "PostByronWallet":
		return PostByronWalletJSONRequestBody(nil)
	case "PutByronWallet":
		return PutByronWalletJSONRequestBody{}
	case "ListByronAddresses":
		return new(ListByronAddressesParams)
	case "CreateAddress":
		return CreateAddressJSONRequestBody{}
	case "ImportAddresses":
		return ImportAddressesJSONRequestBody{}
	case "ByronSelectCoins":
		return ByronSelectCoinsJSONRequestBody{}
	case "MigrateByronWallet":
		return MigrateByronWalletJSONRequestBody{}
	case "PutByronWalletPassphrase":
		return PutByronWalletPassphraseJSONRequestBody{}
	case "PostByronTransactionFee":
		return PostByronTransactionFeeJSONRequestBody{}
	case "ListByronTransactions":
		return new(ListByronTransactionsParams)
	case "PostByronTransaction":
		return PostByronTransactionJSONRequestBody{}
	case "GetNetworkClock":
		return new(GetNetworkClockParams)
	case "PostExternalTransaction":
		return nil // TODO no specific body struct
	case "PutSettings":
		return PutSettingsJSONRequestBody{}
	case "PostSharedWallet":
		return PostSharedWalletJSONRequestBody(nil)
	case "PatchSharedWalletInDelegation":
		return PatchSharedWalletInDelegationJSONRequestBody{}
	case "PatchSharedWalletInPayment":
		return PatchSharedWalletInPaymentJSONRequestBody{}
	case "GetCurrentSmashHealth":
		return new(GetCurrentSmashHealthParams)
	case "ListStakePools":
		return new(ListStakePoolsParams)
	case "QuitStakePool":
		return QuitStakePoolJSONRequestBody{}
	case "PostMaintenanceAction":
		return PostMaintenanceActionJSONRequestBody{}
	case "JoinStakePool":
		return JoinStakePoolJSONRequestBody{}
	case "PostWallet":
		return PostWalletJSONRequestBody(nil)
	case "PutWallet":
		return PutWalletJSONRequestBody{}
	case "ListAddresses":
		return new(ListAddressesParams)
	case "SelectCoins":
		return SelectCoinsJSONRequestBody(nil)
	case "PostAccountKey":
		return PostAccountKeyJSONRequestBody{}
	case "MigrateShelleyWallet":
		return MigrateShelleyWalletJSONRequestBody{}
	case "PutWalletPassphrase":
		return PutWalletPassphraseJSONRequestBody{}
	case "PostTransactionFee":
		return PostTransactionFeeJSONRequestBody(nil)
	case "SignMetadata":
		return SignMetadataJSONRequestBody{}
	case "ListTransactions":
		return new(ListTransactionsParams)
	case "PostTransaction":
		return PostTransactionJSONRequestBody(nil)
	}
	return nil
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
