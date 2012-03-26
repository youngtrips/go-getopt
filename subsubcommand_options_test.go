// Copyright (c) 2011, SoundCloud Ltd., Daniel Bornkessel
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Source code and contact info at http://github.com/kesselborn/go-getopt

package getopt

import (
	"os"
	//"strings"
	"testing"
)

func testingSubSubDefinitions() (ssco SubSubCommandOptions) {
	ssco = SubSubCommandOptions{
		Options{
			"global description",
			Definitions{
				{"config|c", "config file", IsConfigFile | ExampleIsDefault, "/etc/visor.conf"},
				{"server|s", "doozer server", Optional, ""},
				{"scope", "scope", IsSubCommand, ""},
			},
		},
		Scopes{
			"app": {
				Options{
					"app description",
					Definitions{
						{"foo|f", "a param", Optional, ""},
						{"command", "command to execute", IsSubCommand, ""},
					},
				},
				SubCommands{
					"getenv": {
						"app getenv description",
						Definitions{
							{"key", "environment variable's name", IsArg | Required, ""},
						},
					},
					"setenv": {
						"app setenv description",
						Definitions{
							{"name", "app's name", IsArg | Required, ""},
							{"key", "environment variable's name", IsArg | Required, ""},
						},
					},
				},
			},
			"revision": {
				Options{
					"app revision description",
					Definitions{
						{"rev|r", "revision", IsArg | Required, ""},
						{"command", "command to execute", IsSubCommand, ""},
					},
				},
				SubCommands{
					"list": {
						"list revisions",
						Definitions{
							{"all|a", "long list output", Flag, ""},
						},
					},
				},
			},
		},
	}

	return
}
func TestSubSubCommandOptionsConverter(t *testing.T) {
	ssco := testingSubSubDefinitions()

	expectedAppSubOptions := SubCommandOptions{
		Options{
			"app description",
			Definitions{
				{"config|c", "config file", IsConfigFile | ExampleIsDefault, "/etc/visor.conf"},
				{"server|s", "doozer server", Optional, ""},
				{"scope", "scope", IsSubCommand, ""},
				{"foo|f", "a param", Optional, ""},
				{"command", "command to execute", IsSubCommand, ""},
			},
		},
		SubCommands{
			"getenv": {
				"app getenv description",
				Definitions{
					{"key", "environment variable's name", IsArg | Required, ""},
				},
			},
			"setenv": {
				"app setenv description",
				Definitions{
					{"name", "app's name", IsArg | Required, ""},
					{"key", "environment variable's name", IsArg | Required, ""},
				},
			},
		},
	}

	expectedRevisionOptions := SubCommandOptions{
		Options{
			"app revision description",
			Definitions{
				{"config|c", "config file", IsConfigFile | ExampleIsDefault, "/etc/visor.conf"},
				{"server|s", "doozer server", Optional, ""},
				{"scope", "scope", IsSubCommand, ""},
				{"rev|r", "revision", IsArg | Required, ""},
				{"command", "command to execute", IsSubCommand, ""},
			},
		},
		SubCommands{
			"list": {
				"list revisions",
				Definitions{
					{"all|a", "long list output", Flag, ""},
				},
			},
		},
	}

	if _, err := ssco.flattenToSubCommandOptions("app"); err != nil {
		t.Errorf("conversion SuSubCommandOptions -> SubCommandOptions failed (app); \nGot the following error: %s", err.Message)
	}

	if sco, _ := ssco.flattenToSubCommandOptions("app"); equalSubCommandOptions(sco, expectedAppSubOptions) == false {
		t.Errorf("conversion SubSubCommandOptions -> SubCommandOptions failed (app); \nGot\n\t#%#v#\nExpected:\n\t#%#v#\n", sco, expectedAppSubOptions)
	}

	if _, err := ssco.flattenToSubCommandOptions("revision"); err != nil {
		t.Errorf("conversion SuSubCommandOptions -> SubCommandOptions failed (revision); \nGot the following error: %s", err.Message)
	}

	if sco, _ := ssco.flattenToSubCommandOptions("revision"); equalSubCommandOptions(sco, expectedRevisionOptions) == false {
		t.Errorf("conversion SubSubCommandOptions -> SubCommandOptions failed (revision); \nGot\n\t#%#v#\nExpected:\n\t#%#v#\n", sco, expectedRevisionOptions)
	}

	if _, err := ssco.flattenToSubCommandOptions("nonexistantsubcommand"); err.ErrorCode != UnknownSubCommand {
		t.Errorf("non existant sub command didn't throw error")
	}

	expectedAppGetEnvOptions := Options{
		"app getenv description",
		Definitions{
			{"config|c", "config file", IsConfigFile | ExampleIsDefault, "/etc/visor.conf"},
			{"server|s", "doozer server", Optional, ""},
			{"scope", "scope", IsSubCommand, ""},
			{"foo|f", "a param", Optional, ""},
			{"command", "command to execute", IsSubCommand, ""},
			{"key", "environment variable's name", IsArg | Required, ""},
		},
	}

	if _, err := ssco.flattenToOptions("app", "getenv"); err != nil {
		t.Errorf("conversion SubSubCommandOptions -> Options failed (app/getenv); \nGot the following error: %s", err.Message)
	}

	if options, _ := ssco.flattenToOptions("app", "getenv"); equalOptions(options, expectedAppGetEnvOptions) == false {
		t.Errorf("conversion SubCommandOptions -> Options failed (app/getenv); \nGot\n\t#%#v#\nExpected:\n\t#%#v#\n", options, expectedAppGetEnvOptions)
	}

}
func TestSubSubCommandScopeFinder(t *testing.T) {
	ssco := testingSubSubDefinitions()

	os.Args = []string{"prog", "app"}
	if command, _ := ssco.findScope(); command != "app" {
		t.Errorf("did not correctly find subcommand app")
	}

	os.Args = []string{"prog", "-s", "10.20.30.40", "app", "getenv", "key"}
	if command, _ := ssco.findScope(); command != "app" {
		t.Errorf("did not correctly find subcommand app (w/ other options)")
	}

	os.Args = []string{"prog"}
	if _, err := ssco.findScope(); err == nil || err.ErrorCode != NoScope {
		t.Errorf("did not throw error on unknown subcommand")
	}
}

func TestSubSubCommandSubCommand(t *testing.T) {
	ssco := testingSubSubDefinitions()

	os.Args = []string{"prog", "app", "getenv"}
	if scope, command, _ := ssco.findScopeAndSubCommand(); scope != "app" || command != "getenv" {
		t.Errorf("did not correctly find subcommand app / getenv (got: " + scope + " / " + command + ")")
	}

	os.Args = []string{"prog", "-s", "10.20.30.40", "app", "-ffoo", "getenv", "key"}
	if _, _, err := ssco.findScopeAndSubCommand(); err != nil {
		t.Errorf("did not correctly find subcommand app / getenv; Error message: " + err.Message)
	}

	if scope, command, _ := ssco.findScopeAndSubCommand(); scope != "app" || command != "getenv" {
		t.Errorf("did not correctly find subcommand app / getenv (got: " + scope + " / " + command + ")")
	}

	os.Args = []string{"prog"}
	if _, _, err := ssco.findScopeAndSubCommand(); err == nil || err.ErrorCode != NoScope {
		t.Errorf("did not throw error on missing scope")
	}

	os.Args = []string{"prog", "app"}
	if _, _, err := ssco.findScopeAndSubcommand(); err == nil || err.ErrorCode != NoSubcommand {
		t.Errorf("did not throw error on unknown subcommand")
	}
}

//func TestSubSubCommandOptionsParser(t *testing.T) {
//	ssco := testingSubSubDefinitions()
//
//func TestSubcommandOptionsParser(t *testing.T) {
//	sco := SubCommandOptions{
//		"*": {
//			{"command", "command to execute", IsSubcommand, ""},
//			{"foo|f", "some arg", Optional, ""}},
//		"getenv": {
//			{"bar|b", "some arg", Optional, ""},
//			{"name", "app's name", IsArg | Required, ""},
//			{"key", "environment variable's name", IsArg | Required, ""}},
//	}
//
//	os.Args = []string{"prog", "-fbar", "getenv", "--bar=foo", "foo", "bar"}
//	scope, options, arguments, _, _ := sco.ParseCommandLine()
//
//	if scope != "getenv" {
//		t.Errorf("SubCommandOptions parsing: failed to correctly parse scope: Expected: getenv, Got: " + scope)
//	}
//
//	if options["foo"].String != "bar" {
//		t.Errorf("SubCommandOptions parsing: failed to correctly parse option: Expected: bar, Got: " + options["foo"].String)
//	}
//
//	if options["bar"].String != "foo" {
//		t.Errorf("SubCommandOptions parsing: failed to correctly parse option: Expected:  foo, Got: " + options["foo"].String)
//	}
//
//	if scope != "getenv" {
//		t.Errorf("SubCommandOptions parsing: failed to correctly parse sub command: Expected: getenv, Got: " + scope)
//	}
//
//	if arguments[0] != "foo" {
//		t.Errorf("SubCommandOptions parsing: failed to correctly parse arg1: Expected: foo, Got: " + arguments[0])
//	}
//
//	if arguments[1] != "bar" {
//		t.Errorf("SubCommandOptions parsing: failed to correctly parse arg2: Expected: bar, Got: " + arguments[1])
//	}
//}
//
//func TestErrorMessageForMissingArgs(t *testing.T) {
//	sco := SubCommandOptions{
//		"*": {
//			{"foo|f", "some arg", Optional, ""},
//			{"command", "command to execute", IsSubCommand, ""}},
//		"getenv": {
//			{"bar|b", "some arg", Optional, ""},
//			{"name", "app's name", IsArg | Required, ""},
//			{"key", "environment variable's name", IsArg | Required, ""}},
//	}
//
//	os.Args = []string{"prog", "getenv"}
//	_, _, _, _, err := sco.ParseCommandLine()
//
//	if err == nil {
//		t.Errorf("missing arg did not raise error")
//	}
//
//	if expected := "Missing required argument <name>"; err.Message != expected {
//		t.Errorf("Error handling for missing arguments is messed up:\n\tGot     : " + err.Message + "\n\tExpected: " + expected)
//	}
//
//}
//
//func TestSubCommandHelp(t *testing.T) {
//	sco := SubCommandOptions{
//		"*": {
//			{"foo|f", "some arg", Optional, ""},
//			{"command", "command to execute", IsSubCommand, ""}},
//		"getenv": {
//			{"bar|b", "some arg", Optional, ""},
//			{"name", "app's name", IsArg | Required, ""},
//			{"key", "environment variable's name", IsArg | Required, ""}},
//		"register": {
//			{"deploytype|t", "deploy type (one of mount, bazapta, lxc)", NoLongOpt | Optional | ExampleIsDefault, "lxc"},
//			{"name|n", "app's name", IsArg | Required, ""}},
//	}
//
//	os.Args = []string{"prog"}
//	expectedHelp := `Usage: prog [-f <foo>] <command>
//
//this is not a program
//
//Options:
//    -f, --foo=<foo>           some arg
//    -h, --help                usage (-h) / detailed help text (--help)
//
//Available commands:
//    getenv
//    register
//
//`
//	expectedUsage := `Usage: prog [-f <foo>] <command>
//
//`
//
//	if got := sco.Help("this is not a program", "*"); got != expectedHelp {
//		t.Errorf("Usage output not as expected:\ngot:      |" + strings.Replace(got, " ", "_", -1) + "|\nexpected: |" + strings.Replace(expectedHelp, " ", "_", -1) + "|\n")
//	}
//
//	if got := sco.Usage("*"); got != expectedUsage {
//		t.Errorf("Usage output not as expected:\ngot:      |" + strings.Replace(got, " ", "_", -1) + "|\nexpected: |" + strings.Replace(expectedHelp, " ", "_", -1) + "|\n")
//	}
//
//	os.Args = []string{"prog", "register"}
//	expectedHelp = `Usage: prog register [-t <deploytype>] <name>
//
//this is not a program
//
//Options:
//    -t <deploytype>     deploy type (one of mount, bazapta, lxc) (default: lxc)
//    -h, --help          usage (-h) / detailed help text (--help)
//
//Arguments:
//    <name>              app's name
//
//`
//
//	expectedUsage = `Usage: prog register [-t <deploytype>] <name>
//
//`
//
//	if got := sco.Help("this is not a program", "register"); got != expectedHelp {
//		t.Errorf("Usage output not as expected:\ngot:      |" + strings.Replace(got, " ", "_", -1) + "|\nexpected: |" + strings.Replace(expectedHelp, " ", "_", -1) + "|\n")
//	}
//
//	if got := sco.Usage("register"); got != expectedUsage {
//		t.Errorf("Usage output not as expected:\ngot:      |" + strings.Replace(got, " ", "_", -1) + "|\nexpected: |" + strings.Replace(expectedHelp, " ", "_", -1) + "|\n")
//	}
//
//}
