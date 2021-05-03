package main

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/godano/cardano-wallet-client/wallet"
	"github.com/sirupsen/logrus"
)

const fixedArguments = 2 // All relevant methods have: receiver, context

var (
	methodRegex  = regexp.MustCompile("(?P<verb>[A-Z][a-z]+)(?P<object>[[:alpha:]]+)")
	methodPrefix = regexp.MustCompile("^([A-Z][^A-Z]+)")
)

// Hacky way to make the CLI cleaner
var methodNameRemappings = map[string]string{
	"GetShelleyWalletMigrationInfo": "GetWalletMigrationInfo",
}

// No consistent way to remove plural S automatically (example problems: address, statistics)
var objectRemappings = map[string]string{
	"Transactions": "Transaction",
	"Assets":       "Asset",
	"Wallets":      "Wallet",
	"Addresses":    "Address",
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
		if len(match) == 0 {
			continue
		}
		cm := &clientMethod{
			walletCLI:    c,
			method:       method,
			methodName:   methodName,
			methodVerb:   strings.ToLower(match[1]),
			methodObject: match[2],
		}
		if remapping, ok := objectRemappings[cm.methodObject]; ok {
			cm.methodObject = remapping
		}

		if wallet.MethodHasBody[method.Name] || strings.HasSuffix(method.Name, "WithBody") {
			c.log.Debugf("Skipping method with HTTP body: %v", methodName)
			continue
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
	*walletCLI
	method reflect.Method

	methodName   string
	methodObject string
	methodVerb   string
	args         []*clientMethodArg

	// Same method, but with Byron* prefix
	byronVariant *reflect.Method
}

func (c *clientMethod) String() string {
	return fmt.Sprintf("%v %v (method name: %v, args: %v)",
		c.methodObject, c.methodVerb, c.methodName, c.args)
}

func (c *clientMethod) init(method reflect.Method) error {
	methodType := method.Func.Type()
	numArgs := methodType.NumIn()
	if !methodType.IsVariadic() {
		return fmt.Errorf("Method %v is not variadic", method.Name)
	}
	if numArgs < (fixedArguments + 1) {
		// Expect one additional parameter: variadic slice
		return fmt.Errorf("Method %v has unexpected number of arguments: %v",
			method.Name, numArgs)
	}
	numArgs -= fixedArguments + 1 // Ignore receiver, context, and variadic argument

	c.args = make([]*clientMethodArg, numArgs)
	argNames, hasArgNames := wallet.ArgumentNames[method.Name]
	if !hasArgNames || len(argNames) != len(c.args) {
		return fmt.Errorf("Missing argument names (or unexpected number) for method %v", method.Name)
	}
	for i := range c.args {
		argType := methodType.In(i + fixedArguments)
		arg := &clientMethodArg{
			method: c,
			name:   argNames[i],
			typ:    argType,
		}
		if err := arg.checkSupportedArgType(argType); err != nil {
			return err
		}
		arg.initStructValue()
		c.args[i] = arg
	}
	return nil
}

func (c *clientMethod) mergeByronVariant(methods map[string]*clientMethod) {
	// Insert the "Byron" part at the right position
	prefix := methodPrefix.FindString(c.methodName)
	byronName := prefix + "Byron" + c.methodName[len(prefix):]

	if byronMethod, ok := methods[byronName]; ok {
		c.byronVariant = &byronMethod.method
		delete(methods, byronName)
	}
}

func (c *clientMethod) call(receiver interface{}, ctx context.Context, commandArgs []string, useByronVariant bool) (*http.Response, error) {
	method := c.method
	if useByronVariant {
		method = *c.byronVariant
	}

	numArgs := method.Func.Type().NumIn() - 1 // Skip variadic part
	methodArgs := make([]reflect.Value, numArgs)
	methodArgs[0] = reflect.ValueOf(receiver) // Receiver
	methodArgs[1] = reflect.ValueOf(ctx)      // Context
	for i, arg := range c.args {
		// HACK This only works because an optional *struct parameter always comes last and other parameters are always strings
		var argVal interface{}
		if i < len(commandArgs) {
			argVal = commandArgs[i]
		} else {
			argVal = arg.value
		}
		methodArgs[i+fixedArguments] = reflect.ValueOf(argVal)
	}

	if c.log.Level >= logrus.DebugLevel {
		formattedArgs := make([]string, len(methodArgs))
		for i, arg := range methodArgs {
			formattedArgs[i] = fmt.Sprintf("%+v", arg)
		}
		c.log.Debugf("Calling %v with arguments: %v", method.Name, formattedArgs)
	}
	return c.unpackResult(method.Func.Call(methodArgs))
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

type clientMethodArg struct {
	method *clientMethod
	name   string
	typ    reflect.Type

	value interface{} // Only for struct args, fields filled through flags
}

func (a *clientMethodArg) String() string {
	return a.name
}

func (a *clientMethodArg) checkSupportedArgType(argType reflect.Type) error {
	if argType.Kind() == reflect.String {
		return nil
	}
	if a.isStructParameter() {
		return nil
	}

	// TODO support HTTP body parameters

	return fmt.Errorf("Currently only string and pointer-to-struct args supported, cannot handle '%v' of kind %v",
		argType.Name(), argType.Kind())
}

func (a *clientMethodArg) isStructParameter() bool {
	if a.typ.Kind() != reflect.Ptr || a.typ.Elem().Kind() != reflect.Struct {
		return false
	}
	for i := 0; i < a.numFields(); i++ {
		structField := a.typ.Elem().Field(i)
		if structField.Anonymous {
			return false
		}
		fieldType := structField.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		if !a.isSupportedFieldType(fieldType) {
			return false
		}
	}
	return true
}

func (a *clientMethodArg) isSupportedFieldType(typ reflect.Type) bool {
	return typ.Kind() == reflect.Bool || typ.Kind() == reflect.String || typ.Kind() == reflect.Int
}

func (a *clientMethodArg) initStructValue() {
	if a.isStructParameter() {
		a.value = wallet.MakeArgument(a.method.method.Name)
	}
}

func (a *clientMethodArg) numFields() int {
	return a.typ.Elem().NumField()
}
