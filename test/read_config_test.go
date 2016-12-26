package main

import (
	"encoding/json"
	"fmt"
	"github.com/dcodix/grigori/lib/domain/config"
	"github.com/dcodix/grigori/lib/infrastructure/config/read_config_file"
	"io/ioutil"
	"os"
	"testing"
)

func createTmpTestConfigFile(filepath string, testConfig config.ConfigReadedFile) {
	content, _ := json.Marshal(testConfig)
	err := ioutil.WriteFile(filepath, content, 0644)
	if err != nil {
		fmt.Printf("WriteFileJson ERROR: %+v", err)
	}
}

func deleteTestTmpFile(file string) {
	os.Remove(file)
}

func TestReadConfig(test *testing.T) {
	testConfigFile := "/tmp/test_grigori_config_file.cfg.tmp"

	testConfig := config.ConfigReadedFile{ConfigReaded: config.ConfigReaded{Writer: map[string]interface{}{"type": "testwriter", "n_witers": 1}, Resources: []map[string]interface{}{{"reader": "testreader", "processor": "testprocessor", "n_processors": "2", "resource": "testresource"}}}}
	testConfig2 := config.ConfigReadedFile{ConfigReaded: config.ConfigReaded{Writer: map[string]interface{}{"type": "differentwriter", "n_writers": 1}, Resources: []map[string]interface{}{{"reader": "testreader", "processor": "testprocessor", "n_processors": "2", "resource": "testresource"}}}}

	createTmpTestConfigFile(testConfigFile, testConfig)

	config := read_config_file.ReadConfig(testConfigFile)

	if config.Writer["type"] != testConfig.Writer["type"] || config.Writer["n_writers"] != testConfig.Writer["n_writers"] || config.Resources[0]["reader"] != testConfig.Resources[0]["reader"] || config.Resources[0]["processor"] != testConfig.Resources[0]["processor"] {
		test.Error("Config file not read correctly.")
	}
	if config.Writer["type"] == testConfig2.Writer["type"] {
		test.Error("Config file not read correctly.")
	}

	deleteTestTmpFile(testConfigFile)
}
