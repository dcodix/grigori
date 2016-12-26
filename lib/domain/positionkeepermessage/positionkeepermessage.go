package positionkeepermessage

type PositionKeeperMessage struct {
	Action   string
	Path     string
	Position int64
	Channel  chan PositionKeeperMessage
}

type JsonPositions struct {
	Object ObjectType
}

type ObjectType struct {
	Positions map[string]Position
}

type Position struct {
	Position int64
}
