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
	"os"
	"regexp"

	"github.com/go-resty/resty/v2"
	"github.com/misbahulard/index-pattern-synchro/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Space struct {
	Name      string
	Indices   []string
	Timestamp string
}

type Tenant struct {
	Name      string
	Indices   []string
	Timestamp string
}

type SavedObjectPayload struct {
	Type       string     `json:"type"`
	ID         string     `json:"id"`
	Attributes Attributes `json:"attributes"`
}

type Attributes struct {
	Title         string `json:"title"`
	TimeFieldName string `json:"timeFieldName"`
}

type SavedObjectResponse struct {
	SavedObjects []SavedObject `json:"saved_objects"`
}

type SavedObject struct {
	ID               string           `json:"id"`
	Type             string           `json:"type"`
	Error            Error            `json:"error"`
	Attributes       Attributes       `json:"attributes"`
	References       []interface{}    `json:"references"`
	MigrationVersion MigrationVersion `json:"migrationVersion"`
	UpdatedAt        string           `json:"updated_at"`
	Version          string           `json:"version"`
	Namespaces       []string         `json:"namespaces"`
}

type Error struct {
	StatusCode int64  `json:"statusCode"`
	Error      string `json:"error"`
	Message    string `json:"message"`
}

type MigrationVersion struct {
	IndexPattern string `json:"index-pattern"`
}

func getKibanaSpaces(esIndices []string) []Space {
	var cfg config.Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		log.Errorf("An error occured when unmarshal: %s", err)
		os.Exit(1)
	}

	spaces := []Space{}
	for _, s := range cfg.Xpack.Spaces {
		space := Space{}
		space.Name = s.Name
		if s.Timestamp == "" {
			space.Timestamp = "@timestamp"
		} else {
			space.Timestamp = s.Timestamp
		}

		list := []string{}
		for _, index := range esIndices {
			re, err := regexp.Compile(checkAndFixRegex(s.Pattern))
			if err != nil {
				log.Errorf("An error occured when compile regex pattern for index name: %s", err)
				os.Exit(1)
			}

			matched := re.FindAll([]byte(index), -1)

			// if index match with filter
			if len(matched) != 0 {
				list = append(list, index)
			}
		}

		space.Indices = list
		spaces = append(spaces, space)
	}

	log.Debugf("spaces: %+v", spaces)
	return spaces
}

func getKibanaTenants(esIndices []string) []Tenant {
	var cfg config.Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		log.Errorf("An error occured when unmarshal: %s", err)
		os.Exit(1)
	}

	tenants := []Tenant{}
	for _, s := range cfg.Opendistro.Tenants {
		tenant := Tenant{}
		tenant.Name = s.Name
		if s.Timestamp == "" {
			tenant.Timestamp = "@timestamp"
		} else {
			tenant.Timestamp = s.Timestamp
		}

		list := []string{}
		for _, index := range esIndices {
			re, err := regexp.Compile(checkAndFixRegex(s.Pattern))
			if err != nil {
				log.Errorf("An error occured when compile regex pattern for index name: %s", err)
				os.Exit(1)
			}

			matched := re.FindAll([]byte(index), -1)

			// if index match with filter
			if len(matched) != 0 {
				list = append(list, index)
			}
		}

		tenant.Indices = list
		tenants = append(tenants, tenant)
	}

	log.Debugf("tenants: %+v", tenants)
	return tenants
}

func kibanaXpackSavedObjectBulkCreate(data []Space) {
	client := resty.New().
		SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: !viper.GetBool("kibana.ssl_certificate_verification"),
		})

	if viper.GetBool("kibana.auth.enable") {
		client.SetBasicAuth(viper.GetString("kibana.auth.username"), viper.GetString("kibana.auth.password"))
	}

	for _, space := range data {
		savedObjectPayloads := []SavedObjectPayload{}
		for _, index := range space.Indices {
			savedObjectPayload := SavedObjectPayload{}
			savedObjectPayload.Type = "index-pattern"
			savedObjectPayload.ID = index
			savedObjectPayload.Attributes.TimeFieldName = space.Timestamp
			savedObjectPayload.Attributes.Title = index + "*"
			savedObjectPayloads = append(savedObjectPayloads, savedObjectPayload)
		}

		// do http request
		url := viper.GetString("kibana.host") + "/api/saved_objects/_bulk_create"

		if space.Name != "global" {
			url = viper.GetString("kibana.host") + "/s/" + space.Name + "/api/saved_objects/_bulk_create"
		}

		log.Infof("Space: %s", space.Name)
		log.Debugf("Http header: %s", "kbn-xsrf:true")
		log.Debugf("Http payload: %+v", savedObjectPayloads)

		result := SavedObjectResponse{}

		resp, err := client.R().
			SetHeader("kbn-xsrf", "true").
			SetBody(savedObjectPayloads).
			SetResult(&result).
			Post(url)

		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}

		if resp.StatusCode() != 200 {
			log.Errorf("An error occured when create saved objects, got http %d status code", resp.StatusCode())
			log.Errorf("Response: %s", resp.Body())
			os.Exit(1)
		}

		soConflictList := []string{}
		for _, v := range result.SavedObjects {
			if (v.Error != Error{}) {
				if v.Error.StatusCode == 409 {
					soConflictList = append(soConflictList, v.ID)
				}
			}
		}

		if len(soConflictList) != 0 {
			log.Warnf("Saved object %s conflict", soConflictList)
		}

		log.Info("Index pattern has been updated")
	}
}

func kibanaOpendistroSavedObjectBulkCreate(data []Tenant) {
	url := viper.GetString("kibana.host") + "/api/saved_objects/_bulk_create"

	client := resty.New().
		SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: !viper.GetBool("kibana.ssl_certificate_verification"),
		})

	if viper.GetBool("kibana.auth.enable") {
		client.SetBasicAuth(viper.GetString("kibana.auth.username"), viper.GetString("kibana.auth.password"))
	}

	for _, tenant := range data {
		savedObjectPayloads := []SavedObjectPayload{}
		for _, index := range tenant.Indices {
			savedObjectPayload := SavedObjectPayload{}
			savedObjectPayload.Type = "index-pattern"
			savedObjectPayload.ID = index
			savedObjectPayload.Attributes.TimeFieldName = tenant.Timestamp
			savedObjectPayload.Attributes.Title = index + "*"
			savedObjectPayloads = append(savedObjectPayloads, savedObjectPayload)
		}

		// do http request
		log.Infof("Tenant: %s", tenant.Name)
		log.Debugf("Http header: %s", "kbn-xsrf:true")
		log.Debugf("Http header: %s", "securitytenant:"+tenant.Name)
		log.Debugf("Http payload: %+v", savedObjectPayloads)

		result := SavedObjectResponse{}

		resp, err := client.R().
			SetHeader("kbn-xsrf", "true").
			SetHeader("securitytenant", tenant.Name).
			SetBody(savedObjectPayloads).
			SetResult(&result).
			Post(url)

		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}

		if resp.StatusCode() != 200 {
			log.Errorf("An error occured when create saved objects, got http %d status code", resp.StatusCode())
			log.Errorf("Response: %s", resp.Body())
			os.Exit(1)
		}

		soConflictList := []string{}
		for _, v := range result.SavedObjects {
			if (v.Error != Error{}) {
				if v.Error.StatusCode == 409 {
					soConflictList = append(soConflictList, v.ID)
				}
			}
		}

		if len(soConflictList) != 0 {
			log.Warnf("Saved object %s conflict", soConflictList)
		}

		log.Info("Index pattern has been updated")
	}
}
