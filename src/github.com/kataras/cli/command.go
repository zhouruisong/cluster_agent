// Copyright (c) 2016, Gerasimos Maropoulos
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
//    this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//	  this list of conditions and the following disclaimer
//    in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse
//    or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER AND CONTRIBUTOR, GERASIMOS MAROPOULOS
// BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package cli

import (
	goflags "flag"
	"strings"
)

type (
	Action func(Flags) error

	Commands []*command

	command struct { //lowercase in order to have the ability to do cli.Command(...) instead of cli.NewCommand
		Name        string
		Description string
		// Flags are not the arguments was given by the user, but the flags that developer sets to this command
		Flags       Flags
		action      Action
		Subcommands Commands
		flagset     *goflags.FlagSet
	}
)

func DefaultAction(cmdName string) Action {
	return func(a Flags) error { Printf(ErrNoAction, cmdName); return nil }
}

func Command(name string, description string) *command {
	name = strings.Replace(name, "-", "", -1) //removes all - if present, --help -> help
	fset := goflags.NewFlagSet(name, goflags.PanicOnError)
	return &command{Name: name, Description: description, Flags: Flags{}, action: DefaultAction(name), flagset: fset}
}

// Subcommand adds a child command (subcommand)
func (c *command) Subcommand(subCommand *command) *command {
	if c.Subcommands == nil {
		c.Subcommands = Commands{}
	}

	c.Subcommands = append(c.Subcommands, subCommand)
	return c
}

func (c *command) Flag(name string, defaultValue interface{}, usage string) *command {
	if c.Flags == nil {
		c.Flags = Flags{}
	}
	valPointer := requestFlagValue(c.flagset, name, defaultValue, usage)

	newFlag := &flag{name, defaultValue, usage, valPointer, c.flagset}
	c.Flags = append(c.Flags, newFlag)
	return c
}

func (c *command) Action(action Action) *command {
	c.action = action
	return c
}

// Execute returns true if this command has been executed
func (c *command) Execute(parentflagset *goflags.FlagSet) bool {
	var index = -1
	// check if this command has been called from app's arguments
	for idx, a := range parentflagset.Args() {
		if c.Name == a {
			index = idx + 1
		}
	}

	// this command hasn't been called from the user
	if index == -1 {
		return false
	}

	//check if it was help sub command
	wasHelp := parentflagset.Arg(1) == "-h"

	if wasHelp {
		// global -help, --help, help, -h now shows all the help for each subcommand and subflags
		Printf("Please use global '-help' or 'help' without quotes, instead.")
		return true
	}

	if !c.flagset.Parsed() {

		if err := c.flagset.Parse(parentflagset.Args()[index:]); err != nil {
			panic("Panic on command.Execute: " + err.Error())
		}
	}

	if err := c.Flags.Validate(); err == nil {
		c.action(c.Flags)

		for idx := range c.Subcommands {
			if c.Subcommands[idx].Execute(c.flagset) {
				break
			}

		}
	} else {
		Printf(err.Error())
	}
	return true

}
