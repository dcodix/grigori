package main

import (
	"github.com/dcodix/grigori/lib/domain/positionkeepermessage"
	"github.com/dcodix/grigori/lib/infrastructure/positionkeeper/positionkeepermessageconstructor"
	"testing"
)

func TestPositionKeepperConstructor(test *testing.T) {
	//TEST CONSTRUCT
	//Test construct save
	constructor := new(positionkeepermessageconstructor.PositionKeeperMessageConstructor)
	responsechannel := make(chan positionkeepermessage.PositionKeeperMessage, 10)
	position := int64(1764)
	resource := "testresource"
	constructedSaveMessage := constructor.ConstructSavePostitionMessage(resource, position, responsechannel)

	if constructedSaveMessage.Action != "save" {
		test.Error("Save message not well constructed.")
	}
	if constructedSaveMessage.Path != resource {
		test.Error("Save message not well constructed.")
	}
	if constructedSaveMessage.Position != position {
		test.Error("Save message not well constructed.")
	}
	if constructedSaveMessage.Channel != responsechannel {
		test.Error("Save message not well constructed.")
	}

	//Test construct load
	constructor = new(positionkeepermessageconstructor.PositionKeeperMessageConstructor)
	responsechannel = make(chan positionkeepermessage.PositionKeeperMessage, 10)
	resource = "testresource"
	constructedLoadMessage := constructor.ConstructLoadPostitionMessage(resource, responsechannel)
	if constructedLoadMessage.Action != "load" {
		test.Error("Load message not well constructed.")
	}
	if constructedLoadMessage.Path != resource {
		test.Error("Load message not well constructed.")
	}
	if constructedLoadMessage.Channel != responsechannel {
		test.Error("Load message not well constructed.")
	}

	//Test construct position response
	constructor = new(positionkeepermessageconstructor.PositionKeeperMessageConstructor)
	responsechannel = make(chan positionkeepermessage.PositionKeeperMessage, 10)
	position = int64(17984)
	resource = "testresource2"
	constructedResponsePositionMessage := constructor.ConstructResponsePostitionMessage(resource, position, responsechannel)

	if constructedResponsePositionMessage.Action != "response" {
		test.Error("Position Response message not well constructed.")
	}
	if constructedResponsePositionMessage.Path != resource {
		test.Error("Position Response message not well constructed.")
	}
	if constructedResponsePositionMessage.Position != position {
		test.Error("Position Response message not well constructed.")
	}
	if constructedResponsePositionMessage.Channel != responsechannel {
		test.Error("Position Response message not well constructed.")
	}

	//Test construct ack
	constructor = new(positionkeepermessageconstructor.PositionKeeperMessageConstructor)
	responsechannel = make(chan positionkeepermessage.PositionKeeperMessage, 10)
	resource = "testresource3"
	constructedAckMessage := constructor.ConstructAckPostitionMessage(resource, responsechannel)
	if constructedAckMessage.Action != "ack" {
		test.Error("Ack message not well constructed.")
	}
	if constructedAckMessage.Path != resource {
		test.Error("Ack message not well constructed.")
	}
	if constructedAckMessage.Channel != responsechannel {
		test.Error("Ack message not well constructed.")
	}

	//TEST EXTRACT
	constructor = new(positionkeepermessageconstructor.PositionKeeperMessageConstructor)
	responsechannel = make(chan positionkeepermessage.PositionKeeperMessage, 10)
	position = int64(178764)
	resource = "testresource4"
	constructedSaveMessageTestExtract := constructor.ConstructSavePostitionMessage(resource, position, responsechannel)

	//Test Extract position
	if constructor.ExtractPosition(constructedSaveMessageTestExtract) != position {
		test.Error("Error extracting position.")
	}
	//Test Extract resource
	if constructor.ExtractResource(constructedSaveMessageTestExtract) != resource {
		test.Error("Error extracting resource.")
	}
	//Test Extract channel
	if constructor.ExtractChannel(constructedSaveMessageTestExtract) != responsechannel {
		test.Error("Error extracting channel.")
	}
	//Test Extract action
	if constructor.ExtractAction(constructedSaveMessageTestExtract) != "save" {
		test.Error("Error extracting action.")
	}

}
