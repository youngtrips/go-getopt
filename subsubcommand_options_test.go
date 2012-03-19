// Copyright (c) 2011, SoundCloud Ltd., Daniel Bornkessel
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Source code and contact info at http://github.com/kesselborn/go-getopt

package getopt

import (
	//"os"
	//"strings"
	"testing"
)

func equalSubCommandOptions(sco1 SubCommandOptions, sco2 SubCommandOptions) (equal bool) {
	if len(sco1) == len(sco2) {
		for key := range sco1 {
			if !equalOptionsArray(sco1[key], sco2[key]) {
				goto loopend
			}
		}
		equal = true
	}
loopend:

	return
}

func TestSubSubcommandOptionsConverter(t *testing.T) {
  ssco := SubSubCommandOptions{
    Global{
      "global description",
      Definitions{
        {"config|c", "config file", IsConfigFile | ExampleIsDefault, "/etc/visor.conf"},
        {"server|s", "doozer server", Optional, ""},
        {"scope", "scope", IsSubcommand, ""},
      },
    },
    Scopes{
      "app": {
        Global{
          "app description",
          {"command", "command to execute", IsSubcommand, ""},
          {"foo|f", "a param", Optional, ""},
        },
        SubCommands{
          "getenv": {
            "app getenv description",
            Definitions{
              {"name", "app's name", IsArg | Required, ""},
              {"key", "environment variable's name", IsArg | Required, ""},
            },
          },
          "revision": {
            "app revision description",
            Definitions{
              {"rev|r", "revision", IsArg | Required, ""},
            },
          },
        },
      },
    },
  }

  expectedAppSubOptions := SubCommandOptions{
    Global{
      "app description",
      Definitions{
        {"config|c", "config file", IsConfigFile | ExampleIsDefault, "/etc/visor.conf"},
        {"server|s", "doozer server", Optional, ""},
        {"scope", "scope", IsSubcommand, ""},
        {"command", "command to execute", IsSubcommand, ""},
      },
    },
    SubCommands{
      "getenv": {
        "app getenv description",
        Definitions{
          {"name", "app's name", IsArg | Required, ""},
          {"key", "environment variable's name", IsArg | Required, ""},
        },
      },
      "revision": {
        "app revision description",
        Definitions{
          {"rev|r", "revision", IsArg | Required, ""},
        },
      },
    },
  }

	expectedRevisionOptions := Options{
		"*": {
			{"config|c", "config file", IsConfigFile | ExampleIsDefault, "/etc/visor.conf"},
			{"server|s", "doozer server", Optional, ""},
			{"scope", "scope", IsSubcommand, ""},
			{"command", "command to execute", IsSubcommand, ""}},
		"list": {
			{"app|n", "app's name", IsArg | Required, ""}},
	}

	if _, err := sco.flattenToSubcommandOptions("app"); err != nil {
		t.Errorf("conversion SuSubCommandOptions -> SubCommandOptions failed (app); \nGot the following error: %s", err.Message)
	}

	if options, _ := sco.flattenToSubcommandOptions("app"); equalSubCommandOptions(options, expectedAppSubOptions) == false {
		t.Errorf("conversion SubSubCommandOptions -> SubCommandOptions failed (app); \nGot\n\t#%#v#\nExpected:\n\t#%#v#\n", options, expectedAppSubOptions)
	}

	if _, err := sco.flattenToSubcommandOptions("revision"); err != nil {
		t.Errorf("conversion SuSubCommandOptions -> SubCommandOptions failed (revision); \nGot the following error: %s", err.Message)
	}

	if options, _ := sco.flattenToSubcommandOptions("revision"); equalSubCommandOptions(options, expectedRevisionOptions) == false {
		t.Errorf("conversion SubSubCommandOptions -> SubCommandOptions failed (revision); \nGot\n\t#%#v#\nExpected:\n\t#%#v#\n", options, expectedRevisionOptions)
	}

	if _, err := sco.flattenToOptions("nonexistantsubcommand"); err.ErrorCode != UnknownSubcommand {
		t.Errorf("non existant sub command didn't throw error")
	}

}

//func TestSubcommandOptionsSubCommandFinder(t *testing.T) {
//	sco := SubCommandOptions{
//		"*": {
//			{"command", "command to execute", IsSubcommand, ""},
//			{"foo|f", "some arg", Optional, ""}},
//		"getenv": {
//			{"name", "app's name", IsArg | Required, ""},
//			{"key", "environment variable's name", IsArg | Required, ""}},
//	}
//
//	os.Args = []string{"prog", "getenv"}
//	if command, _ := sco.findSubcommand(); command != "getenv" {
//		t.Errorf("did not correctly find subcommand getenv")
//	}
//
//	os.Args = []string{"prog", "-f", "bar", "getenv", "name", "key"}
//	if command, _ := sco.findSubcommand(); command != "getenv" {
//		t.Errorf("did not correctly find subcommand getenv")
//	}
//
//	os.Args = []string{"prog"}
//	if _, err := sco.findSubcommand(); err == nil || err.ErrorCode != NoSubcommand {
//		t.Errorf("did not throw error on unknown subcommand")
//	}
//}
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
//			{"command", "command to execute", IsSubcommand, ""}},
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
//			{"command", "command to execute", IsSubcommand, ""}},
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