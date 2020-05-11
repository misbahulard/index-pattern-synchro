package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// Config : read yaml config file
func Config(cfg *Conf) {
	f, err := os.Open("config.yaml")

	if err != nil {
		log.Fatalln(err)
	}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
}
