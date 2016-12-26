package main

import (
	"encoding/json"
	"fmt"
	"github.com/dcodix/grigori/lib/domain/message"
	"github.com/dcodix/grigori/lib/infrastructure/processor/processor_plain_to_json"
	"testing"
)

func TestReaderProcessorPlainToJson(test *testing.T) {
	messagesReaded := make(chan message.Message, 1000)
	messages := make(chan message.Message, 1000)

	processor := new(processor_plain_to_json.ProcessorPlainToJson)
	tags := []interface{}{"tagtet1", "tagtest2"}
	file := "thefile"
	log_type := "thetype"
	config := make(map[string]interface{})
	config["tags"] = tags
	config["resource"] = file
	config["type"] = log_type
	processor.Config(config)

	go func() {
		processor.Process(messagesReaded, messages)

	}()

	var messageSend message.Message
	var messageReaded message.Message
	messageSend.Message = "test"

	messagesReaded <- messageSend
	messageReaded = <-messages

	messageToWriteObject := new(processor_plain_to_json.JsonMessage)
	err := json.Unmarshal([]byte(messageReaded.Message), &messageToWriteObject)
	if err != nil {
		fmt.Println(err)
	}

	if messageToWriteObject.Message != messageSend.Message {
		test.Error("Message processed different than expected (message).")
	}

	if messageToWriteObject.Tags[0] != tags[0] {
		test.Error("Message processed different than expected (tags).")
	}

	if messageToWriteObject.File != file {
		test.Error("Message processed different than expected (tags).")
	}

	if messageToWriteObject.Type != log_type {
		test.Error("Message processed different than expected (tags).")
	}

	close(messagesReaded)
	close(messages)
}
