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
	"fmt"
	"io/ioutil"
	"os"

	"github.com/misbahulard/index-pattern-synchro/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var output string

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the config file for quick start.",
	Long: `Initialize the config file for quick start, by default it will create
in your home directory '$HOME/.synchro/config.yaml',
or you can define we should place the file.

By default, we read the config file automatically in /etc/synchro, /opt/synchro,
$HOME/.synchro, or the current working directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		generateConfig()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&output, "output", "o", "", "Define where should we put the generated config file.")
}

func generateConfig() {
	destination := "config.yaml"

	if output != "" {
		destination = output
	}

	if _, err := os.Stat(destination); !os.IsNotExist(err) {
		fmt.Println("File already exist, skip.")
		os.Exit(0)
	}

	// generate config from struct
	cfg := config.Config{
		Interval: "15m",
		Log: config.Log{
			Debug: false,
			Type:  "text",
			File: config.File{
				Enable: false,
				Path:   "default.log",
			},
		},
		Elasticsearch: config.Elasticsearch{
			Host:                       "http://localhost:9200",
			SslCertificateVerification: true,
			Auth: config.Auth{
				Enable:   false,
				Username: "elastic",
				Password: "secret",
			},
			Indices: []config.Index{
				config.Index{
					Name:            "*",
					RolloverPattern: "[0-9]{4}.[0-9]{2}.[0-9]{2}",
				},
			},
		},
		Kibana: config.Kibana{
			Host:                       "http://localhost:5601",
			SslCertificateVerification: true,
			Auth: config.Auth{
				Enable:   false,
				Username: "elastic",
				Password: "secret",
			},
		},
		Xpack: config.Xpack{
			Enable: false,
			Spaces: []config.Space{
				config.Space{Name: "global", Pattern: "*", Timestamp: "timestamp"},
			},
		},
		Opendistro: config.Opendistro{
			Enable: false,
			Tenants: []config.Tenant{
				config.Tenant{Name: "global", Pattern: "*", Timestamp: "timestamp"},
			},
		},
	}

	fmt.Printf("Create the config file: %s ", destination)

	out, err := yaml.Marshal(cfg)
	if err != nil {
		fmt.Printf("An error occured when marshal to yaml: %s\n", err)
		os.Exit(1)
	}

	err = ioutil.WriteFile(destination, out, 0644)
	if err != nil {
		fmt.Printf("An error occured when write file: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("[ok]")
}
