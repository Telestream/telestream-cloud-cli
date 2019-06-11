package cli

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type CommandActionMock struct {
	mock.Mock
}

func (m *CommandActionMock) action() {

	m.Called()
}

func TestCommand(t *testing.T) {

	var testVector = []struct {
		name       string
		input      []string
		cmdName    string
		result     bool
		actionDone bool
	}{
		{"no arguments", []string{"program_name"}, "command", false, false},
		{"wrong arguments", []string{"program_name", "some_invalid_cmd"}, "command", false, false},
		{"arguments to call command", []string{"program_name", "command"}, "command", true, true},
		{"more arguments to call command", []string{"program_name", "command", "command2"}, "command", true, true},
	}

	for _, testEl := range testVector {

		t.Run(testEl.name, func(t *testing.T) {

			mock := new(CommandActionMock)

			mock.On("action").Return()
			cmd := NewCommand(testEl.cmdName, mock.action, "")

			res, _ := cmd.checkAndParse(testEl.input, 1)
			assert.Equal(t, res, testEl.result)

			numCalls := 0

			if testEl.actionDone {
				numCalls = 1
			}

			mock.AssertNumberOfCalls(t, "action", numCalls)
			assert.Equal(t, cmd.getName(), testEl.cmdName)
			cmd.printDescription()
		})
	}
}

type FlaggedActionMock struct {
	mock.Mock
}

func (m *FlaggedActionMock) action(flagMap FlagMap) {

	m.Called(flagMap)
}

func TestFlaggedCommand(t *testing.T) {

	var testVector = []struct {
		name        string
		input       []string
		cmdName     string
		flags       map[string]bool
		result      bool
		actionDone  bool
		valWoutFlag string
	}{
		{"no arguments", []string{"program_name"}, "fcommand", map[string]bool{"fflag": true, "sflag": false}, false, false, ""},
		{"no flag passed", []string{"program_name", "fcommand"}, "fcommand", map[string]bool{"fflag": true, "sflag": false}, true, false, ""},
		{"sflag", []string{"program_name", "fcommand", "-sflag", "sflag_value"}, "fcommand", map[string]bool{"fflag": true, "sflag": false}, true, false, ""},
		{"fflag", []string{"program_name", "fcommand", "-fflag", "fflag_value"}, "fcommand", map[string]bool{"fflag": true, "sflag": false}, true, true, ""},
		{"fflag", []string{"program_name", "fcommand", "-sflag", "sflag_value", "-fflag", "fflag_value"}, "fcommand",
			map[string]bool{"fflag": true, "sflag": false}, true, true, ""},
		{"fflag without flag", []string{"program_name", "fcommand", "fflag_value", "-sflag", "sflag_value"}, "fcommand",
			map[string]bool{"fflag": true, "sflag": false}, true, true, "fflag"},
		{"print help", []string{"program_name", "fcommand", "help"}, "fcommand",
			map[string]bool{"fflag": true, "sflag": false}, true, false, "fflag"},
		{"no flag required", []string{"program_name", "fcommand"}, "fcommand",
			map[string]bool{"fflag": false, "sflag": false}, true, true, ""},
	}

	for _, testEl := range testVector {
		t.Run(testEl.name, func(t *testing.T) {

			mockFact := new(FlaggedActionMock)
			mockFact.On("action", mock.Anything).Return()
			cmd := NewFlaggedCommand(testEl.cmdName, testEl.valWoutFlag, mockFact.action, testEl.flags, "")

			res, _ := cmd.checkAndParse(testEl.input, 1)
			assert.Equal(t, res, testEl.result)

			numCalls := 0

			if testEl.actionDone {
				numCalls = 1
			}

			mockFact.AssertNumberOfCalls(t, "action", numCalls)

			cmd.printFlags(true)
			cmd.printFlags(false)
			assert.Equal(t, cmd.getName(), testEl.cmdName)
			cmd.printDescription()
		})
	}
}

func TestSubCommand(t *testing.T) {

	var testVector = []struct {
		name       string
		input      []string
		cmdName    string
		nexCommads []CommandBaseInterface
		result     bool
	}{
		{"no arguments", []string{"program_name"}, "scommand", []CommandBaseInterface{}, false},
		{"call subcommand", []string{"program_name", "scommand"}, "scommand", []CommandBaseInterface{}, true},
		{"call subcommand", []string{"program_name", "scommand"}, "scommand",
			[]CommandBaseInterface{NewCommand("command", nil, "")}, true},
		{"call command", []string{"program_name", "scommand", "command"}, "scommand",
			[]CommandBaseInterface{NewCommand("command", nil, "")}, true},
		{"call command2", []string{"program_name", "scommand", "command2"}, "scommand",
			[]CommandBaseInterface{NewCommand("command", nil, ""), NewCommand("command2", nil, "")}, true},
	}

	for _, testEl := range testVector {
		t.Run(testEl.name, func(t *testing.T) {

			cmd := NewSubCommand(testEl.cmdName, testEl.nexCommads, "")
			res, _ := cmd.checkAndParse(testEl.input, 1)
			assert.Equal(t, res, testEl.result)

			cmd.printSubcommands()
			assert.Equal(t, cmd.getName(), testEl.cmdName)
			cmd.printDescription()
		})
	}
}

func TestAdditionalFlags(t *testing.T) {

	var testVector = []struct {
		name        string
		input       []string
		flags       map[string]string
		argvOutput  []string
		flagsOutput map[string]string
	}{
		// No flag passed test
		{"no flag", []string{"program_name", "first_command", "second_command"}, map[string]string{"some_flag": "some_flag_desc"},
			[]string{"program_name", "first_command", "second_command"}, map[string]string{"some_flag": ""}},
		// One flag passed on the beginning (--)
		{"one flag on begin", []string{"program_name", "--some_flag", "some_flag_value", "first_command", "second_command"}, map[string]string{"some_flag": "some_flag_desc"},
			[]string{"program_name", "first_command", "second_command"}, map[string]string{"some_flag": "some_flag_value"}},
		// One flag passed on the beginning (-some_flag=some_flag_value)
		{"one flag on begin -=", []string{"program_name", "-some_flag=some_flag_value", "first_command", "second_command"}, map[string]string{"some_flag": "some_flag_desc"},
			[]string{"program_name", "first_command", "second_command"}, map[string]string{"some_flag": "some_flag_value"}},
		// One flag one the middle
		{"one flag on middle", []string{"program_name", "first_command", "-some_flag=some_flag_value", "second_command"}, map[string]string{"some_flag": "some_flag_desc"},
			[]string{"program_name", "first_command", "second_command"}, map[string]string{"some_flag": "some_flag_value"}},
		// One flag at the end
		{"one flag at the end", []string{"program_name", "first_command", "second_command", "-some_flag=some_flag_value"}, map[string]string{"some_flag": "some_flag_desc"},
			[]string{"program_name", "first_command", "second_command"}, map[string]string{"some_flag": "some_flag_value"}},
		// check if program will fail when no value passed at the end
		{"check if fail", []string{"program_name", "first_command", "second_command", "-some_flag"}, map[string]string{"some_flag": "some_flag_desc"},
			[]string{"program_name", "first_command", "second_command"}, map[string]string{"some_flag": ""}},
		// Two flags
		{"two flags", []string{"program_name", "--flag1", "flag1_value", "--flag2", "flag2_value", "first_command", "second_command"},
			map[string]string{"flag1": "flag1 description", "flag2": "flag2 description"},
			[]string{"program_name", "first_command", "second_command"}, map[string]string{"flag1": "flag1_value", "flag2": "flag2_value"}},
	}

	for _, testEl := range testVector {

		t.Run(testEl.name, func(t *testing.T) {

			argvOutput, flags := GetAdditionalFlags(testEl.input, testEl.flags)

			if !reflect.DeepEqual(argvOutput, testEl.argvOutput) {

				t.Errorf("%v vs %v", argvOutput, testEl.argvOutput)
				t.Errorf("Wrong output argv")
			}

			for key, val := range flags {

				if *val != testEl.flagsOutput[key] {

					t.Errorf("Wrong output flags")
				}
			}

		})
	}

}

func TestCommandHandler(t *testing.T) {

	var testVector = []struct {
		name     string
		input    []string
		commands []CommandBaseInterface
		addFlags map[string]string
	}{
		{"no arguments", []string{"program_name"}, []CommandBaseInterface{}, map[string]string{}},
		{"help", []string{"program_name", "help"}, []CommandBaseInterface{NewCommand("command", nil, "")}, map[string]string{}},
		{"call command", []string{"program_name", "command"}, []CommandBaseInterface{NewCommand("command", nil, "")}, map[string]string{"flag": ""}},
		{"list command", []string{"program_name", "command1"}, []CommandBaseInterface{NewCommand("command", nil, "")}, map[string]string{"flag": ""}},
	}

	for _, testEl := range testVector {

		t.Run(testEl.name, func(t *testing.T) {

			argv, _ := GetAdditionalFlags(testEl.input, testEl.addFlags)
			cmd := NewCommandHandler("program_name", testEl.commands, testEl.addFlags)
			cmd.ParseArgs(argv)
			cmd.printAdditionalFlags()
			cmd.printCommands()
		})
	}
}
