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

type File struct {
	Enable bool   `mapstructure:"enable"`
	Path   string `mapstructure:"path"`
}

type Log struct {
	Debug bool `mapstructure:"debug"`
	File  File `mapstructure:"file"`
}

type Auth struct {
	Enable   bool   `mapstructure:"enable"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Index struct {
	Name            string `mapstructure:"name"`
	RolloverPattern string `mapstructure:"rollover_pattern" yaml:"rollover_pattern"`
}

type Elasticsearch struct {
	Host                       string  `mapstructure:"host"`
	SslCertificateVerification bool    `mapstructure:"ssl_certificate_verification" yaml:"ssl_certificate_verification"`
	Auth                       Auth    `mapstructure:"auth"`
	Indices                    []Index `mapstructure:"indices"`
}

type Kibana struct {
	Host                       string `mapstructure:"host"`
	SslCertificateVerification bool   `mapstructure:"ssl_certificate_verification" yaml:"ssl_certificate_verification"`
	Auth                       Auth   `mapstructure:"auth"`
}

type Tenant struct {
	Name      string `mapstructure:"name"`
	Pattern   string `mapstructure:"pattern"`
	Timestamp string `mapstructure:"timestamp"`
}

type Opendistro struct {
	Enable  bool     `mapstructure:"enable"`
	Tenants []Tenant `mapstructure:"tenants"`
}

type Space struct {
	Name      string `mapstructure:"name"`
	Pattern   string `mapstructure:"pattern"`
	Timestamp string `mapstructure:"timestamp"`
}

type Xpack struct {
	Enable bool    `mapstructure:"enable"`
	Spaces []Space `mapstructure:"spaces"`
}

type Config struct {
	Interval      string        `mapstructure:"interval"`
	Log           Log           `mapstructure:"log"`
	Elasticsearch Elasticsearch `mapstructure:"elasticsearch"`
	Kibana        Kibana        `mapstructure:"kibana"`
	Xpack         Xpack         `mapstructure:"xpack"`
	Opendistro    Opendistro    `mapstructure:"opendistro"`
}
