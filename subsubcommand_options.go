// Copyright (c) 2011, SoundCloud Ltd., Daniel Bornkessel
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Source code and contact info at http://github.com/kesselborn/go-getopt

package getopt

import (
//"os"
//"path/filepath"
//"sort"
//"strings"
//"fmt"
)

type Scopes map[string]SubCommandOptions

type SubSubCommandOptions struct {
	Global Options
	Scopes Scopes
}

func (ssco SubSubCommandOptions) flattenToSubCommandOptions(scope string) (sco SubCommandOptions, err *GetOptError) {
	globalCommand := ssco.Global
	var present bool

	if sco, present = ssco.Scopes[scope]; present == true {
		//print("\n======= " + scope + "========\n"); print(strings.Replace(fmt.Sprintf("%#v", sco.SubCommands),  "getopt", "\n", -1)); print("\n================")
		sco.Global.Definitions = append(globalCommand.Definitions, sco.Global.Definitions...)
	} else {
		err = &GetOptError{UnknownSubCommand, "Unknown scope: " + scope}
	}

	return
}

func (ssco SubSubCommandOptions) flattenToOptions(scope string, subCommand string) (options Options, err *GetOptError) {
	if sco, err := ssco.flattenToSubCommandOptions(scope); err == nil {
		options, err = sco.flattenToOptions(subCommand)
	}

	return
}

func (ssco SubSubCommandOptions) findScope() (scope string, err *GetOptError) {
	options := ssco.Global
	scope = "*"

	_, arguments, _, _ := options.ParseCommandLine()

	if len(arguments) < 1 {
		err = &GetOptError{NoScope, "Couldn't find scope"}
	} else {
		scope = arguments[0]
		if _, present := ssco.Scopes[scope]; present != true {
			err = &GetOptError{UnknownScope, "Given scope '" + scope + "' not defined"}
		}
	}

	return
}

func (ssco SubSubCommandOptions) findScopeAndSubCommand() (scope string, subCommand string, err *GetOptError) {
	if scope, err = ssco.findScope(); err == nil {
		var sco SubCommandOptions

		if sco, err = ssco.flattenToSubCommandOptions(scope); err == nil {
			var arguments []string
			if _, arguments, _, err = sco.Global.ParseCommandLine(); len(arguments) > 1 {
				subCommand = arguments[1]
				if _, present := sco.SubCommands[subCommand]; present != true {
					err = &GetOptError{UnknownSubCommand, "Given sub command '" + subCommand + "' not defined"}
				}
			} else {
				err = &GetOptError{NoSubCommand, "Couldn't find sub command"}
			}
		}
	}

	return
}

//func (sco SubCommandOptions) ParseCommandLine() (subCommand string, options map[string]OptionValue, arguments []string, passThrough []string, err *GetOptError) {
//func (sco SubCommandOptions) ParseCommandLine() (subCommand string, options map[string]OptionValue, arguments []string, passThrough []string, err *GetOptError) {
//
//	if subCommand, err = sco.findSubCommand(); err == nil {
//		var flattenedOptions Options
//		if flattenedOptions, err = sco.flattenToOptions(subCommand); err == nil {
//			options, arguments, passThrough, err = flattenedOptions.ParseCommandLine()
//			arguments = arguments[1:]
//		}
//	}
//
//	return
//}
//
//func (sco SubCommandOptions) Usage(scope string) (output string) {
//	return sco.UsageCustomArg0(scope, filepath.Base(os.Args[0]))
//}
//
//func (sco SubCommandOptions) UsageCustomArg0(scope string, arg0 string) (output string) {
//	subCommand, err := sco.findSubCommand()
//
//	if err != nil {
//		subCommand = "*"
//	} else {
//		arg0 = arg0 + " " + scope
//	}
//
//	if flattenedOptions, present := sco[subCommand]; present {
//		output = flattenedOptions.UsageCustomArg0(arg0)
//	}
//
//	return
//}
//
//func (sco SubCommandOptions) Help(description string, scope string) (output string) {
//	return sco.HelpCustomArg0(description, scope, filepath.Base(os.Args[0]))
//}
//
//func (sco SubCommandOptions) HelpCustomArg0(description string, scope string, arg0 string) (output string) {
//	subCommand, err := sco.findSubCommand()
//
//	if err != nil {
//		subCommand = "*"
//	} else {
//		arg0 = arg0 + " " + scope
//	}
//
//	if flattenedOptions, present := sco[subCommand]; present {
//		output = flattenedOptions.HelpCustomArg0(description, arg0)
//	}
//
//	if subCommand == "*" {
//		output = output + "Available commands:\n"
//
//		keys := make([]string, len(sco))
//		i := 0
//
//		for k := range sco {
//			keys[i] = k
//			i = i + 1
//		}
//		sort.Strings(keys)
//
//		for _, key := range keys {
//			if key != "*" {
//				output = output + "    " + key + "\n"
//			}
//		}
//		output = output + "\n"
//	}
//
//	return
//}
