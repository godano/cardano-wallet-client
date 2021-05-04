package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unsafe"

	"github.com/godano/cardano-wallet-client/wallet"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var dryRunErr = errors.New("Request dry-run")

func (c *walletCLI) initRootCommand() {
	c.rootCmd = &cobra.Command{
		Use:   "godano-wallet-cli object operation",
		Short: "CLI for the cardano-wallet REST API",
		Long: `godano-wallet-cli connects to the REST API of a cardano-wallet process and
translates the CLI commands and parameters to appropriate REST API calls`,
	}

	flags := c.rootCmd.PersistentFlags()
	flags.StringVarP(&c.serverAddress, "server", "s", c.serverAddress, "Endpoint of the cardano-wallet process to connect to")
	flags.BoolVarP(&c.logQuiet, "quiet", "q", c.logQuiet, "Set the log level to Warning")
	flags.BoolVarP(&c.logVeryQuiet, "quieter", "Q", c.logVeryQuiet, "Set the log level to Error")
	flags.BoolVarP(&c.logVerbose, "verbose", "v", c.logVerbose, "Set the log level to Debug")
	flags.BoolVarP(&c.logTrace, "trace", "V", c.logTrace, "Set the log level to Trace")
	flags.BoolVarP(&c.dryRun, "dry-run", "n", c.dryRun, "Show the resulting request instead of executing it")
	flags.BoolVarP(&c.outputYAML, "yaml", "y", c.outputYAML, "Output responses as YAML instead of JSON (more compact)")
}

func (c *walletCLI) initByronCommand() {
	c.byronCmd = &cobra.Command{
		Use:   "Byron",
		Short: "Commands for Byron-era objects",
		Long: `This sub-command bundles functionality for objects from the Byron era.
All Byron-era commands have equivalents in the main command.`,
		Aliases: []string{"byron"},
	}
	c.rootCmd.AddCommand(c.byronCmd)
}

func (c *walletCLI) objectCommand(m *methodCommand, allObjectVerbs []string) *cobra.Command {
	objectCommands := c.objectCommands
	parentCmd := c.rootCmd
	if m.method.isByronMethod {
		objectCommands = c.byronObjectCommands
		parentCmd = c.byronCmd
	}

	objectCommand, ok := objectCommands[m.method.object]
	if !ok {
		// These messages cover the case that there are multiple sub-commands for this object
		var joinedVerbs string
		if len(allObjectVerbs) == 1 {
			joinedVerbs = allObjectVerbs[0]
		} else if len(allObjectVerbs) == 2 {
			joinedVerbs = fmt.Sprintf("%v or %v", allObjectVerbs[0], allObjectVerbs[1])
		} else {
			joinedVerbs = strings.Join(allObjectVerbs[:len(allObjectVerbs)-1], ", ") + ", or " + allObjectVerbs[len(allObjectVerbs)-1]
		}

		shortMessage := fmt.Sprintf("%v %v objects", joinedVerbs, m.objectStr())
		longMessage := shortMessage

		objectCommand = &cobra.Command{
			Use:     m.method.object,
			Short:   shortMessage,
			Long:    longMessage,
			Aliases: []string{strings.ToLower(m.method.object)},
		}
		parentCmd.AddCommand(objectCommand)
		objectCommands[m.method.object] = objectCommand
	}
	return objectCommand
}

type methodCommand struct {
	cli    *walletCLI
	method *clientMethod

	// Only if method.hasParams or method.hasBody
	extraArg interface{}

	// Only if method.hasBody
	bodyFile    string
	bodyContent string
}

func (c *methodCommand) verbCommand(allObjectVerbs []string) {
	objectCommand := c.cli.objectCommand(c, allObjectVerbs)
	isOnlyCommand := len(allObjectVerbs) == 1

	var cmd *cobra.Command
	if isOnlyCommand {
		// For objects with only one verb, there is no sub-command
		c.configureCommand(objectCommand)
		cmd = objectCommand
	} else {
		cmd = new(cobra.Command)
		cmd.Use = c.method.verb
		c.configureCommand(cmd)
		objectCommand.AddCommand(cmd)
	}

	cmd.Short = fmt.Sprintf("%v %v objects", c.method.verb, c.objectStr())
	cmd.Long = fmt.Sprintf("%v operation for %v objects", c.method.verb, c.objectStr())
}

func (c *methodCommand) objectStr() string {
	res := c.method.object
	if c.method.isByronMethod {
		res = "Byron-era " + res
	}
	return res
}

func (c *methodCommand) configureCommand(cmd *cobra.Command) {
	cmd.Args = cobra.ExactArgs(len(c.method.stringArgs))
	for _, arg := range c.method.stringArgs {
		cmd.Use += fmt.Sprintf(" <%v>", arg)
	}

	if c.method.hasParams || c.method.hasBody {
		methodName := c.method.method.Name
		c.extraArg = wallet.MakeArgument(methodName)
		if c.extraArg == nil {
			// If this happens, the wallet.MakeArgument is outdated
			panic(fmt.Errorf("Failed to create argument for method %v", methodName))
		}
		if c.method.hasParams {
			c.addParamsFlags(cmd.Flags())
		} else if c.method.hasBody {
			c.addBodyFlags(cmd.Flags())
		}
	}

	cmd.Run = func(cmd *cobra.Command, args []string) {
		c.callMethod(args)
	}
}

func (c *methodCommand) callMethod(stringArgs []string) {
	args, err := c.buildMethodArguments(stringArgs)
	c.cli.checkErr(err)

	// Log all parameters for debugging
	if c.cli.log.Level >= logrus.DebugLevel {
		formattedArgs := make([]string, len(args))
		for i, arg := range args {
			formattedArgs[i] = fmt.Sprintf("%+v", arg)
		}
		c.cli.log.Debugf("Calling %v with arguments: %v", c.method.method.Name, formattedArgs)
	}

	// Log the request in case of dry-run
	if c.cli.dryRun {
		reqEditor := func(ctx context.Context, req *http.Request) error {
			c.cli.outputDryRunRequest(req)
			return dryRunErr
		}

		// Add a request-editor as variadic argument to prevent the request from actually executing
		args = append(args, reqEditor)
	}

	// Actually run the request method now
	res, err := c.method.call(args)
	if c.cli.dryRun && err == dryRunErr {
		return
	}
	c.cli.checkErr(err)
	c.cli.outputResponse(res)
}

func (c *methodCommand) buildMethodArguments(stringArgs []string) ([]interface{}, error) {
	// Load the optional body first, before connecting to the server
	if c.method.hasBody {
		if err := c.loadBody(); err != nil {
			return nil, err
		}
	}

	// Connect
	client, err := c.cli.connectClient()
	if err != nil {
		return nil, err
	}

	// First to arguments: receiver and context
	args := []interface{}{client, c.cli.ctx}

	// Afterwards come the string arguments
	for _, stringArg := range stringArgs {
		args = append(args, stringArg)
	}

	// Finally, an optional extra non-string argument (params or body)
	if c.method.hasParams || c.method.hasBody {
		extraArg := c.extraArg
		if c.method.extraArgumentType.Kind() != reflect.Ptr {
			// Need to get an interface{} for the struct, instead of the pointer
			// The struct type, that c.extraArg points to, is correct, therefore resolve the pointer
			extraArg = reflect.ValueOf(extraArg).Elem().Interface()
		}
		args = append(args, extraArg)
	}

	return args, nil
}

func (c *methodCommand) addBodyFlags(flags *pflag.FlagSet) {
	flags.StringVarP(&c.bodyContent, "body", "b", "", "Specify JSON-encoded content to send as request body")
	flags.StringVarP(&c.bodyFile, "body-file", "B", "", "Specify a JSON file to send as request body")
}

func (c *methodCommand) loadBody() error {
	if c.bodyFile != "" && c.bodyContent != "" {
		return fmt.Errorf("Cannot specify both --body/-b and --body-file/-B")
	}

	// Read the provided data into a buffer
	var bodyBytes []byte
	if c.bodyContent != "" {
		bodyBytes = []byte(c.bodyContent)
	} else if c.bodyFile != "" {
		var err error
		bodyBytes, err = ioutil.ReadFile(c.bodyFile)
		if err != nil {
			return err
		}
	}

	// Parse the buffer into the body object (already initialized).
	// This requires that wallet.MakeArgument() always returns pointers.
	// If necessary, the c.extraArg pointer is resolved in buildMethodArguments().
	if len(bodyBytes) > 0 {
		err := json.Unmarshal(bodyBytes, c.extraArg)
		if err != nil {
			return fmt.Errorf("Failed to parse body into type %T: %v", c.extraArg, err)
		}
	}
	return nil
}

var whitespaceRegex = regexp.MustCompile("\\s*")

func (c *methodCommand) addParamsFlags(flags *pflag.FlagSet) {
	val := reflect.ValueOf(c.extraArg)
	argType := c.method.extraArgumentType.Elem()
	for i := 0; i < argType.NumField(); i++ {
		structField := argType.Field(i)
		fieldType := structField.Type
		fieldName := structField.Name
		flagName := strings.ToLower(whitespaceRegex.ReplaceAllString(fieldName, ""))
		flagUsage := fieldName + " parameter"
		field := val.Elem().FieldByName(fieldName)
		flagPointer := unsafe.Pointer(field.Addr().Pointer())

		if fieldType.Kind() == reflect.Ptr {
			// For pointer values, get a pointer from the pflags library and set it directly into the struct
			var flagValue pflag.Value
			switch fieldType.Elem().Kind() {
			case reflect.Bool:
				indirectPtr := (**bool)(flagPointer)
				flagValue = &boolValue{target: indirectPtr}
			case reflect.String:
				indirectPtr := (**string)(flagPointer)
				flagValue = &stringValue{target: indirectPtr}
			case reflect.Int:
				indirectPtr := (**int)(flagPointer)
				flagValue = &intValue{target: indirectPtr}
			}
			flags.Var(flagValue, flagName, flagUsage)
		} else {
			// For non-pointer values, give the pflags library the address of the value inside the struct
			switch fieldType.Kind() {
			case reflect.Bool:
				flags.BoolVar((*bool)(flagPointer), flagName, false, flagUsage)
			case reflect.String:
				flags.StringVar((*string)(flagPointer), flagName, "", flagUsage)
			case reflect.Int:
				flags.IntVar((*int)(flagPointer), flagName, 0, flagUsage)
			}
		}
	}
}

// The types below are copied and modified from github.com/spf13/pflag and are necessary to avoid setting
// default values (such as empty strings) when the user does not specify the respective flag.
// This makes the double-pointer necessary.

type stringValue struct {
	target **string
}

func (s *stringValue) Set(val string) error {
	*s.target = &val
	return nil
}

func (s *stringValue) Type() string {
	return "string"
}

func (s *stringValue) String() string {
	if s.target == nil {
		return "<double-none>"
	}
	if *s.target == nil {
		return "<none>"
	}
	return string(**s.target)
}

type boolValue struct {
	target **bool
}

func (b *boolValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	*b.target = &v
	return err
}

func (b *boolValue) Type() string {
	return "bool"
}

func (b *boolValue) String() string {
	if *b.target == nil {
		return "<none>"
	}
	return strconv.FormatBool(bool(**b.target))
}

type intValue struct {
	target **int
}

func (i *intValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	vInt := int(v)
	*i.target = &vInt
	return err
}

func (i *intValue) Type() string {
	return "int"
}

func (i *intValue) String() string {
	if *i.target == nil {
		return "<none>"
	}
	return strconv.Itoa(int(**i.target))
}
