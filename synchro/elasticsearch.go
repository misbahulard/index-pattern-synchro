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
package synchro

import (
	"crypto/tls"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/misbahulard/index-pattern-synchro/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type EsIndex []struct {
	Health       string `json:"health"`
	Status       string `json:"status"`
	Index        string `json:"index"`
	UUID         string `json:"uuid"`
	Pri          string `json:"pri"`
	Rep          string `json:"rep"`
	DocsCount    string `json:"docs.count"`
	DocsDeleted  string `json:"docs.deleted"`
	StoreSize    string `json:"store.size"`
	PriStoreSize string `json:"pri.store.size"`
}

func getEsIndex() []string {
	url := viper.GetString("elasticsearch.host") + "/_cat/indices?format=json&pri=true"

	client := resty.New().
		SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: !viper.GetBool("elasticsearch.ssl_certificate_verification"),
		})

	if viper.GetBool("elasticsearch.auth.enable") {
		client.SetBasicAuth(viper.GetString("elasticsearch.auth.username"), viper.GetString("elasticsearch.auth.password"))
	}

	resp, err := client.R().Get(url)

	if err != nil {
		log.Fatalln(err)
	}

	var esIndex EsIndex
	json.Unmarshal(resp.Body(), &esIndex)

	list := []string{}
	debugFilteredIndex := []string{}
	debugFilteredIndexWithRolloverPatten := []string{}

	var cfg config.Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Errorf("An error occured when unmarshal: %s", err)
	}

	// filter elasticsearch index
	for _, c := range cfg.Elasticsearch.Indices {
		re, err := regexp.Compile(checkAndFixRegex(c.Name))
		if err != nil {
			log.Errorf("An error occured when compile regex pattern for index name: %s", err)
		}

		// for each elastic index apply the filter
		for _, e := range esIndex {
			// skip system index
			if string(e.Index[0]) == "." {
				continue
			}

			matched := re.FindAll([]byte(e.Index), -1)

			// if index match with filter
			if len(matched) != 0 {
				indexName := e.Index
				debugFilteredIndex = append(debugFilteredIndex, indexName)

				// if rollover pattern is defined, filter with the elastic index with the pattern before append to the list
				if c.RolloverPattern != "" {
					re, err := regexp.Compile(checkAndFixRegex(c.RolloverPattern))
					if err != nil {
						log.Errorf("An error occured when compile regex pattern for index rollover pattern: %s", err)
					}

					patternBytes := re.FindAll([]byte(e.Index), -1)

					// if index match with rollover pattern
					if len(patternBytes) != 0 {
						pattern := string(patternBytes[0])
						indexWithoutRolloverSuffix := strings.Split(e.Index, pattern)[0]
						lastChar := indexWithoutRolloverSuffix[len(indexWithoutRolloverSuffix)-1:]

						// save index name
						indexName = indexWithoutRolloverSuffix

						// remove the special char in the last char when it exist
						if stringInSlice(lastChar, []string{"-", "_", "."}) {
							indexName = indexWithoutRolloverSuffix[0 : len(indexWithoutRolloverSuffix)-1]
						}

						debugFilteredIndexWithRolloverPatten = append(debugFilteredIndexWithRolloverPatten, indexName)
					}
				}

				list = append(list, indexName)
			}
		}
	}

	// print debug values
	log.Debugf("Filtered Raw Index: %+v", debugFilteredIndex)
	log.Debugf("Filtered Index with Rollover Pattern: %+v", unique(debugFilteredIndexWithRolloverPatten))
	log.Debugf("Final Index: %+v", unique(list))

	return unique(list)
}
