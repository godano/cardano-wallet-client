package wallet

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ParameterNamesTestSuite struct {
	suite.Suite
	*require.Assertions
}

func TestParameterNames(t *testing.T) {
	testSuite := new(ParameterNamesTestSuite)
	suite.Run(t, testSuite)
}

func (s *ParameterNamesTestSuite) SetupSuite() {
	s.Assertions = s.Require()
}

// TestParameterNames tests that `ArgumentNames` and `ArgumentNamesWithResponse` correspond to the methods in
// `Client` and `ClientWithResponse`. This test will report, if after updating the generated code, the manually
// defined parameter names must be updated.
func (s *ParameterNamesTestSuite) TestParameterNames() {
	client := new(ClientWithResponses)
	methods := getMethods(client)

	// Make a copy since we modify the map below
	paramNames := make(map[string][]string, len(ArgumentNames)+len(ArgumentNamesWithResponse))
	for method, names := range ArgumentNames {
		paramNames[method] = names
	}
	for method, names := range ArgumentNamesWithResponse {
		paramNames[method] = names
	}

	for _, method := range methods {
		s.Contains(paramNames, method.Name, "Every method in ClientWithResponses must be included in ArgumentNames or ArgumentNamesWithResponse")
		actualNumParams := len(paramNames[method.Name])
		actualNumParams++ // For the receiver
		actualNumParams++ // For the context
		if !strings.HasSuffix(method.Name, "WithResponse") {
			actualNumParams++ // For variadic reqEditors parameter
		}

		expectedNumParams := method.Type.NumIn()

		s.Equal(expectedNumParams, actualNumParams,
			"Method %v reported %v parameters, but actually has %v", method.Name, actualNumParams, expectedNumParams)
		delete(paramNames, method.Name)
	}

	s.Empty(paramNames, "ArgumentNames or ArgumentNamesWithResponse contain non-existing methods")
}

func getMethods(obj interface{}) []reflect.Method {
	clientType := reflect.TypeOf(obj)
	methods := make([]reflect.Method, clientType.NumMethod())
	for i := range methods {
		methods[i] = clientType.Method(i)
	}
	return methods
}
