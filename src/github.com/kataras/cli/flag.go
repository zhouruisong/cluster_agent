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
	"reflect"
	"strings"
)

type Flags []*flag

type flag struct { //lowercase in order to have the ability to do cli.Flag(...) instead of cli.NewFlag(...)
	Name    string
	Default interface{}
	Usage   string
	Value   interface{}
	Raw     *goflags.FlagSet
}

func Flag(name string, defaultValue interface{}, usage string, raw *goflags.FlagSet) *flag {
	return &flag{name, defaultValue, usage, nil, raw}
}

func (f flag) Alias() string {
	if len(f.Name) > 1 {
		return f.Name[0:1]
	}
	return f.Name
}

// Get returns a flag by it's name, if flag not found returns nil
func (c Flags) Get(name string) *flag {
	for idx, v := range c {
		if v.Name == name {
			return c[idx]
		}
	}

	return nil
}

// String returns the flag's value as string by it's name, if not found returns empty string ""
func (c Flags) String(name string) string {
	f := c.Get(name)
	if f == nil {
		return ""
	}
	return *f.Value.(*string) //*f.Value if string
}

// Bool returns the flag's value as bool by it's name, if not found returns false
func (c Flags) Bool(name string) bool {
	f := c.Get(name)
	if f != nil {
		return *f.Value.(*bool)
	}
	return false
}

// Int returns the flag's value as int by it's name, if can't parse int then returns -1
func (c Flags) Int(name string) int {
	f := c.Get(name)
	if f == nil {
		return -1
	}
	return *f.Value.(*int)
}

// IsValid returns true if flags are valid, otherwise false
func (c Flags) IsValid() bool {
	if c.Validate() != nil {
		return false
	}
	return true
}

// Validate returns nil if this flags are valid, otherwise returns an error message
func (c Flags) Validate() error {
	var notFilled []string
	for _, v := range c {
		// if no value given (nil) for required flag then it is not valid
		isRequired := v.Default == nil
		val := reflect.ValueOf(v.Value).Elem().String()
		if isRequired && val == "" {
			notFilled = append(notFilled, v.Name)
		}
	}

	if len(notFilled) > 0 {
		if len(notFilled) == 1 {
			return fmt.Errorf("Required flag [-%s] is missing.\n", notFilled[0])
		} else {
			return fmt.Errorf("Required flags [%s] are missing.\n", strings.Join(notFilled, ","))
		}

	}
	return nil

}

func (c Flags) ToString() (summary string) {
	for idx, v := range c {
		summary += "-" + v.Alias()
		if idx < len(c)-1 {
			summary += ", "
		}
	}

	if len(summary) > 0 {
		summary = "[" + summary + "]"
	}

	return
}

func requestFlagValue(flagset *goflags.FlagSet, name string, defaultValue interface{}, usage string) interface{} {
	if defaultValue == nil { // if it's nil then set it to a string because we will get err: interface is nil, not string if we pass a required flag
		defaultValue = ""
	}
	switch defaultValue.(type) {
	case int:
		{
			valPointer := flagset.Int(name, defaultValue.(int), usage)

			// it's not h (-h) for example but it's host, then assign it's alias also
			if len(name) > 1 {
				alias := name[0:1]
				flagset.IntVar(valPointer, alias, defaultValue.(int), usage)
			}
			return valPointer
		}
	case bool:
		{
			valPointer := flagset.Bool(name, defaultValue.(bool), usage)

			// it's not h (-h) for example but it's host, then assign it's alias also
			if len(name) > 1 {
				alias := name[0:1]
				flagset.BoolVar(valPointer, alias, defaultValue.(bool), usage)
			}
			return valPointer
		}
	default:
		valPointer := flagset.String(name, defaultValue.(string), usage)

		// it's not h (-h) for example but it's host, then assign it's alias also
		if len(name) > 1 {
			alias := name[0:1]
			flagset.StringVar(valPointer, alias, defaultValue.(string), usage)
		}

		return valPointer

	}
}
