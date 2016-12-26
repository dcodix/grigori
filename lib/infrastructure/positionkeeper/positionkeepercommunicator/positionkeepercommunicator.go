package positionkeepercommunicator

import (
	"github.com/dcodix/grigori/lib/domain/positionkeeper"
	"github.com/dcodix/grigori/lib/domain/positionkeepermessage"
	"time"
)

type PositionKeeperComunicator struct {
	mainChannel     chan positionkeepermessage.PositionKeeperMessage
	responseChannel chan positionkeepermessage.PositionKeeperMessage
}

func (p *PositionKeeperComunicator) SetMainChannel(channel chan positionkeepermessage.PositionKeeperMessage) {
	p.mainChannel = channel
}

func (p *PositionKeeperComunicator) SetResponseChannel(channel chan positionkeepermessage.PositionKeeperMessage) {
	p.responseChannel = channel
}

func (p *PositionKeeperComunicator) CreateResponseCommunicator() positionkeeper.PositionKeeperComunicator {
	newCommunicator := new(PositionKeeperComunicator)
	responseChannel := make(chan positionkeepermessage.PositionKeeperMessage, 10)
	newCommunicator.SetMainChannel(p.responseChannel)
	newCommunicator.SetResponseChannel(responseChannel)
	return newCommunicator
}

func (p *PositionKeeperComunicator) GetMainChannel() (channel chan positionkeepermessage.PositionKeeperMessage) {
	return p.mainChannel
}

func (p *PositionKeeperComunicator) GetResponseChannel() (channel chan positionkeepermessage.PositionKeeperMessage) {
	return p.responseChannel
}

func (p *PositionKeeperComunicator) GetMessage() (message positionkeepermessage.PositionKeeperMessage, notMessage bool) {
	notMessage = false

	timeout := time.NewTicker(time.Duration(2) * time.Second)
	defer timeout.Stop()

	select {
	case messageReceived := <-p.responseChannel:
		message = messageReceived
	case <-timeout.C:
		notMessage = true
	}

	return message, notMessage
}

func (p *PositionKeeperComunicator) SendMessage(message positionkeepermessage.PositionKeeperMessage, sendChannel chan positionkeepermessage.PositionKeeperMessage) {
	sendChannel <- message
}
