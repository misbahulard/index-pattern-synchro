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
	"os"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Run() {
	job()

	interval := "@every " + viper.GetString("interval")

	log.Infof("Run the job every %s", viper.GetString("interval"))

	c := cron.New()
	c.AddFunc(interval, job)
	c.Run()
}

func job() {
	log.Info("Get Elasticsearch indices")
	esIndex := getEsIndex()

	if (viper.GetBool("xpack.enable") && viper.GetBool("opendistro.enable")) || (!viper.GetBool("xpack.enable") && !viper.GetBool("opendistro.enable")) {
		log.Error("You need to enable one of Kibana Xpack Spaces or Kibana Opendistro Tenants, please choose one of them.")
		os.Exit(1)
	}

	// log.Debug(esIndex)

	kibanaSpaces := []Space{}
	if viper.GetBool("xpack.enable") {
		log.Info("Get Kibana spaces")
		kibanaSpaces = getKibanaSpaces(esIndex)
	}

	kibanaXpackSavedObjectBulkCreate(kibanaSpaces)

	kibanaTenants := []Tenant{}
	if viper.GetBool("opendistro.enable") {
		log.Info("Get Kibana tenants")
		kibanaTenants = getKibanaTenants(esIndex)
	}

	kibanaOpendistroSavedObjectBulkCreate(kibanaTenants)
}
