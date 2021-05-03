package main

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unsafe"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func (c *walletCLI) rootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "godano-wallet-cli object operation",
		Short: "CLI for the cardano-wallet REST API",
		Long: `godano-wallet-cli connects to the REST API of a cardano-wallet process and
translates the CLI commands and parameters to appropriate REST API calls`,
	}

	flags := rootCmd.PersistentFlags()
	flags.StringVarP(&c.serverAddress, "server", "s", c.serverAddress, "Endpoint of the cardano-wallet process to connect to")
	flags.BoolVarP(&c.logQuiet, "quiet", "q", c.logQuiet, "Set the log level to Warning")
	flags.BoolVarP(&c.logVeryQuiet, "quieter", "Q", c.logVeryQuiet, "Set the log level to Error")
	flags.BoolVarP(&c.logVerbose, "verbose", "v", c.logVerbose, "Set the log level to Debug")
	flags.BoolVarP(&c.outputYAML, "yaml", "y", c.outputYAML, "Output responses as YAML instead of JSON (more compact)")
	return rootCmd
}

func (c *clientMethod) getObjectCommand(rootCmd *cobra.Command, objectVerbs map[string][]string) *cobra.Command {
	objectCommand, ok := c.objectCommands[c.methodObject]
	if !ok {
		// These messages cover the case that there are multiple sub-commands for this object
		verbs := objectVerbs[c.methodObject]
		sort.Strings(verbs)
		var joinedVerbs string
		if len(verbs) == 1 {
			joinedVerbs = verbs[0]
		} else if len(verbs) == 2 {
			joinedVerbs = fmt.Sprintf("%v or %v", verbs[0], verbs[1])
		} else {
			joinedVerbs = strings.Join(verbs[:len(verbs)-1], ", ") + ", or " + verbs[len(verbs)-1]
		}

		shortMessage := fmt.Sprintf("%v %v objects", joinedVerbs, c.methodObject)
		longMessage := shortMessage

		objectCommand = &cobra.Command{
			Use:     c.methodObject,
			Short:   shortMessage,
			Long:    longMessage,
			Aliases: []string{strings.ToLower(c.methodObject)},
		}
		rootCmd.AddCommand(objectCommand)
		c.objectCommands[c.methodObject] = objectCommand
	}
	return objectCommand
}

func (c *clientMethod) registerCommand(rootCmd *cobra.Command, clientObj interface{}, objectVerbs map[string][]string) {
	objectCommand := c.getObjectCommand(rootCmd, objectVerbs)
	isOnlyCommand := len(objectVerbs[c.methodObject]) == 1

	var cmd *cobra.Command
	if isOnlyCommand {
		// For objects with only one verb, there is no sub-command
		c.configureCommand(objectCommand)
		cmd = objectCommand
	} else {
		cmd = new(cobra.Command)
		cmd.Use = c.methodVerb
		c.configureCommand(cmd)
		objectCommand.AddCommand(cmd)
	}

	cmd.Short = fmt.Sprintf("%v %v", c.methodVerb, c.methodObject)
	cmd.Long = fmt.Sprintf("%v operation for %v objects", c.methodVerb, c.methodObject)
}

func (c *clientMethod) configureCommand(cmd *cobra.Command) {
	numArgs := 0
	for _, arg := range c.args {
		if arg.isStructParameter() {
			arg.addCommandFlags(cmd.Flags())
		} else {
			numArgs++
			cmd.Use += fmt.Sprintf(" <%v>", arg.name)
		}
	}
	cmd.Args = cobra.ExactArgs(numArgs)

	var useByronVariant bool
	cmd.Run = func(cmd *cobra.Command, args []string) {
		client, err := c.connectClient()
		c.checkErr(err)
		res, err := c.call(client, c.ctx, args, useByronVariant)
		c.checkErr(err)
		if err == nil {
			c.outputResponse(res)
		}
	}

	if c.byronVariant != nil {
		cmd.Use += " [--byron]"
		cmd.Flags().BoolVar(&useByronVariant, "byron", false, "Use the Byron variant of this command")
	}
}

var whitespaceRegex = regexp.MustCompile("\\s*")

func (a *clientMethodArg) addCommandFlags(flags *pflag.FlagSet) {
	val := reflect.ValueOf(a.value)
	for i := 0; i < a.numFields(); i++ {
		structField := a.typ.Elem().Field(i)
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
// The double-pointer

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
