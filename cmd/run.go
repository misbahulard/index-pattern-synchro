/*
Copyright © 2021 Misbahul Ardani <misbahulard@gmail.com>

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
package cmd

import (
	"github.com/misbahulard/index-pattern-synchro/logger"
	"github.com/misbahulard/index-pattern-synchro/synchro"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var debug bool

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run index pattern synchro.",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Configure()
		synchro.Run()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().BoolVar(&debug, "debug", false, "Enable debug log level.")
	viper.BindPFlag("log.debug", runCmd.Flags().Lookup("debug"))
}
