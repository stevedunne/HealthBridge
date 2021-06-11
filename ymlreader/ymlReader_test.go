package ymlreader

import (
	"fmt"
	"healthBridge/assert"
	"testing"
)

func _TestWriteConf(t *testing.T) {

	res, err := writeConf()

	if err != nil {
		t.Errorf("Ah fuck %v", err)
	}

	if res != "" {
		t.Logf("Success: \r\n%s\r\n", res)
	} else {
		t.Errorf("Expected 20 but got %v", res)
	}
}

// func TestBytes(t *testing.T) {
// 	sample := "global:\r\n  scrape_interval: 15\r\nhealth_configs:\r\n"
// 	//     "- monitor_name: 'serverpings'"
// 	//     "monitor_item:"
// 	//    "" uri:
// 	//     - 'server1/ping'
// 	//     type: 'ping'"

// 	var sampleBuffer = []byte(sample)

// 	fileBuffer, err := ioutil.ReadFile("C:\\src\\HealthBridge\\test.yml")
// 	if err != nil {
// 		fmt.Print(fmt.Errorf("Error %v", err))
// 	}

// 	assert.IntEqual(len(sampleBuffer), len(fileBuffer), "", t)

// 	for i, val := range sampleBuffer {
// 		fmt.Printf("%v : %v (%s) %v (%s)\r\n", i, val, string(val), fileBuffer[i], string(fileBuffer[i]))
// 	}

// }

func _TestSimpleGlobalSample(t *testing.T) {
	sample := "global:\r\n  scrape_interval: 15\r\n"

	var sampleBuffer = []byte(sample)

	config, err := unmarshalFileContent(sampleBuffer)
	if err != nil {
		fmt.Print(fmt.Errorf("Error %v", err))
	}

	assert.NotNil(config, "", t)

	if config.GlobalConfig == (GlobalConfig{}) {
		t.Errorf("Expected GlobalConfig not nil")
	} else {
		t.Log("Passed")
	}
	assert.IntEqual(15, int(config.GlobalConfig.ScrapeInterval), t)
	assert.IntEqual(0, int(config.GlobalConfig.ScrapeTimeout), t)
}

func _TestUnmarshallFileContents(t *testing.T) {
	sample := "global:\r\n    scrape_interval: 15 \r\nhealth_configs:\r\n    - monitor_name: 'serverpings'\r\r    monitor_item:\r\n    uri:\r\n    - 'server1/ping'\r\n    type: 'ping'"

	var buffer = []byte(sample)

	config, err := unmarshalFileContent(buffer)

	assert.Error(err, t)
	// assert.NotNil(config, "", t)
	// assert.NotNil(config.GlobalConfig, "", t)
	assert.IntEqual(0, int(config.GlobalConfig.ScrapeInterval), t)
	assert.IntEqual(0, int(config.GlobalConfig.ScrapeTimeout), t)

	//	assert.NotNil(config.HealthConfigs, "", t)
	assert.IntEqual(1, len(config.HealthConfigs), t)
	assert.StrEqual("serverpings", config.HealthConfigs[0].MonitorName, "", t)
	assert.StrEqual("ping", config.HealthConfigs[0].MonitorType.Type, "", t)
	assert.IntEqual(1, len(config.HealthConfigs[0].MonitorType.URI), t)
	assert.StrEqual("server1/ping", config.HealthConfigs[0].MonitorType.URI[0], "", t)
}

func _TestMarshallFileContents(t *testing.T) {

	config, err := writeConf()
	if err != nil {
		fmt.Print(fmt.Errorf("Error %v", err))
	}

	fmt.Printf("%v", config)
}
