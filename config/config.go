package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

type (
	Scraper struct {
		Sites   []string `yaml:"sites"`
		Baseurl string   `yaml:"baseurl"`
		Timeout int      `yaml:"timeout"`
	}
)

func LoadConfiguration(configFile string) (Scraper, error) {
	var config Scraper

	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		return config, errors.New(fmt.Sprintf(
			"Failed to read configuration file %s: %s",
			configFile,
			err,
		))
	}
	unMarshalError := yaml.Unmarshal(content, &config)
	if unMarshalError != nil {
		return config, errors.New(fmt.Sprintf(
			"Failed to load configuration file %s: %s",
			configFile,
			unMarshalError,
		))
	}

	log.Debugf("Config read for scraper - Sites: %q, Baseurl: %s, Timeout: %s",
		config.Sites,
		config.Baseurl,
		config.Timeout,
	)
	return config, nil
}
