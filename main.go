/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"os"

	"github.com/op/go-logging"
	"github.com/spf13/cobra"

	_ "net/http/pprof"
)

var logger = logging.MustGetLogger("main")

// Constants go here.
const cmdRoot = "core"

// The main command describes the service and
// defaults to printing the help message.
var mainCmd = &cobra.Command{
	Use: "",
	Run: func(cmd *cobra.Command, args []string) {
		run(args)
	},
}

var orderingEndpoint string

//command line params
var (
	//-c <config json file> defaults to "ccchecker.json"
	configFile string

	//-e <config env file>, defaults to "env.json"
	envFile string
)

func main() {
	mainFlags := mainCmd.PersistentFlags()
	mainFlags.StringVarP(&configFile, "config", "c", "ccchecker.json", "CC Checker config file ")
	mainFlags.StringVarP(&envFile, "env", "e", "env.json", "CC Checker's env file ")

	// On failure Cobra prints the usage message and error string, so we only
	// need to exit with a non-0 status
	if mainCmd.Execute() != nil {
		os.Exit(1)
	}
}

func run(args []string) {
	//initialize the env
	env := InitCCCheckerEnv(configFile, envFile)

	CCCheckerInit()

	//TODO make parameters out of report and verbose
	CCCheckerRun(env, true, true)
	fmt.Printf("Test complete\n")
	return
}
