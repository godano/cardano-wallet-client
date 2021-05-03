package wallet

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var log = logrus.New()

// To execute the tests, either set `testWalletServerAddress` to a valid cardano-wallet endpoint,
// or set the `GODANO_WALLET_CLIENT_SERVER_ADDRESS` env-var to that value.
const testWalletServerAddress = ""

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

// CardanoWalletTestSuite executed all read-only requests against a cardano-wallet endpoint.
// Some requests can only be executed with some pre-conditions, e.g. that a wallet or transaction is actually available.
// The long comment list below shows the covered and missing requests.
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

	address := testWalletServerAddress
	if envAddress := os.Getenv(EnvVarWalletServerAddress); envAddress != "" {
		address = envAddress
	}
	if address == "" {
		s.Fail("To execute tests, either define WalletServerAddress in client_test.go, or set the GODANO_WALLET_CLIENT_SERVER_ADDRESS environment variable")
	}

	tlsConfig, err := MakeTLSConfig()
	s.NoError(err)

	s.client, err = NewHTTPSClientWithResponses(address, tlsConfig)
	s.NoError(err)
	s.NotNil(s.client)
}

func (s *CardanoWalletTestSuite) TestSetup() {
	// Nothing - simply run the SetupSuite method
}

// === Get, Inspect, List requests
// [x] InspectAddress(ctx context.Context, addressId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [x] ListByronWallets(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
// [x] GetByronWallet(ctx context.Context, walletId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [x] ListByronAddresses(ctx context.Context, walletId string, params *ListByronAddressesParams, reqEditors ...RequestEditorFn) (*http.Response, error)
// [x] ListByronAssets(ctx context.Context, walletId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [x] GetByronAssetDefault(ctx context.Context, walletId string, policyId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [x] GetByronAsset(ctx context.Context, walletId string, policyId string, assetName string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [x] GetByronWalletMigrationInfo(ctx context.Context, walletId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [x] GetByronUTxOsStatistics(ctx context.Context, walletId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [x] ListByronTransactions(ctx context.Context, walletId string, params *ListByronTransactionsParams, reqEditors ...RequestEditorFn) (*http.Response, error)
// [x] GetByronTransaction(ctx context.Context, walletId string, transactionId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [x] GetNetworkClock(ctx context.Context, params *GetNetworkClockParams, reqEditors ...RequestEditorFn) (*http.Response, error)
// [x] GetNetworkInformation(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
// [x] GetNetworkParameters(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
// [x] GetSettings(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] GetSharedWallet(ctx context.Context, walletId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [x] GetCurrentSmashHealth(ctx context.Context, params *GetCurrentSmashHealthParams, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] ListStakePools(ctx context.Context, params *ListStakePoolsParams, reqEditors ...RequestEditorFn) (*http.Response, error)
// [x] GetMaintenanceActions(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
// [x] ListWallets(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
// [x] GetWallet(ctx context.Context, walletId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [x] ListAddresses(ctx context.Context, walletId string, params *ListAddressesParams, reqEditors ...RequestEditorFn) (*http.Response, error)
// [x] ListAssets(ctx context.Context, walletId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] GetAssetDefault(ctx context.Context, walletId string, policyId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [x] GetAsset(ctx context.Context, walletId string, policyId string, assetName string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [x] GetDelegationFee(ctx context.Context, walletId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] GetWalletKey(ctx context.Context, walletId string, role string, index string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] GetShelleyWalletMigrationInfo(ctx context.Context, walletId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [x] GetUTxOsStatistics(ctx context.Context, walletId string, reqEditors ...RequestEditorFn) (*http.Response, error)
// [x] ListTransactions(ctx context.Context, walletId string, params *ListTransactionsParams, reqEditors ...RequestEditorFn) (*http.Response, error)
// [x] GetTransaction(ctx context.Context, walletId string, transactionId string, reqEditors ...RequestEditorFn) (*http.Response, error)

// === Select, Migrate, Quit, Join, Sign requests
// [ ] ByronSelectCoinsWithBody(ctx context.Context, walletId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] ByronSelectCoins(ctx context.Context, walletId string, body ByronSelectCoinsJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] SelectCoinsWithBody(ctx context.Context, walletId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] SelectCoins(ctx context.Context, walletId string, body SelectCoinsJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] MigrateShelleyWalletWithBody(ctx context.Context, walletId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] MigrateShelleyWallet(ctx context.Context, walletId string, body MigrateShelleyWalletJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] SignMetadataWithBody(ctx context.Context, walletId string, role string, index string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
// [ ] SignMetadata(ctx context.Context, walletId string, role string, index string, body SignMetadataJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

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

	addresses := *resp.JSON200
	if len(addresses) == 0 {
		log.Warnf("Cannot test Byron Address requests: no addresses available")
		return
	}
	addressId := addresses[0].Id
	s.Run(fmt.Sprintf("InspectAddress(%v)", addressId), func() {
		resp, err := s.client.InspectAddressWithResponse(s.ctx, addressId)
		s.NoError(err)
		s.resp(err, http.StatusOK, resp.HTTPResponse, resp.Body)
		s.NotNil(resp.JSON200)
		s.Nil(resp.JSON400)
		s.logObject("InspectAddress", resp.JSON200)
	})
}

func (s *ByronWalletTestSuite) TestByronListAssets() {
	resp, err := s.client.ListByronAssetsWithResponse(s.ctx, s.walletId)
	s.NoError(err)
	s.resp(err, http.StatusOK, resp.HTTPResponse, resp.Body)
	s.NotNil(resp.JSON200)
	s.Nil(resp.JSON406)
	s.logObject("ListByronAssets", resp.JSON200)

	assets := *resp.JSON200
	if len(assets) == 0 {
		log.Warnf("Cannot test ByronAsset requests: no assets available")
		return
	}
	policyId := assets[0].PolicyId
	assetName := assets[0].AssetName
	s.Run(fmt.Sprintf("GetByronAsset(%v)", policyId), func() {
		resp, err := s.client.GetByronAssetWithResponse(s.ctx, s.walletId, policyId, assetName)
		s.NoError(err)
		s.resp(err, http.StatusOK, resp.HTTPResponse, resp.Body)
		s.NotNil(resp.JSON200)
		s.Nil(resp.JSON404)
		s.Nil(resp.JSON406)
		s.logObject("GetByronAsset", resp.JSON200)
	})
	s.Run(fmt.Sprintf("GetByronAssetDefault(%v)", policyId), func() {
		resp, err := s.client.GetByronAssetDefaultWithResponse(s.ctx, s.walletId, policyId)
		s.NoError(err)
		s.resp(err, http.StatusOK, resp.HTTPResponse, resp.Body)
		s.NotNil(resp.JSON200)
		s.Nil(resp.JSON404)
		s.Nil(resp.JSON406)
		s.logObject("GetByronAssetDefault", resp.JSON200)
	})
}

func (s *ByronWalletTestSuite) TestGetByronWalletMigrationInfo() {
	resp, err := s.client.GetByronWalletMigrationInfoWithResponse(s.ctx, s.walletId)
	s.NoError(err)
	s.resp(err, http.StatusOK, resp.HTTPResponse, resp.Body)
	s.NotNil(resp.JSON200)
	s.Nil(resp.JSON403)
	s.Nil(resp.JSON404)
	s.Nil(resp.JSON406)
	s.logObject("GetByronWalletMigrationInfo", resp.JSON200)
}

func (s *ByronWalletTestSuite) TestGetByronUTxOsStatistics() {
	resp, err := s.client.GetByronUTxOsStatisticsWithResponse(s.ctx, s.walletId)
	s.NoError(err)
	s.resp(err, http.StatusOK, resp.HTTPResponse, resp.Body)
	s.NotNil(resp.JSON200)
	s.Nil(resp.JSON404)
	s.Nil(resp.JSON406)
	s.logObject("GetByronUTxOsStatistics", resp.JSON200)
}

func (s *ByronWalletTestSuite) TestByronListTransactions() {
	resp, err := s.client.ListByronTransactionsWithResponse(s.ctx, s.walletId, new(ListByronTransactionsParams))
	s.NoError(err)
	s.resp(err, http.StatusOK, resp.HTTPResponse, resp.Body)
	s.NotNil(resp.JSON200)
	s.Nil(resp.JSON400)
	s.Nil(resp.JSON404)
	s.Nil(resp.JSON406)
	s.logObject("ListByronTransactions", resp.JSON200)

	transactions := *resp.JSON200
	if len(transactions) == 0 {
		log.Warnf("Cannot test ByronTransaction requests: no transactions available")
		return
	}
	transactionId := transactions[0].Id
	s.Run(fmt.Sprintf("GetByronTransaction(%v)", transactionId), func() {
		resp, err := s.client.GetByronTransactionWithResponse(s.ctx, s.walletId, transactionId)
		s.NoError(err)
		s.resp(err, http.StatusOK, resp.HTTPResponse, resp.Body)
		s.NotNil(resp.JSON200)
		s.Nil(resp.JSON404)
		s.Nil(resp.JSON406)
		s.logObject("GetByronTransaction", resp.JSON200)
	})
}

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

	addresses := *resp.JSON200
	if len(addresses) == 0 {
		log.Warnf("Cannot test Address requests: no addresses available")
		return
	}
	addressId := addresses[0].Id
	s.Run(fmt.Sprintf("InspectAddress(%v)", addressId), func() {
		resp, err := s.client.InspectAddressWithResponse(s.ctx, addressId)
		s.NoError(err)
		s.resp(err, http.StatusOK, resp.HTTPResponse, resp.Body)
		s.NotNil(resp.JSON200)
		s.Nil(resp.JSON400)
		s.logObject("InspectAddress", resp.JSON200)
	})
}

func (s *SingleWalletTestSuite) TestGetDelegationFee() {
	resp, err := s.client.GetDelegationFeeWithResponse(s.ctx, s.walletId)
	s.NoError(err)
	s.resp(err, http.StatusOK, resp.HTTPResponse, resp.Body)
	s.NotNil(resp.JSON200)
	s.Nil(resp.JSON403)
	s.Nil(resp.JSON404)
	s.Nil(resp.JSON406)
	s.logObject("GetDelegationFee", resp.JSON200)
}

func (s *SingleWalletTestSuite) TestListAssets() {
	resp, err := s.client.ListAssetsWithResponse(s.ctx, s.walletId)
	s.NoError(err)
	s.resp(err, http.StatusOK, resp.HTTPResponse, resp.Body)
	s.NotNil(resp.JSON200)
	s.Nil(resp.JSON406)
	s.logObject("ListAssets", resp.JSON200)

	assets := *resp.JSON200
	if len(assets) == 0 {
		log.Warnf("Cannot test Asset requests: no assets available")
		return
	}
	policyId := assets[0].PolicyId
	assetName := assets[0].AssetName
	s.Run(fmt.Sprintf("GetAsset(%v)", policyId), func() {
		resp, err := s.client.GetAssetWithResponse(s.ctx, s.walletId, policyId, assetName)
		s.NoError(err)
		s.resp(err, http.StatusOK, resp.HTTPResponse, resp.Body)
		s.NotNil(resp.JSON200)
		s.Nil(resp.JSON404)
		s.Nil(resp.JSON406)
		s.logObject("GetAsset", resp.JSON200)
	})
	// TODO does not work - not sure why
	// s.Run(fmt.Sprintf("GetAssetDefault(%v)", policyId), func() {
	// 	resp, err := s.client.GetAssetDefaultWithResponse(s.ctx, s.walletId, policyId)
	// 	s.NoError(err)
	// 	s.resp(err, http.StatusOK, resp.HTTPResponse, resp.Body)
	// 	s.NotNil(resp.JSON200)
	// 	s.Nil(resp.JSON404)
	// 	s.Nil(resp.JSON406)
	// 	s.logObject("GetAssetDefault", resp.JSON200)
	// })
}

// TODO this endpoint currently leads to 404
// func (s *SingleWalletTestSuite) TestGetShelleyWalletMigrationInfo() {
// 	resp, err := s.client.GetShelleyWalletMigrationInfoWithResponse(s.ctx, s.walletId)
// 	s.NoError(err)
// 	s.resp(err, http.StatusOK, resp.HTTPResponse, resp.Body)
// 	s.NotNil(resp.JSON200)
// 	s.Nil(resp.JSON403)
// 	s.Nil(resp.JSON404)
// 	s.Nil(resp.JSON406)
// 	s.logObject("GetShelleyWalletMigrationInfo", resp.JSON200)
// }

func (s *SingleWalletTestSuite) TestGetUTxOsStatistics() {
	resp, err := s.client.GetUTxOsStatisticsWithResponse(s.ctx, s.walletId)
	s.NoError(err)
	s.resp(err, http.StatusOK, resp.HTTPResponse, resp.Body)
	s.NotNil(resp.JSON200)
	s.Nil(resp.JSON404)
	s.Nil(resp.JSON406)
	s.logObject("GetUTxOsStatistics", resp.JSON200)
}

func (s *SingleWalletTestSuite) TestListTransactions() {
	resp, err := s.client.ListTransactionsWithResponse(s.ctx, s.walletId, new(ListTransactionsParams))
	s.NoError(err)
	s.resp(err, http.StatusOK, resp.HTTPResponse, resp.Body)
	s.NotNil(resp.JSON200)
	s.Nil(resp.JSON400)
	s.Nil(resp.JSON404)
	s.Nil(resp.JSON406)
	s.logObject("ListTransactions", resp.JSON200)

	transactions := *resp.JSON200
	if len(transactions) == 0 {
		log.Warnf("Cannot test Transaction requests: no transactions available")
		return
	}
	transactionId := transactions[0].Id
	s.Run(fmt.Sprintf("GetTransaction(%v)", transactionId), func() {
		resp, err := s.client.GetTransactionWithResponse(s.ctx, s.walletId, transactionId)
		s.NoError(err)
		s.resp(err, http.StatusOK, resp.HTTPResponse, resp.Body)
		s.NotNil(resp.JSON200)
		s.Nil(resp.JSON404)
		s.Nil(resp.JSON406)
		s.logObject("GetTransaction", resp.JSON200)
	})
}

// TODO This leads to 404 in the current Daedalus cardano-wallet process
// func (s *SingleWalletTestSuite) TestGetSharedWallet() {
// 	resp, err := s.client.GetSharedWalletWithResponse(s.ctx, s.walletId)
// 	s.NoError(err)
// 	s.resp(err, http.StatusOK, resp.HTTPResponse, resp.Body)
// 	s.NotNil(resp.JSON200)
// 	s.Nil(resp.JSON404)
// 	s.Nil(resp.JSON406)
// 	s.logObject("GetSharedWallet", resp.JSON200)
// }
