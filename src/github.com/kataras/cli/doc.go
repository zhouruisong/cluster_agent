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

/*

package main

import (
	"strconv"
	"github.com/kataras/cli"
)

func main() {
	app := cli.NewApp("httpserver", "converts current directory into http server", "0.0.1")

	app.Flag("directory", "C:/users/myfiles", "specify a current working directory") // $executable -d, $executable --directory $the_dir

	// create a command that can be used by more than one app if you want it
	listenCommand := cli.Command("listen", "starts the server") //  $executable listen

	listenCommand.Flag("host", "127.0.0.1", "specify an address listener")      // $executable listen -h, $executable listen --host $the_host
	listenCommand.Flag("port", 8080, "specify a port to listen")                // $executable listen -p, $executable listen --port $the_port
	listenCommand.Flag("dir", "", "current working directory")                  // $executable listen -d, $executable listen --dir $the_dir
	listenCommand.Flag("req", nil, "a required flag because nil default given") // $executable listen -r , $executable listen --req $the_req

	listenCommand.Action(listen)

	app.Command(listenCommand) //register this command to the app.

	app.Run(run)
}

func run(args cli.Flags) error {
	println("Running ONLY APP with -d = " + args.String("directory"))
	return nil
}

func listen(args cli.Flags) error {
	println("EXECUTE ONLY 'listen' with args\n1 host: ", args.String("host"))
	println("2 port: ", strconv.Itoa(args.Int("port")))
	return nil
}

*/
/* OR */
/*

package main

import (
	"strconv"
	"github.com/kataras/cli"
)

func main() {

	var globalFlags = cli.Flags{
		cli.Flag("directory", "C:/users/myfiles", "specify a current working directory"),
	}

	cli.App{"httpserver", "converts current directory into http server", "0.0.1", cli.Commands{
		cli.Command("listen", "starts the server").
			Flag("host", "127.0.0.1", "specify an address listener").
			Flag("port", 8080, "specify a port to listen").
			Flag("dir", "", "current working directory"). // It's defaults to empty string it is not required flag
			Flag("req", nil, "a required flag because nil default given"). // required flag because nil default given
			Action(listen),
	}, globalFlags}.Run(run)

}

func run(args cli.Flags) error {
	println("Running ONLY APP with -d = " + args.String("directory"))
	return nil
}

func listen(args cli.Flags) error {
	println("EXECUTE ONLY 'listen' with args\n1 host: ", args.String("host"))
	println("2 port: ", strconv.Itoa(args.Int("port")))
	return nil
}

*/

//Note that: --help (or -help, help, -h) global flag is automatically used and displays help message.
