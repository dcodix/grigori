package positionkeeper

import (
	"github.com/dcodix/grigori/lib/domain/positionkeepermessage"
)

type PositionKeeper interface {
	Load(resource string) (position int64, notFound bool)
	Save(resource string, position int64, saveToBAckend bool)
	Config(config map[string]string)
}

type PositionKeeperRepository interface {
	SetConstructor(constructor PositionKeeperMessageConstructor)
	SetCommunicator(communicator PositionKeeperComunicator)
	SetPositionKeeper(positionkeeper PositionKeeper)
	Run(quit chan bool)
}

type PositionKeeperComunicator interface {
	SetMainChannel(channel chan positionkeepermessage.PositionKeeperMessage)
	SetResponseChannel(channel chan positionkeepermessage.PositionKeeperMessage)
	CreateResponseCommunicator() PositionKeeperComunicator
	GetMainChannel() (channel chan positionkeepermessage.PositionKeeperMessage)
	GetResponseChannel() (channel chan positionkeepermessage.PositionKeeperMessage)
	GetMessage() (positionkeepermessage.PositionKeeperMessage, bool)
	SendMessage(message positionkeepermessage.PositionKeeperMessage, channel chan positionkeepermessage.PositionKeeperMessage)
}

type PositionKeeperMessageConstructor interface {
	ConstructSavePostitionMessage(resource string, position int64, responsechannel chan positionkeepermessage.PositionKeeperMessage) positionkeepermessage.PositionKeeperMessage
	ConstructLoadPostitionMessage(resource string, responsechannel chan positionkeepermessage.PositionKeeperMessage) positionkeepermessage.PositionKeeperMessage
	ConstructResponsePostitionMessage(resource string, position int64, responsechannel chan positionkeepermessage.PositionKeeperMessage) positionkeepermessage.PositionKeeperMessage
	ConstructAckPostitionMessage(resource string, responsechannel chan positionkeepermessage.PositionKeeperMessage) positionkeepermessage.PositionKeeperMessage
	ExtractPosition(positionKeeperMessage positionkeepermessage.PositionKeeperMessage) int64
	ExtractResource(positionKeeperMessage positionkeepermessage.PositionKeeperMessage) string
	ExtractChannel(positionKeeperMessage positionkeepermessage.PositionKeeperMessage) chan positionkeepermessage.PositionKeeperMessage
	ExtractAction(positionKeeperMessage positionkeepermessage.PositionKeeperMessage) string
}
