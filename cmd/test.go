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
	"crypto/tls"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Verify all function and connection is ok.",
	Run: func(cmd *cobra.Command, args []string) {
		test()
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}

func test() {
	// test elasticsearch
	url := viper.GetString("elasticsearch.host")

	client := resty.New().
		SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: !viper.GetBool("elasticsearch.ssl_certificate_verification"),
		})

	if viper.GetBool("elasticsearch.auth.enable") {
		client.SetBasicAuth(viper.GetString("elasticsearch.auth.username"), viper.GetString("elasticsearch.auth.password"))
	}

	fmt.Printf("Elasticsearch.. ")

	resp, err := client.R().Get(url)
	if err != nil {
		fmt.Printf("An error occured when doing request to server: %s\n", err)
		os.Exit(1)
	}

	if resp.StatusCode() != 200 {
		fmt.Println("error")
		fmt.Printf("An error occured when connecting to Elasticsearch: %s\n", resp.Body())
		os.Exit(1)
	}

	fmt.Println("ok")

	// test kibana
	url = viper.GetString("kibana.host")

	client = resty.New().
		SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: !viper.GetBool("kibana.ssl_certificate_verification"),
		})

	if viper.GetBool("kibana.auth.enable") {
		client.SetBasicAuth(viper.GetString("kibana.auth.username"), viper.GetString("kibana.auth.password"))
	}

	fmt.Printf("Kibana.. ")

	resp, err = client.R().SetBody(`[{"type":"index-pattern","id":"test"}]`).Get(url)
	if err != nil {
		fmt.Printf("An error occured when doing request to server: %s\n", err)
		os.Exit(1)
	}

	if resp.StatusCode() != 200 {
		fmt.Println("error")
		fmt.Printf("An error occured when connecting to Kibana: %s\n", resp.Body())
		os.Exit(1)
	}

	fmt.Println("ok")
}
