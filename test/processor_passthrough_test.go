package main

import (
	"github.com/dcodix/grigori/lib/domain/message"
	"github.com/dcodix/grigori/lib/infrastructure/processor/processor_passthrough"
	"testing"
)

func TestReaderProcessorPassthrough(test *testing.T) {
	messagesReaded := make(chan message.Message, 1000)
	messages := make(chan message.Message, 1000)

	processor := new(processor_passthrough.ProcessorPassthrough)

	go func() {
		processor.Process(messagesReaded, messages)

	}()

	var messageSend message.Message
	var messageReaded message.Message
	messageSend.Message = "test"

	messagesReaded <- messageSend
	messageReaded = <-messages

	if messageSend != messageReaded {
		test.Error("Message processed different than sended.")
	}

	close(messagesReaded)
	close(messages)
}
