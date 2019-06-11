package cli

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var commandNameC = color.New(color.FgHiWhite).Add(color.Bold)
var subCommandsC = color.New(color.FgRed).Add(color.Bold)
var requiredFlagC = color.New(color.FgHiBlue).Add(color.Bold)
var notRequiredFlagC = color.New(color.FgMagenta).Add(color.Bold)
var valueWFlagC = color.New(color.FgHiGreen).Add(color.Bold)
var commandDescriptionC = color.New(color.FgWhite).Add(color.Italic)

// CommandBase - holds basic command name and description
type CommandBase struct {
	name        string
	description string
}

// CommandBaseInterface - basic command interface
type CommandBaseInterface interface {
	checkAndParse(argv []string, argDepth int) (bool, error)
	getName() string
	printDescription()
}

// SubCommand - sub command type of basic command, in addition it can hold a slice of next commands
type SubCommand struct {
	CommandBase
	nextCommands []CommandBaseInterface
}

// CommandAction - defines function type that is called on command
type CommandAction func()

// Command - command type of basic command, in addition it holds a pointer on function (defined by CommandAction)
type Command struct {
	CommandBase
	cAction CommandAction
}

// FlagProperties - holds poiner that points to given flag value and information if this specisic flag is required
type FlagProperties struct {
	Value      *string
	IsRequired bool
}

// FlagMap key - is command name
type FlagMap map[string]FlagProperties

// ParsedAction - defines function type that is called on FlaggedCommand
type ParsedAction func(flagMap FlagMap)

// FlaggedCommand - command that specifies command additional flags
type FlaggedCommand struct {
	CommandBase
	valueWithoutFlag string
	flagMap          FlagMap
	pAction          ParsedAction
	pFlag            *flag.FlagSet
}

func isNextHelp(argv []string, argDepth int) bool {

	if argDepth+1 < len(argv) {
		if argv[argDepth+1] == "help" {

			return true
		}
	}

	return false
}

func isNextFlag(argv []string, argDepth int) bool {

	if argDepth+1 < len(argv) {
		if argv[argDepth+1][0] == '-' {

			return true
		}
	}

	return false
}

// Creates new sub command with specific name, next commands and description
func NewSubCommand(name string, commands []CommandBaseInterface, description string) *SubCommand {

	subCommand := new(SubCommand)
	subCommand.name = name
	subCommand.description = description
	subCommand.nextCommands = commands

	return subCommand
}

func (sCmd *SubCommand) checkAndParse(argv []string, argDepth int) (bool, error) {

	if argDepth < len(argv) && argv[argDepth] == sCmd.name {
		for _, cmd := range sCmd.nextCommands {
			if res, err := cmd.checkAndParse(argv, argDepth+1); res {

				return true, err
			}
		}

		sCmd.printSubcommands()
		return true, errors.New("Cannot match any sub command")
	}

	return false, nil
}

func (sCmd *SubCommand) printSubcommands() {

	commandNameC.Println()
	commandNameC.Println("Available commands for " + sCmd.name + ":")
	commandNameC.Println()

	for _, cmd := range sCmd.nextCommands {
		fmt.Print(strings.Repeat(" ", 3))
		cmd.printDescription()
	}

	commandNameC.Println()
}

func (sCmd *SubCommand) getName() string {

	return sCmd.name
}

func (sCmd *SubCommand) printDescription() {

	commandNameC.Print(sCmd.name + " ")
	subCommands := "("

	for idx, cmd := range sCmd.nextCommands {
		if 0 == idx {

			subCommands += cmd.getName()
		} else {

			subCommands += "|" + cmd.getName()
		}
	}
	subCommands += ") "

	subCommandsC.Print(subCommands)
	commandDescriptionC.Println("- " + sCmd.description)
}

// Creates Command with specific name, action (called when command passed) and desctiption
func NewCommand(name string, action CommandAction, description string) *Command {

	command := new(Command)
	command.name = name
	command.cAction = action
	command.description = description

	return command
}

func (cmd *Command) checkAndParse(argv []string, argDepth int) (bool, error) {

	if argDepth < len(argv) && argv[argDepth] == cmd.name {

		if argv[argDepth] == cmd.name {

			if cmd.cAction != nil {
				cmd.cAction()
			}
			return true, nil
		}
	}

	return false, nil
}

func (cmd *Command) printDescription() {

	commandNameC.Print(cmd.name + " ")
	commandDescriptionC.Println("- " + cmd.description)
}

func (cmd *Command) getName() string {

	return cmd.name
}

// Creates flagged command with specific name, action, attrubutes and description
func NewFlaggedCommand(name string, valueWithoutFlag string, action ParsedAction, flags map[string]bool, description string) *FlaggedCommand {

	flaggedCommand := new(FlaggedCommand)
	flaggedCommand.name = name
	flaggedCommand.valueWithoutFlag = valueWithoutFlag
	flaggedCommand.pAction = action
	flaggedCommand.description = description
	flaggedCommand.pFlag = flag.NewFlagSet(name, flag.ContinueOnError)
	flaggedCommand.pFlag.SetOutput(ioutil.Discard)
	flaggedCommand.flagMap = FlagMap{}

	if valueWithoutFlag != "" {

		flags[valueWithoutFlag] = true
	}

	for key, val := range flags {

		pVal := flaggedCommand.pFlag.String(key, "", "Required: "+strconv.FormatBool(val))
		flaggedCommand.flagMap[key] = FlagProperties{pVal, val}
	}

	return flaggedCommand
}

func (fCmd *FlaggedCommand) checkAndParse(argv []string, argDepth int) (bool, error) {

	actualArgIdx := 1

	if argDepth+actualArgIdx < len(argv) && argv[argDepth] == fCmd.name {

		// help passed -> print all possible flags
		if isNextHelp(argv, argDepth) {

			fCmd.printFlags(true)
			return true, errors.New("Print help")
		}

		if !isNextFlag(argv, argDepth) {

			nextArgIdx := argDepth + actualArgIdx
			fCmd.flagMap[fCmd.valueWithoutFlag] = FlagProperties{&argv[nextArgIdx], fCmd.flagMap[fCmd.valueWithoutFlag].IsRequired}
			actualArgIdx++
		}

		if argDepth+actualArgIdx < len(argv) {
			fCmd.pFlag.Parse(argv[argDepth+actualArgIdx:])
		}

		if fCmd.isAnyRequiredNotSet() {

			fCmd.printFlags(false)
			return true, errors.New("Some required flag not set")
		}

		fCmd.pAction(fCmd.flagMap)
		return true, nil

	} else if argDepth < len(argv) && argv[argDepth] == fCmd.name {

		if fCmd.isAnyFlagRequired() {

			fCmd.printFlags(false)
			return true, errors.New("No flag set")

		} else {

			fCmd.pAction(fCmd.flagMap)
			return true, nil
		}
	}

	return false, nil
}

func (fCmd *FlaggedCommand) printFlags(printAll bool) {

	commandNameC.Println()

	if printAll {

		commandNameC.Println("Flags for " + fCmd.name + ":")
	} else {

		commandNameC.Println("Required flags for " + fCmd.name + ":")
	}

	commandNameC.Println()

	for key, val := range fCmd.flagMap {
		if val.IsRequired || printAll {

			fmt.Print(strings.Repeat(" ", 3))

			if val.IsRequired {
				requiredFlagC.Println("-" + key + " " + "<" + strings.ToUpper(key) + "> (required)")
			} else {
				notRequiredFlagC.Println("-" + key + " " + "<" + strings.ToUpper(key) + "> ")
			}
		}
	}

	commandNameC.Println()
}

func (fCmd *FlaggedCommand) getName() string {

	return fCmd.name
}

func (fCmd *FlaggedCommand) printDescription() {

	commandNameC.Print(fCmd.name + " ")

	subCommandsC.Print("(help|) ")

	flags := ""

	if fCmd.valueWithoutFlag != "" {

		valueWFlagC.Print("<" + strings.ToUpper(fCmd.valueWithoutFlag) + ">/")
		requiredFlagC.Print("-" + fCmd.valueWithoutFlag + " " + "<" + strings.ToUpper(fCmd.valueWithoutFlag) + "> ")
	}

	for key, val := range fCmd.flagMap {

		if val.IsRequired && key != fCmd.valueWithoutFlag {
			flags += "-" + key + " " + "<" + strings.ToUpper(key) + "> "
		}
	}

	requiredFlagC.Print(flags)
	commandDescriptionC.Println("- " + fCmd.description)
}

func (fCmd *FlaggedCommand) isAnyFlagRequired() bool {

	for _, val := range fCmd.flagMap {
		if val.IsRequired {

			return true
		}
	}

	return false
}

func (fCmd *FlaggedCommand) isAnyRequiredNotSet() bool {

	for _, val := range fCmd.flagMap {
		if val.IsRequired && *val.Value == "" {

			return true
		}
	}

	return false
}

type CommandHandler struct {
	cmds            []CommandBaseInterface
	additionalFlags map[string]string
	programName     string
}

// Creates new command handler
func NewCommandHandler(programName string, cmds []CommandBaseInterface, additionalFlags map[string]string) *CommandHandler {

	commandHandler := new(CommandHandler)
	commandHandler.programName = programName
	commandHandler.cmds = cmds
	commandHandler.additionalFlags = additionalFlags

	return commandHandler
}

func GetAdditionalFlags(argv []string, nameDesc map[string]string) ([]string, map[string]*string) {

	parsedArgs := make(map[string]*string)
	pFlag := flag.NewFlagSet("", flag.ContinueOnError)
	pFlag.SetOutput(ioutil.Discard)
	retArgv := argv

	for key, val := range nameDesc {

		parsedArgs[key] = pFlag.String(key, "", val)
	}

	for idx, _ := range argv {

		pFlag.Parse(argv[idx:])
	}

	// remove additional flags from command line arguments
	for key, argVal := range parsedArgs {

		for idx, val := range retArgv {

			elsToDelete := 0

			if "-"+key+"="+*argVal == val || "--"+key+"="+*argVal == val {

				elsToDelete = 1
			} else if "-"+key == val || "--"+key == val {

				elsToDelete = 1

				if *argVal != "" {

					elsToDelete = 2
				}
			}

			if elsToDelete > 0 {

				retArgv = append(retArgv[:idx], retArgv[(idx+elsToDelete):]...)
				break
			}
		}
	}

	return retArgv, parsedArgs
}

// calls command maching to os.Args
func (cmdHndl *CommandHandler) ParseArgs(argv []string) {

	depth := 1

	if depth < len(argv) {
		if argv[depth] == "add_flags" {

			cmdHndl.printAdditionalFlags()
			return
		}
	}

	for _, cmd := range cmdHndl.cmds {

		if res, _ := cmd.checkAndParse(argv, depth); res {

			return
		}
	}

	if 0 == len(cmdHndl.cmds) {

		fmt.Println("Error - program has no commands")
		return
	}

	cmdHndl.printCommands()
}

func (cmdHndl *CommandHandler) printAdditionalFlags() {

	commandNameC.Println()
	commandNameC.Println("Additional flags all commands:")
	commandNameC.Println()

	for key, val := range cmdHndl.additionalFlags {

		fmt.Print(strings.Repeat(" ", 3))
		requiredFlagC.Print("-" + key + " " + "<" + strings.ToUpper(key) + "> ")
		commandDescriptionC.Println("- " + val)
	}
	commandNameC.Println()
}

func (cmdHndl *CommandHandler) printCommands() {

	commandNameC.Println()
	commandNameC.Println("Available commands for " + cmdHndl.programName + ":")
	commandNameC.Println()

	for _, cmd := range cmdHndl.cmds {

		fmt.Print(strings.Repeat(" ", 3))
		cmd.printDescription()
	}

	fmt.Print(strings.Repeat(" ", 3))
	commandNameC.Print("add_flags")
	commandDescriptionC.Println(" - prints additional flags for all commands")

	commandNameC.Println()
}
