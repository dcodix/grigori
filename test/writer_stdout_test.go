package main

import (
	"github.com/dcodix/grigori/lib/domain/message"
	"github.com/dcodix/grigori/lib/infrastructure/writer/writer_stdout"
	"testing"
)

func TestReaderWriterStdout(test *testing.T) {
	messages := make(chan message.Message, 1000)

	var messageSend message.Message
	messageSend.Message = "test"

	writer := new(writer_stdout.WriterStdout)

	go func() {
		writer.Write(messages)

	}()

	messages <- messageSend
	close(messages)
}
