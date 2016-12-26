package main

import (
	"github.com/dcodix/grigori/lib"
	"github.com/dcodix/grigori/lib/domain/config"
	"github.com/dcodix/grigori/lib/domain/message"
	"github.com/dcodix/grigori/lib/domain/positionkeeper"
	"github.com/dcodix/grigori/lib/domain/positionkeepermessage"
	"github.com/dcodix/grigori/lib/domain/processor"
	"github.com/dcodix/grigori/lib/domain/reader"
	"github.com/dcodix/grigori/lib/domain/writer"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"testing"
)

//DUMMY POSITION KEEPER
type testPositionKeeper struct{}

func (p *testPositionKeeper) Load(resource string) (position int64, notFound bool) {
	return int64(1), false
}

func (p *testPositionKeeper) Save(resource string, position int64, saveToBackend bool) {}

func (p *testPositionKeeper) Config(config map[string]string) {}

//DUMMY REPOSITORY
type testPositionKeeperRepository struct{}

func (p *testPositionKeeperRepository) Run(quit chan bool) {
	for {
		select {
		case <-quit:
			return
		default:
		}
	}

}
func (p *testPositionKeeperRepository) SetConstructor(constructor positionkeeper.PositionKeeperMessageConstructor) {
}

func (p *testPositionKeeperRepository) SetCommunicator(communicator positionkeeper.PositionKeeperComunicator) {
}

func (p *testPositionKeeperRepository) SetPositionKeeper(positionkeeper positionkeeper.PositionKeeper) {
}

//DUMMY MESSAGE POSITION KEEPER MESSAGE CONSTRUCTOR
type testPositionKeeperMessageConstructor struct{}

func (p *testPositionKeeperMessageConstructor) ConstructSavePostitionMessage(resource string, position int64, responsechannel chan positionkeepermessage.PositionKeeperMessage) positionkeepermessage.PositionKeeperMessage {
	var positionKeeperMessage positionkeepermessage.PositionKeeperMessage
	return positionKeeperMessage
}

func (p *testPositionKeeperMessageConstructor) ConstructLoadPostitionMessage(resource string, responsechannel chan positionkeepermessage.PositionKeeperMessage) positionkeepermessage.PositionKeeperMessage {
	var positionKeeperMessage positionkeepermessage.PositionKeeperMessage
	return positionKeeperMessage
}

func (p *testPositionKeeperMessageConstructor) ConstructResponsePostitionMessage(resource string, position int64, responsechannel chan positionkeepermessage.PositionKeeperMessage) positionkeepermessage.PositionKeeperMessage {
	var positionKeeperMessage positionkeepermessage.PositionKeeperMessage
	return positionKeeperMessage
}

func (p *testPositionKeeperMessageConstructor) ConstructAckPostitionMessage(resource string, responsechannel chan positionkeepermessage.PositionKeeperMessage) positionkeepermessage.PositionKeeperMessage {
	var positionKeeperMessage positionkeepermessage.PositionKeeperMessage
	return positionKeeperMessage
}

func (p *testPositionKeeperMessageConstructor) ExtractPosition(positionKeeperMessage positionkeepermessage.PositionKeeperMessage) int64 {
	return int64(1)
}

func (p *testPositionKeeperMessageConstructor) ExtractResource(positionKeeperMessage positionkeepermessage.PositionKeeperMessage) string {
	return "test"
}

func (p *testPositionKeeperMessageConstructor) ExtractChannel(positionKeeperMessage positionkeepermessage.PositionKeeperMessage) chan positionkeepermessage.PositionKeeperMessage {
	channel := make(chan positionkeepermessage.PositionKeeperMessage, 10)
	return channel
}

func (p *testPositionKeeperMessageConstructor) ExtractAction(positionKeeperMessage positionkeepermessage.PositionKeeperMessage) string {
	return "test"
}

//DUMMY COMMUNICATOR
type testPositionKeeperComunicator struct{}

func (p *testPositionKeeperComunicator) SetMainChannel(channel chan positionkeepermessage.PositionKeeperMessage) {
}

func (p *testPositionKeeperComunicator) SetResponseChannel(channel chan positionkeepermessage.PositionKeeperMessage) {
}

func (p *testPositionKeeperComunicator) CreateResponseCommunicator() positionkeeper.PositionKeeperComunicator {
	newCommunicator := new(testPositionKeeperComunicator)
	return newCommunicator
}

func (p *testPositionKeeperComunicator) GetMainChannel() (channel chan positionkeepermessage.PositionKeeperMessage) {
	channel = make(chan positionkeepermessage.PositionKeeperMessage, 10)
	return channel
}

func (p *testPositionKeeperComunicator) GetResponseChannel() (channel chan positionkeepermessage.PositionKeeperMessage) {
	channel = make(chan positionkeepermessage.PositionKeeperMessage, 10)
	return channel
}

func (p *testPositionKeeperComunicator) GetMessage() (positionkeepermessage.PositionKeeperMessage, bool) {
	var positionKeeperMessage positionkeepermessage.PositionKeeperMessage
	return positionKeeperMessage, false
}

func (p *testPositionKeeperComunicator) SendMessage(message positionkeepermessage.PositionKeeperMessage, channel chan positionkeepermessage.PositionKeeperMessage) {
}

//DUMMY READER
type testReader struct{}

func (k *testReader) Read(messages chan message.Message, quit chan bool) {
	b := CreateInMessagesArray(400)
	for _, message := range b {
		messages <- message
		select {
		case <-quit:
			return
		default:

		}
	}
}

func (k *testReader) Config(config map[string]interface{}) {
}

func (k *testReader) SetPositionKeeperComunicator(communicator positionkeeper.PositionKeeperComunicator) {
}

func (k *testReader) SetPPositionKeeperMessageConstructor(constructor positionkeeper.PositionKeeperMessageConstructor) {
}

func (k *testReader) Clone() reader.Reader {
	clonedReader := new(testReader)
	return clonedReader
}

func (k *testReader) GetResource() string {
	return "none"
}

//DUMMY PROCESSOR
type testProcessor struct {}

func (p *testProcessor) Process(messagesIn chan message.Message, messagesOut chan message.Message) {

	for message := range messagesIn {
		messagesOut <- message
	}
}

func (p *testProcessor) Config(map[string]interface{}) {}

func (p *testProcessor) Clone() processor.Processor {
	clonedProcessor := new(testProcessor)
	return clonedProcessor

}

//DUMMY WRITER
type testWriter struct {
	data []message.Message
}

func (w *testWriter) Config(config map[string]interface{}) {}

func (w *testWriter) Clone() writer.Writer {
	clonedWriter := new(testWriter)
	return clonedWriter
}

func (w *testWriter) Write(messages chan message.Message) {
	for message := range messages {
		w.data = append(w.data, message)
	}
}

//TEST
func CreateStringArray(max int) []string {
	var string_array []string
	for i := 0; i < max; i++ {
		string_array = append(string_array, strconv.Itoa(i))
	}
	return string_array

}

func CreateInMessagesArray(max int) []message.Message {
	var messages []message.Message
	for _, s := range CreateStringArray(max) {
		messages = append(messages, message.Message{s})
	}
	return messages
}

func CreateMessagesArray(max int) []message.Message {
	var messages []message.Message
	for _, s := range CreateStringArray(max) {
		messages = append(messages, message.Message{s})
	}
	return messages
}

func messageInSlice(message message.Message, messages []message.Message) bool {
	for _, i := range messages {
		if message == i {
			return true
		}
	}
	return false
}

func TestGoroutinesConfig(test *testing.T) {
	reader := new(testReader)
	perocessorTest := new(testProcessor)
	writerTest := new(testWriter)
	positionkeeper := new(testPositionKeeper)
	communicator := new(testPositionKeeperComunicator)
	constructor := new(testPositionKeeperMessageConstructor)
	repository := new(testPositionKeeperRepository)
	repository.SetCommunicator(communicator)
	repository.SetConstructor(constructor)
	repository.SetPositionKeeper(positionkeeper)

	configPosKeeper := make(map[string]string, 10)
	configPosKeeper["test"] = "test"
	positionkeeper.Config(configPosKeeper)

	var processors []processor.Processor
	processors = append(processors, perocessorTest)

	var writers []writer.Writer
	writers = append(writers, writerTest)

	resource := config.Resource{reader, processors}
	resources := []config.Resource{resource}
	var configReadedTest config.ConfigReadedFile
	con := &config.Config{resources, writers, repository, constructor, communicator, configReadedTest}

	//Listen for signals to stop app
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	//Run Grigori
	grigori.Run(con, sigs)

	//TEST
	testSlice := CreateMessagesArray(400)
	testSlicePointer := &testSlice

	writerOut, _ := con.Writers[0].(*testWriter)

	exist := true
	for _, i := range writerOut.data {
		if !messageInSlice(i, *testSlicePointer) {
			exist = false
		}
	}
	if !exist {
		test.Error("Writter did not receive expectet value.")
	}

	exist = true
	for _, i := range *testSlicePointer {
		if !messageInSlice(i, writerOut.data) {
			exist = false
		}
	}
	if !exist {
		test.Error("Writter did not receive expectet value.")
	}

}
