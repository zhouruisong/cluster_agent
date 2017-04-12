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
	"fmt"
	"os"
	"text/template"
)

var Output = os.Stdout // the output is the same for all Apps, atm.

func Printf(format string, args ...interface{}) {
	fmt.Fprintf(Output, format, args...)
}

func HelpMe(app App) {
	tmplStr := appTmpl
	tmpl, err := template.New(app.Name).Parse(tmplStr)
	if err != nil {
		panic("Panic: " + err.Error())
	}

	tmpl.Execute(Output, app)

}

type App struct {
	Name        string
	Description string
	Version     string
	Commands    Commands
	Flags       Flags
}

func NewApp(name string, description string, version string) *App {
	return &App{name, description, version, nil, nil}
}

// Command adds a  command to the app
func (a *App) Command(cmd *command) *App {
	if a.Commands == nil {
		a.Commands = Commands{}
	}

	a.Commands = append(a.Commands, cmd)
	return a
}

func (a *App) Flag(name string, defaultValue interface{}, usage string) *App {
	if a.Flags == nil {
		a.Flags = Flags{}
	}

	a.Flags = append(a.Flags, &flag{name, defaultValue, usage, nil, nil})
	return a
}

func (a App) help() {
	HelpMe(a)
	os.Exit(1)
}

func (a App) HasCommands() bool {
	return a.Commands != nil && len(a.Commands) > 0
}

func (a App) HasFlags() bool {
	return a.Flags != nil && len(a.Flags) > 0
}

func (a App) Run(appAction Action) {

	flagset := goflags.NewFlagSet(a.Name, goflags.PanicOnError)
	flagset.SetOutput(Output)

	if a.Flags != nil {

		//now, get the args and set the flags
		for idx, arg := range a.Flags {
			valPointer := requestFlagValue(flagset, arg.Name, arg.Default, arg.Usage)
			a.Flags[idx].Value = valPointer
		}
	}

	if len(os.Args) <= 1 {
		a.help()
	}

	// if help argument/flag is passed
	if len(os.Args) > 1 && (os.Args[1] == "help" || os.Args[1] == "-help" || os.Args[1] == "--help") || os.Args[1] == "-h" {
		a.help()

	}
	// if flag parsing failed, yes we check it after --help.
	if err := flagset.Parse(os.Args[1:]); err != nil {
		a.help()
	}

	//first we check for commands, if any command executed then  app action should NOT be executed

	var ok = false
	for idx := range a.Commands {
		if ok = a.Commands[idx].Execute(flagset); ok {
			break
		}
	}

	if !ok {
		if err := a.Flags.Validate(); err == nil {
			if err = appAction(a.Flags); err == nil {
				ok = true
			} else {
				Printf(err.Error())
				return
			}
		} else {
			Printf(err.Error())
			return
		}
	}

	if !ok {
		a.help()
	}

}

func (a *App) Printf(format string, args ...interface{}) {
	Printf(format, args...)
}

var appTmpl = `NAME:
   {{.Name}} - {{.Description}}

USAGE:
{{- if .HasFlags}}
   {{.Name}} [global arguments...]
{{end -}}
{{ if .HasCommands }}
   {{.Name}} command [arguments...]
{{ end }}
VERSION:
   {{.Version}}

{{ if.HasFlags}}
GLOBAL ARGUMENTS:
{{ range $idx,$flag := .Flags }}
   -{{$flag.Alias }}        {{$flag.Usage}} (default '{{$flag.Default}}')
{{ end }}
{{end -}}
{{ if .HasCommands }}
COMMANDS:
{{ range $index, $cmd := .Commands }}
   {{$cmd.Name }} {{$cmd.Flags.ToString}}        {{$cmd.Description}}
     {{ range $index, $subcmd := .Subcommands }}
     {{$subcmd.Name}}        {{$subcmd.Description}}
	 {{ end }}
     {{ range $index, $subflag := .Flags }}
      -{{$subflag.Alias }}        {{$subflag.Usage}} (default '{{$subflag.Default}}')
	 {{ end }}
{{ end }}
{{ end }}
`
