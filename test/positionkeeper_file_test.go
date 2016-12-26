package main

import (
	"fmt"
	"github.com/dcodix/grigori/lib/infrastructure/positionkeeper/positionkeeper_file"
	"os"
	"testing"
)

func TestPositionKeepperFile(test *testing.T) {
	positionkeeper_file_test := "/tmp/test_positionkeeper_839829383.tmp"
	positionkeeper := new(positionkeeper_file.PositionKeeperFile)
	config := make(map[string]string, 100)
	config["path"] = positionkeeper_file_test
	positionkeeper.Config(config)

	resource1 := "/path/to/testresource1"
	position1 := int64(82)
	resource2 := "/path/to/testresource2"
	position2 := int64(43)
	resource_not_saved := "/path/to/testresource3"

	positionkeeper.Save(resource1, position1, true)
	positionkeeper.Save(resource2, position2, true)

	readed_position1, notFound := positionkeeper.Load(resource1)
	if notFound || readed_position1 != position1 {
		test.Error("Save or Load not working fine. Retrieved position incorrect.")
	}
	readed_position2, notFound := positionkeeper.Load(resource2)
	if notFound || readed_position2 != position2 {
		test.Error("Save or Load not working fine. Retrieved position incorrect.")
	}

	_, notFound = positionkeeper.Load(resource_not_saved)
	if !notFound {
		test.Error("Not exit with error when resource don't exist")
	}

	err := os.Remove(positionkeeper_file_test)

	if err != nil {
		fmt.Println(err)
		return
	}

}
