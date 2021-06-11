package ymlreader

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

// Config holds the global defaults and an array of health configs
// https://github.com/prometheus/prometheus/blob/0a8912433a457b54a871df75e72afc3e45a31d5c/config/config.go
type Config struct {
	GlobalConfig  GlobalConfig    `yaml:"global"`
	HealthConfigs []*HealthConfig `yaml:"health_configs,omitempty"`
}

//GlobalConfig represents global settings for the health monitors
type GlobalConfig struct {
	// How frequently to scrape targets by default.
	ScrapeInterval int32 `yaml:"scrape_interval,omitempty"`
	// The default timeout when scraping targets.
	ScrapeTimeout int32 `yaml:"scrape_timeout,omitempty"`
}

//HealthConfig represents an individual monitors settings
type HealthConfig struct {
	// The job name to which the job label is set by default.
	MonitorName string `yaml:"monitor_name"`

	MonitorType MonitorItem `yaml:"monitor_item"`
}

//MonitorItem contains the health monitor type and an array of the Uris
type MonitorItem struct {
	URI  []string `yaml:"uri"`
	Type string   `yaml:"type"`
}

// writeConf creates a sample configuration and exports it as yml
func writeConf() (string, error) {
	item := HealthConfig{}
	item.MonitorName = "test monitor name"
	item.MonitorType.Type = "ping"
	var uris [2]string
	uris[0] = "test"
	uris[1] = "other"
	item.MonitorType.URI = uris[:]

	item2 := HealthConfig{}
	item2.MonitorName = "othername"
	item2.MonitorType.Type = "ping"
	var uris2 [1]string
	uris2[0] = "first"
	item2.MonitorType.URI = uris2[:]

	var configs [2]*HealthConfig
	configs[0] = &item
	configs[1] = &item2

	config := Config{}
	config.GlobalConfig.ScrapeInterval = 10
	config.GlobalConfig.ScrapeTimeout = 60

	config.HealthConfigs = configs[:]

	data, err := yaml.Marshal(config)

	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
	} else {
		return string(data), nil
	}
	return "", nil
}

// ReadConf takes a filename and reads the settings
func ReadConf(filename string) (*Config, error) {
	fmt.Printf("Reading file: %s", filename)

	buf, err := getFileContent(filename)
	if err != nil {
		return nil, err
	}

	fmt.Println("Read file contents")

	c, err := unmarshalFileContent(buf)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func getFileContent(filename string) ([]byte, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func unmarshalFileContent(buffer []byte) (*Config, error) {
	c := &Config{}
	err := yaml.Unmarshal(buffer, c)
	if err != nil {
		return nil, fmt.Errorf("in decoding: %v", err)
	}
	return c, nil
}
