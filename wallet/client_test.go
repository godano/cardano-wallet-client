package wallet

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type testSuiteBase struct {
	suite.Suite
	*require.Assertions

	client ClientWithResponsesInterface
	ctx    context.Context
}

func (s *testSuiteBase) resp(err error, expectedCode int, resp *http.Response, body []byte) {
	s.NoError(err)
	s.NotNil(resp)
	if expectedCode != resp.StatusCode {
		s.Fail("Unexpected http code", "Received code %v, expected %v.\nResponse body: %s",
			resp.StatusCode, expectedCode, body)
	}
}

func (s *testSuiteBase) matches(valueName string, value string, regexStr string) {
	regex, err := regexp.Compile(regexStr)
	s.NoError(err)
	s.True(regex.MatchString(value), "Unexpected %s: '%v', does not match regex: %v", valueName, value, regexStr)
}

func (s *testSuiteBase) logObject(name string, obj interface{}) {
	marshalled, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		log.Errorf("Failed to JSON-marshal response object: %v", err)
		return
	}
	scanner := bufio.NewScanner(bytes.NewReader(marshalled))
	scanner.Split(bufio.ScanLines)
	log.Infof("%v:", name)
	for scanner.Scan() {
		log.Info(scanner.Text())
	}
}

type CardanoWalletTestSuite struct {
	testSuiteBase
}

func TestCardanoWalletClient(t *testing.T) {
	testSuite := new(CardanoWalletTestSuite)
	testSuite.ctx = context.Background()
	suite.Run(t, testSuite)
}

func (s *CardanoWalletTestSuite) SetupSuite() {
	s.testSuiteBase.Assertions = s.Require()

	var err error
	s.client, err = NewWalletClient()
	s.NoError(err)
	s.NotNil(s.client)
}

func (s *CardanoWalletTestSuite) TestSetup() {
	// Nothing - simply run the SetupSuite method
}

// === Get, Inspect, List requests
// [ ] InspectAddress(ctx context.Context, addressId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [x] ListByronWallets(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
// [x] GetByronWallet(ctx context.Context, walletId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [x] ListByronAddresses(ctx context.Context, walletId string, params *ListByronAddressesParams, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] ListByronAssets(ctx context.Context, walletId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] GetByronAssetDefault(ctx context.Context, walletId string, policyId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] GetByronAsset(ctx context.Context, walletId string, policyId string, assetName string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] GetByronWalletMigrationInfo(ctx context.Context, walletId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] GetByronUTxOsStatistics(ctx context.Context, walletId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] ListByronTransactions(ctx context.Context, walletId string, params *ListByronTransactionsParams, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] GetByronTransaction(ctx context.Context, walletId string, transactionId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [x] GetNetworkClock(ctx context.Context, params *GetNetworkClockParams, reqEditors ...RequestEditorFn) (*http.Response, error)
// [x] GetNetworkInformation(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
// [x] GetNetworkParameters(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
// [x] GetSettings(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] GetSharedWallet(ctx context.Context, walletId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [x] GetCurrentSmashHealth(ctx context.Context, params *GetCurrentSmashHealthParams, reqEditors ...RequestEditorFn) (*http.Response, error)
// [x] ListStakePools(ctx context.Context, params *ListStakePoolsParams, reqEditors ...RequestEditorFn) (*http.Response, error)
// [x] GetMaintenanceActions(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] ListWallets(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] GetWallet(ctx context.Context, walletId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] ListAddresses(ctx context.Context, walletId string, params *ListAddressesParams, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] ListAssets(ctx context.Context, walletId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] GetAssetDefault(ctx context.Context, walletId string, policyId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] GetAsset(ctx context.Context, walletId string, policyId string, assetName string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] GetDelegationFee(ctx context.Context, walletId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] GetWalletKey(ctx context.Context, walletId string, role string, index string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] GetShelleyWalletMigrationInfo(ctx context.Context, walletId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] GetUTxOsStatistics(ctx context.Context, walletId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] ListTransactions(ctx context.Context, walletId string, params *ListTransactionsParams, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] GetTransaction(ctx context.Context, walletId string, transactionId string, reqEditors ...RequestEditorFn) (*http.Response, error)

// === Create, Post, Put, Import, Patch requests
// PostAnyAddressWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
// PostAnyAddress(ctx context.Context, body PostAnyAddressJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
// PutByronWalletWithBody(ctx context.Context, walletId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
// PutByronWallet(ctx context.Context, walletId string, body PutByronWalletJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
// PostByronWalletWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
// PostByronWallet(ctx context.Context, body PostByronWalletJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
// CreateAddressWithBody(ctx context.Context, walletId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
// CreateAddress(ctx context.Context, walletId string, body CreateAddressJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
// ImportAddressesWithBody(ctx context.Context, walletId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
// ImportAddresses(ctx context.Context, walletId string, body ImportAddressesJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
// ImportAddress(ctx context.Context, walletId string, addressId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// PutByronWalletPassphraseWithBody(ctx context.Context, walletId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
// PutByronWalletPassphrase(ctx context.Context, walletId string, body PutByronWalletPassphraseJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
// PostByronTransactionFeeWithBody(ctx context.Context, walletId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
// PostByronTransactionFee(ctx context.Context, walletId string, body PostByronTransactionFeeJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
// PostByronTransactionWithBody(ctx context.Context, walletId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
// PostByronTransaction(ctx context.Context, walletId string, body PostByronTransactionJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
// PostExternalTransactionWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
// PutSettingsWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
// PutSettings(ctx context.Context, body PutSettingsJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
// PostSharedWalletWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
// PostSharedWallet(ctx context.Context, body PostSharedWalletJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
// PatchSharedWalletInDelegationWithBody(ctx context.Context, walletId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
// PatchSharedWalletInDelegation(ctx context.Context, walletId string, body PatchSharedWalletInDelegationJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
// PatchSharedWalletInPaymentWithBody(ctx context.Context, walletId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
// PatchSharedWalletInPayment(ctx context.Context, walletId string, body PatchSharedWalletInPaymentJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
// PostMaintenanceActionWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
// PostMaintenanceAction(ctx context.Context, body PostMaintenanceActionJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
// PostWalletWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
// PostWallet(ctx context.Context, body PostWalletJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
// PutWalletWithBody(ctx context.Context, walletId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
// PutWallet(ctx context.Context, walletId string, body PutWalletJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
// PostAccountKeyWithBody(ctx context.Context, walletId string, index string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
// PostAccountKey(ctx context.Context, walletId string, index string, body PostAccountKeyJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
// PutWalletPassphraseWithBody(ctx context.Context, walletId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
// PutWalletPassphrase(ctx context.Context, walletId string, body PutWalletPassphraseJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
// PostTransactionFeeWithBody(ctx context.Context, walletId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
// PostTransactionFee(ctx context.Context, walletId string, body PostTransactionFeeJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
// PostTransactionWithBody(ctx context.Context, walletId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
// PostTransaction(ctx context.Context, walletId string, body PostTransactionJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

// === Delete requests
// DeleteByronWallet(ctx context.Context, walletId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// DeleteByronTransaction(ctx context.Context, walletId string, transactionId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// DeleteSharedWallet(ctx context.Context, walletId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// DeleteWallet(ctx context.Context, walletId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// DeleteTransaction(ctx context.Context, walletId string, transactionId string, reqEditors ...RequestEditorFn) (*http.Response, error)

// === Select, Migrate, Quit, Join, Sign requests
// ByronSelectCoinsWithBody(ctx context.Context, walletId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
// ByronSelectCoins(ctx context.Context, walletId string, body ByronSelectCoinsJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
// MigrateByronWalletWithBody(ctx context.Context, walletId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
// MigrateByronWallet(ctx context.Context, walletId string, body MigrateByronWalletJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
// QuitStakePoolWithBody(ctx context.Context, walletId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
// QuitStakePool(ctx context.Context, walletId string, body QuitStakePoolJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
// JoinStakePoolWithBody(ctx context.Context, stakePoolId string, walletId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
// JoinStakePool(ctx context.Context, stakePoolId string, walletId string, body JoinStakePoolJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
// SelectCoinsWithBody(ctx context.Context, walletId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
// SelectCoins(ctx context.Context, walletId string, body SelectCoinsJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
// MigrateShelleyWalletWithBody(ctx context.Context, walletId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
// MigrateShelleyWallet(ctx context.Context, walletId string, body MigrateShelleyWalletJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
// SignMetadataWithBody(ctx context.Context, walletId string, role string, index string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
// SignMetadata(ctx context.Context, walletId string, role string, index string, body SignMetadataJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

func (s *CardanoWalletTestSuite) TestGetSettings() {
	resp, err := s.client.GetSettingsWithResponse(s.ctx)
	s.NoError(err)
	s.resp(err, http.StatusOK, resp.HTTPResponse, resp.Body)
	s.NotNil(resp.JSON200)
	s.logObject("GetSettings", resp.JSON200)

	// Regex taken from Swagger doc
	s.matches("PoolMetadataSource", resp.JSON200.PoolMetadataSource,
		"^(none|direct|https?://[a-zA-Z0-9-_~.]+(:[0-9]+)?/?)$")
}

func (s *CardanoWalletTestSuite) TestGetNetworkParameters() {
	resp, err := s.client.GetNetworkParametersWithResponse(s.ctx)
	s.NoError(err)
	s.resp(err, http.StatusOK, resp.HTTPResponse, resp.Body)
	s.NotNil(resp.JSON200)
	s.Nil(resp.JSON406)
	s.logObject("GetNetworkParameters", resp.JSON200)
}

func (s *CardanoWalletTestSuite) TestGetNetworkInformation() {
	resp, err := s.client.GetNetworkInformationWithResponse(s.ctx)
	s.NoError(err)
	s.resp(err, http.StatusOK, resp.HTTPResponse, resp.Body)
	s.NotNil(resp.JSON200)
	s.Nil(resp.JSON406)
	s.logObject("GetNetworkInformation", resp.JSON200)
}

func (s *CardanoWalletTestSuite) TestGetNetworkClock() {
	resp, err := s.client.GetNetworkClockWithResponse(s.ctx, new(GetNetworkClockParams))
	s.NoError(err)
	s.resp(err, http.StatusOK, resp.HTTPResponse, resp.Body)
	s.NotNil(resp.JSON200)
	s.Nil(resp.JSON406)
	s.logObject("GetNetworkClock", resp.JSON200)
}

func (s *CardanoWalletTestSuite) TestGetCurrentSmashHealth() {
	resp, err := s.client.GetCurrentSmashHealthWithResponse(s.ctx, new(GetCurrentSmashHealthParams))
	s.NoError(err)
	s.resp(err, http.StatusOK, resp.HTTPResponse, resp.Body)
	s.NotNil(resp.JSON200)
	s.Nil(resp.JSON406)
	s.Nil(resp.JSON400)
	s.logObject("GetCurrentSmashHealth", resp.JSON200)
}

// TODO this currently fails with "unexpected EOF"
// Probably too long-running?
// func (s *CardanoWalletTestSuite) TestListStakePools() {
// 	resp, err := s.client.ListStakePoolsWithResponse(s.ctx,
// 		&ListStakePoolsParams{Stake: 100 * 1000 * 1000}) // 100 Ada
// 	s.NoError(err)
// 	s.resp(err, http.StatusOK, resp.HTTPResponse, resp.Body)
// 	s.NotNil(resp.JSON200)
// 	s.Nil(resp.JSON400)
// 	s.logObject("ListStakePools", resp.JSON200)
// }

func (s *CardanoWalletTestSuite) TestGetMaintenanceActions() {
	resp, err := s.client.GetMaintenanceActionsWithResponse(s.ctx)
	s.NoError(err)
	s.resp(err, http.StatusOK, resp.HTTPResponse, resp.Body)
	s.NotNil(resp.JSON200)
	s.logObject("ListStakePools", resp.JSON200)
}

func (s *CardanoWalletTestSuite) TestListByronWallets() {
	resp, err := s.client.ListByronWalletsWithResponse(s.ctx)
	s.NoError(err)
	s.resp(err, http.StatusOK, resp.HTTPResponse, resp.Body)
	s.NotNil(resp.JSON200)
	s.Nil(resp.JSON406)
	s.logObject("ListByronWallets", resp.JSON200)

	wallets := *resp.JSON200
	if len(wallets) == 0 {
		log.Warnf("Cannot test ByronWallet requests: no Byron wallets available")
		return
	}

	// Further test Byron wallets, if possible
	walletId := wallets[0].Id
	subSuite := &ByronWalletTestSuite{
		CardanoWalletTestSuite: s,
		walletId:               walletId,
	}
	s.Run(subSuite.String(), func() {
		suite.Run(s.T(), subSuite)
	})
}

type ByronWalletTestSuite struct {
	*CardanoWalletTestSuite
	walletId string
}

func (s *ByronWalletTestSuite) String() string {
	return fmt.Sprintf("ByronWalletTestSuite(%s)", s.walletId)
}

func (s *ByronWalletTestSuite) TestGetByronWallet() {
	resp, err := s.client.GetByronWalletWithResponse(s.ctx, s.walletId)
	s.NoError(err)
	s.resp(err, http.StatusOK, resp.HTTPResponse, resp.Body)
	s.NotNil(resp.JSON200)
	s.Nil(resp.JSON404)
	s.Nil(resp.JSON406)
	s.logObject("GetByronWallet", resp.JSON200)
}

func (s *ByronWalletTestSuite) TestListByronAddresses() {
	resp, err := s.client.ListByronAddressesWithResponse(s.ctx, s.walletId, new(ListByronAddressesParams))
	s.NoError(err)
	s.resp(err, http.StatusOK, resp.HTTPResponse, resp.Body)
	s.NotNil(resp.JSON200)
	s.Nil(resp.JSON400)
	s.Nil(resp.JSON404)
	s.Nil(resp.JSON406)
	s.logObject("ListByronAddresses", resp.JSON200)
}

// TODO ByronWallet related requests
// [ ] GetByronUTxOsStatistics(ctx context.Context, walletId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] ListByronAssets(ctx context.Context, walletId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] GetByronWalletMigrationInfo(ctx context.Context, walletId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] ListByronTransactions(ctx context.Context, walletId string, params *ListByronTransactionsParams, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] GetByronAssetDefault(ctx context.Context, walletId string, policyId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] GetByronAsset(ctx context.Context, walletId string, policyId string, assetName string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] GetByronTransaction(ctx context.Context, walletId string, transactionId string, reqEd

func (s *CardanoWalletTestSuite) TestListWallets() {
	resp, err := s.client.ListWalletsWithResponse(s.ctx)
	s.NoError(err)
	s.resp(err, http.StatusOK, resp.HTTPResponse, resp.Body)
	s.NotNil(resp.JSON200)
	s.Nil(resp.JSON406)
	s.logObject("ListWallets", resp.JSON200)

	wallets := *resp.JSON200
	if len(wallets) == 0 {
		log.Warnf("Cannot test Wallet requests: no wallets available")
		return
	}

	// Further test wallets, if possible
	walletId := wallets[0].Id
	subSuite := &SingleWalletTestSuite{
		testSuiteBase: s.testSuiteBase,
		walletId:      walletId,
	}
	s.Run(subSuite.String(), func() {
		suite.Run(s.T(), subSuite)
	})
}

type SingleWalletTestSuite struct {
	testSuiteBase
	walletId string
}

func (s *SingleWalletTestSuite) String() string {
	return fmt.Sprintf("SingleWalletTestSuite(%s)", s.walletId)
}

func (s *SingleWalletTestSuite) TestGetWallet() {
	resp, err := s.client.GetWalletWithResponse(s.ctx, s.walletId)
	s.NoError(err)
	s.resp(err, http.StatusOK, resp.HTTPResponse, resp.Body)
	s.NotNil(resp.JSON200)
	s.Nil(resp.JSON404)
	s.Nil(resp.JSON406)
	s.logObject("GetWallet", resp.JSON200)
}

func (s *SingleWalletTestSuite) TestListAddresses() {
	resp, err := s.client.ListAddressesWithResponse(s.ctx, s.walletId, new(ListAddressesParams))
	s.NoError(err)
	s.resp(err, http.StatusOK, resp.HTTPResponse, resp.Body)
	s.NotNil(resp.JSON200)
	s.Nil(resp.JSON400)
	s.Nil(resp.JSON404)
	s.Nil(resp.JSON406)
	s.logObject("ListAddresses", resp.JSON200)
}

// TODO test wallet requests
// [ ] ListAssets(ctx context.Context, walletId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] GetAssetDefault(ctx context.Context, walletId string, policyId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] GetAsset(ctx context.Context, walletId string, policyId string, assetName string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] GetDelegationFee(ctx context.Context, walletId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] GetWalletKey(ctx context.Context, walletId string, role string, index string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] GetShelleyWalletMigrationInfo(ctx context.Context, walletId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] GetUTxOsStatistics(ctx context.Context, walletId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] ListTransactions(ctx context.Context, walletId string, params *ListTransactionsParams, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] GetTransaction(ctx context.Context, walletId string, transactionId string, reqEditors ...RequestEditorFn) (*http.Response, error)
