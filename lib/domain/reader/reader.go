package reader

import (
	"github.com/dcodix/grigori/lib/domain/message"
	"github.com/dcodix/grigori/lib/domain/positionkeeper"
)

type Reader interface {
	Config(config map[string]interface{})
	SetPositionKeeperComunicator(communicator positionkeeper.PositionKeeperComunicator)
	SetPPositionKeeperMessageConstructor(constructor positionkeeper.PositionKeeperMessageConstructor)
	Clone() (reader Reader)
	Read(messages chan message.Message, quit chan bool)
	GetResource() (resoure string)
}
