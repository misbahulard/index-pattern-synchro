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
package config

type ConfigFile struct {
	Enable bool   `yaml:"enable"`
	Path   string `yaml:"path"`
}

type ConfigLog struct {
	Debug bool       `yaml:"debug"`
	File  ConfigFile `yaml:"file"`
}

type ConfigAuth struct {
	Enable   bool   `yaml:"enable"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type ConfigElasticsearch struct {
	Host                       string     `yaml:"host"`
	SslCertificateVerification bool       `yaml:"ssl_certificate_verification"`
	Auth                       ConfigAuth `yaml:"auth"`
	Indices                    []string   `yaml:"indices"`
}

type ConfigKibana struct {
	Host                       string     `yaml:"host"`
	SslCertificateVerification bool       `yaml:"ssl_certificate_verification"`
	Auth                       ConfigAuth `yaml:"auth"`
}

type Config struct {
	Interval      string              `yaml:"interval"`
	Log           ConfigLog           `yaml:"log"`
	Elasticsearch ConfigElasticsearch `yaml:"elasticsearch"`
	Kibana        ConfigKibana        `yaml:"kibana"`
}
