package main

import (
	"github.com/dcodix/grigori/lib/domain/positionkeepermessage"
	"github.com/dcodix/grigori/lib/infrastructure/positionkeeper/positionkeepercommunicator"
	"github.com/dcodix/grigori/lib/infrastructure/positionkeeper/positionkeepermessageconstructor"
	"github.com/dcodix/grigori/lib/infrastructure/positionkeeper/positionkeeperrepository"
	"testing"
)

//Dummy position keeper for testing purposes
type TestBlankPositionKeeper struct{}

func (t *TestBlankPositionKeeper) Save(resource string, position int64, saveToBackend bool) {}

func (t *TestBlankPositionKeeper) Load(resource string) (int64, bool) {
	return int64(873), false
}

func (t *TestBlankPositionKeeper) Config(config map[string]string) {}

//TEST
func TestPositionKeepperRepository(test *testing.T) {
	constructor := new(positionkeepermessageconstructor.PositionKeeperMessageConstructor)
	communicatorPositionKeeper := new(positionkeepercommunicator.PositionKeeperComunicator)
	positionkeeper := new(TestBlankPositionKeeper)
	repository := new(positionkeeperrepository.PositionKeeperRepository)
	channelPositionKeeper := make(chan positionkeepermessage.PositionKeeperMessage, 10)
	quit := make(chan bool, 1)

	communicatorPositionKeeper.SetMainChannel(channelPositionKeeper)
	communicatorPositionKeeper.SetResponseChannel(channelPositionKeeper)

	communicatorClient := communicatorPositionKeeper.CreateResponseCommunicator()

	repository.SetConstructor(constructor)
	repository.SetCommunicator(communicatorPositionKeeper)
	repository.SetPositionKeeper(positionkeeper)

	go func() {
		repository.Run(quit)
	}()

	//Test Save
	position := int64(1)
	resource := "testresource"
	action := "ack"
	saveMessage := constructor.ConstructSavePostitionMessage(resource, position, communicatorClient.GetResponseChannel())
	communicatorClient.SendMessage(saveMessage, channelPositionKeeper)
	ackMessage, empty := communicatorClient.GetMessage()
	if empty || ackMessage.Action != action || ackMessage.Path != resource || ackMessage.Channel != communicatorClient.GetResponseChannel() {
		//	test.Error("Ack message not correctly created or lost.")
	}

	//Test Load
	resource = "testresource2"
	action = "response"
	loadMessage := constructor.ConstructLoadPostitionMessage(resource, communicatorClient.GetResponseChannel())
	communicatorClient.SendMessage(loadMessage, channelPositionKeeper)
	responseMessage, empty := communicatorClient.GetMessage()
	if empty || responseMessage.Action != action || responseMessage.Path != resource || responseMessage.Channel != communicatorClient.GetResponseChannel() {
		test.Error("Response message not correctly created or lost.")
	}

	quit <- true
}
