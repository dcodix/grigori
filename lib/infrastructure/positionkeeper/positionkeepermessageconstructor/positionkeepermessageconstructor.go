package positionkeepermessageconstructor

import (
	"github.com/dcodix/grigori/lib/domain/positionkeepermessage"
)

type PositionKeeperMessageConstructor struct{}

func (p *PositionKeeperMessageConstructor) ConstructSavePostitionMessage(resource string, position int64, responsechannel chan positionkeepermessage.PositionKeeperMessage) positionkeepermessage.PositionKeeperMessage {
	var positionKeeperMessage positionkeepermessage.PositionKeeperMessage
	positionKeeperMessage.Action = "save"
	positionKeeperMessage.Path = resource
	positionKeeperMessage.Position = position
	positionKeeperMessage.Channel = responsechannel
	return positionKeeperMessage
}

func (p *PositionKeeperMessageConstructor) ConstructLoadPostitionMessage(resource string, responsechannel chan positionkeepermessage.PositionKeeperMessage) positionkeepermessage.PositionKeeperMessage {
	var positionKeeperMessage positionkeepermessage.PositionKeeperMessage
	positionKeeperMessage.Action = "load"
	positionKeeperMessage.Path = resource
	positionKeeperMessage.Channel = responsechannel
	return positionKeeperMessage
}

func (p *PositionKeeperMessageConstructor) ConstructResponsePostitionMessage(resource string, position int64, responsechannel chan positionkeepermessage.PositionKeeperMessage) positionkeepermessage.PositionKeeperMessage {
	var positionKeeperMessage positionkeepermessage.PositionKeeperMessage
	positionKeeperMessage.Action = "response"
	positionKeeperMessage.Path = resource
	positionKeeperMessage.Position = position
	positionKeeperMessage.Channel = responsechannel
	return positionKeeperMessage
}

func (p *PositionKeeperMessageConstructor) ConstructAckPostitionMessage(resource string, responsechannel chan positionkeepermessage.PositionKeeperMessage) positionkeepermessage.PositionKeeperMessage {
	var positionKeeperMessage positionkeepermessage.PositionKeeperMessage
	positionKeeperMessage.Action = "ack"
	positionKeeperMessage.Path = resource
	positionKeeperMessage.Channel = responsechannel
	return positionKeeperMessage
}

func (p *PositionKeeperMessageConstructor) ExtractPosition(positionKeeperMessage positionkeepermessage.PositionKeeperMessage) int64 {
	return positionKeeperMessage.Position
}

func (p *PositionKeeperMessageConstructor) ExtractResource(positionKeeperMessage positionkeepermessage.PositionKeeperMessage) string {
	return positionKeeperMessage.Path
}

func (p *PositionKeeperMessageConstructor) ExtractChannel(positionKeeperMessage positionkeepermessage.PositionKeeperMessage) chan positionkeepermessage.PositionKeeperMessage {
	return positionKeeperMessage.Channel
}

func (p *PositionKeeperMessageConstructor) ExtractAction(positionKeeperMessage positionkeepermessage.PositionKeeperMessage) string {
	return positionKeeperMessage.Action
}
