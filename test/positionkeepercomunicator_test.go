package main

import (
	"github.com/dcodix/grigori/lib/domain/positionkeepermessage"
	"github.com/dcodix/grigori/lib/infrastructure/positionkeeper/positionkeepercommunicator"
	"testing"
)

func TestPositionKeepperComunicator(test *testing.T) {
	channel := make(chan positionkeepermessage.PositionKeeperMessage, 10)
	var message positionkeepermessage.PositionKeeperMessage
	communicator := new(positionkeepercommunicator.PositionKeeperComunicator)

	//Tests gets/sets
	communicator.SetMainChannel(channel)
	communicator.SetResponseChannel(channel)

	if communicator.GetMainChannel() != channel {
		test.Error("Main channel incorrecly set or get.")
	}
	if communicator.GetResponseChannel() != channel {
		test.Error("Response channel incorrecly set or get.")
	}

	//Test Success
	message.Action = "testaction1"

	communicator.SendMessage(message, channel)
	messageReceived, notMessage := communicator.GetMessage()
	if notMessage || messageReceived.Action != message.Action {
		test.Error("Message not send and receiver correctly.")
	}

	//Test Success getting channel from communicator
	message.Action = "testaction2"

	communicator.SendMessage(message, communicator.GetMainChannel())
	messageReceived, notMessage = communicator.GetMessage()
	if notMessage || messageReceived.Action != message.Action {
		test.Error("Message not send and receiver correctly using GetPostitionKeeperMainChannel function.")
	}

	//Test fail
	message.Action = "testaction3"

	communicator.SendMessage(message, communicator.GetMainChannel())
	messageReceived, notMessage = communicator.GetMessage()
	if notMessage || messageReceived.Action == "testaction4" {
		test.Error("Message received as correct when it should not be.")
	}

	//Test timeout
	messageReceived, notMessage = communicator.GetMessage()
	if !notMessage {
		test.Error("Timeout did not work.")
	}

	//Test communicator constructor
	clientCommunicator := communicator.CreateResponseCommunicator()

	if clientCommunicator.GetMainChannel() != channel {
		test.Error("Main channel incorrectly set to new communicator")
	}
	if clientCommunicator.GetResponseChannel() == channel {
		test.Error("Response channel incorrectly set to new communicator")
	}

}
