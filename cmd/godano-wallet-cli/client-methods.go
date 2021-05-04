package main

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/godano/cardano-wallet-client/wallet"
)

const fixedArguments = 2 // All relevant methods have: receiver, context

var (
	methodRegex         = regexp.MustCompile("(?P<verb>[A-Z][a-z]+)(?P<object>[[:alpha:]]+)")
	ignoredMethodsRegex = regexp.MustCompile("WithBody$")
	methodPrefix        = regexp.MustCompile("^([A-Z][^A-Z]+)")

	byronMethodRegex = regexp.MustCompile("Byron") // Simple regex: contains "Byron" anywhere
)

// Hacky way to make the CLI cleaner
var methodNameRemappings = map[string]string{
	"GetShelleyWalletMigrationInfo": "GetWalletMigrationInfo",
	"MigrateShelleyWallet":          "MigrateWallet",
	"ByronSelectCoins":              "SelectByronCoins",
	"ImportAddresses":               "ImportAddressBatch",
}

// No consistent way to remove plural S automatically (example problems: address, statistics)
var objectRemappings = map[string]string{
	"Transactions":       "Transaction",
	"Assets":             "Asset",
	"Wallets":            "Wallet",
	"Addresses":          "Address",
	"AnyAddress":         "Address",
	"MaintenanceActions": "MaintenanceAction",
	"StakePools":         "StakePool",
}

func (c *walletCLI) inspectClientInterface(clientObj interface{}) map[string]*clientMethod {
	methods := make(map[string]*clientMethod)
	objType := reflect.TypeOf(clientObj)
	for i := 0; i < objType.NumMethod(); i++ {
		method := objType.Method(i)
		methodName := method.Name
		if remappedName, ok := methodNameRemappings[methodName]; ok {
			methodName = remappedName
		}

		match := methodRegex.FindStringSubmatch(methodName)
		if len(match) == 0 || ignoredMethodsRegex.MatchString(methodName) {
			continue
		}

		cm := &clientMethod{
			method: method,
			name:   methodName,
			verb:   strings.ToLower(match[1]),
			object: match[2],
		}
		if byronMethodRegex.MatchString(cm.object) {
			cm.isByronMethod = true
			cm.object = byronMethodRegex.ReplaceAllString(cm.object, "")
		}
		if remapping, ok := objectRemappings[cm.object]; ok {
			cm.object = remapping
		}

		err := cm.init(method)
		if err != nil {
			c.log.Debugf("Skipping method %v: %v", methodName, err)
			continue
		}
		c.log.Debugf("Found method %v", cm)
		methods[methodName] = cm
	}
	return methods
}

type clientMethod struct {
	method reflect.Method

	name          string // Some method names are modified by methodNameRemappings
	object        string
	verb          string
	isByronMethod bool
	stringArgs    []string

	hasParams         bool
	hasBody           bool
	extraArgumentType reflect.Type
}

func (c *clientMethod) String() string {
	return fmt.Sprintf("%v %v (method name: %v, byron: %v, hasParams: %v, hasBody: %v, args: %v)",
		c.object, c.verb, c.method.Name, c.isByronMethod, c.hasParams, c.hasBody, c.stringArgs)
}

// The following method signatures are accepted:
// func (receiver) MethodName(ctx, [string, ]*, [params *struct or body struct], [...])
// If the method is variadic, the variadic part is ignored
func (c *clientMethod) init(method reflect.Method) error {
	methodType := method.Func.Type()

	// Calculate, how many arguments we have to provide
	numArgs := methodType.NumIn()
	specialArgs := fixedArguments // Special arguments are receiver, context, and optional variadic part
	if methodType.IsVariadic() {
		specialArgs++
	}
	if numArgs < specialArgs {
		return fmt.Errorf("Method %v has not enough arguments: %v", method.Name, numArgs)
	}
	numArgs -= specialArgs

	argNames, hasArgNames := wallet.ArgumentNames[method.Name]
	if !hasArgNames || len(argNames) != numArgs {
		return fmt.Errorf("Missing argument names (or unexpected number) for method %v", method.Name)
	}

	// Validate all arguments
	for i := 0; i < numArgs; i++ {
		argName := argNames[i]
		argIndex := i + fixedArguments
		argType := methodType.In(argIndex)
		stringArg, err := c.initArgument(argName, argIndex, argType)
		if err != nil {
			return err
		}
		if !stringArg && i != numArgs-1 {
			// Method is not supported: non-string argument must be last argument
			return fmt.Errorf("Method does not have expected signature")
		}
	}
	return nil
}

func (c *clientMethod) initArgument(name string, index int, typ reflect.Type) (bool, error) {
	if isStringArg(typ) {
		// If we are still scanning string args, record the argument name
		c.stringArgs = append(c.stringArgs, name)
		return true, nil
	} else {
		// Otherwise, check the optional params/body argument
		switch {
		case wallet.MethodHasParamsStruct[c.method.Name] &&
			name == wallet.ParamsArgName &&
			isParamsArg(typ):
			c.hasParams = true
			c.extraArgumentType = typ
			return false, nil
		case wallet.MethodHasBody[c.method.Name] &&
			name == wallet.BodyArgName &&
			isBodyArg(typ):
			c.hasBody = true
			c.extraArgumentType = typ
			return false, nil
		default:
			return false, fmt.Errorf("Unexpected non-string parameter %v (%v), type: %v", index, name, typ)
		}
	}
}

func isStringArg(argType reflect.Type) bool {
	return argType.Kind() == reflect.String
}

func isParamsArg(argType reflect.Type) bool {
	if argType.Kind() != reflect.Ptr || argType.Elem().Kind() != reflect.Struct {
		return false
	}
	for i := 0; i < argType.Elem().NumField(); i++ {
		structField := argType.Elem().Field(i)
		if structField.Anonymous {
			return false
		}
		fieldType := structField.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		if !isSupportedParamsFieldType(fieldType) {
			return false
		}
	}
	return true
}

func isSupportedParamsFieldType(typ reflect.Type) bool {
	return typ.Kind() == reflect.Bool || typ.Kind() == reflect.String || typ.Kind() == reflect.Int
}

func isBodyArg(argType reflect.Type) bool {
	// Accept almost everything here, since we parse user-provided JSON data into it
	return argType.Kind() == reflect.Struct || argType.Kind() == reflect.Interface
}

func (c *clientMethod) call(arguments []interface{}) (*http.Response, error) {
	// Wrap all arguments as reflect.Value
	reflectArgs := make([]reflect.Value, len(arguments))
	for i, arg := range arguments {
		reflectArgs[i] = reflect.ValueOf(arg)
	}
	return c.unpackResult(c.method.Func.Call(reflectArgs))
}

func (c *clientMethod) unpackResult(result []reflect.Value) (*http.Response, error) {
	if len(result) != 2 {
		return nil, fmt.Errorf("Unexpected number of method outputs (%v): %v", len(result), result)
	} else {
		respValue, errValue := result[0], result[1]
		var resp *http.Response
		var ok bool
		if !respValue.IsNil() {
			resp, ok = respValue.Interface().(*http.Response)
			if !ok {
				return nil, fmt.Errorf("Unexpected first response value (expected *http.Response): %v", respValue.Interface())
			}
		}
		var err error
		if !errValue.IsNil() {
			err, ok = errValue.Interface().(error)
			if !ok {
				return nil, fmt.Errorf("Unexpected second response value (expected error): %v", errValue.Interface())
			}
		}
		return resp, err
	}
}
