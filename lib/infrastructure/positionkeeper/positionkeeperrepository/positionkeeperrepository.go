package positionkeeperrepository

import (
	"fmt"
	"github.com/dcodix/grigori/lib/domain/positionkeeper"
	"time"
)

type PositionKeeperRepository struct {
	constructor    positionkeeper.PositionKeeperMessageConstructor
	communicator   positionkeeper.PositionKeeperComunicator
	positionkeeper positionkeeper.PositionKeeper
}

func (p *PositionKeeperRepository) SetConstructor(constructor positionkeeper.PositionKeeperMessageConstructor) {
	p.constructor = constructor
}

func (p *PositionKeeperRepository) SetCommunicator(communicator positionkeeper.PositionKeeperComunicator) {
	p.communicator = communicator
}

func (p *PositionKeeperRepository) SetPositionKeeper(positionkeeper positionkeeper.PositionKeeper) {
	p.positionkeeper = positionkeeper
}

func saveToBackend(ticker *time.Ticker) (saveToBAckend bool) {
	select {
	case <-ticker.C:
		saveToBAckend = true
	default:
		saveToBAckend = false
	}
	return saveToBAckend
}

func (p *PositionKeeperRepository) Run(quit chan bool) {
	ticker := time.NewTicker(time.Duration(5) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-quit:
			return
		default:

		}

		saveToBackend := saveToBackend(ticker)

		msg, empty := p.communicator.GetMessage()
		if !empty {
			resource := p.constructor.ExtractResource(msg)
			channel := p.constructor.ExtractChannel(msg)
			switch msg.Action {
			case "save":
				position := p.constructor.ExtractPosition(msg)
				p.positionkeeper.Save(resource, position, saveToBackend)
			case "load":
				position, _ := p.positionkeeper.Load(resource)
				responseMessage := p.constructor.ConstructResponsePostitionMessage(resource, position, p.constructor.ExtractChannel(msg))
				p.communicator.SendMessage(responseMessage, channel)
			default:
				fmt.Println("do nothing")
			}
		}

	}

}
