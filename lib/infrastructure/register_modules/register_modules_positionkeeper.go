package register_modules

import (
	"github.com/dcodix/grigori/lib/domain/positionkeeper"
	"log"
)

var PositionKeepers = make(map[string]positionkeeper.PositionKeeper)

// Register is called to register a reader for use by the program.
func RegisterPositionKeeper(positionKeeperType string, positionKeeper positionkeeper.PositionKeeper) {
	if _, exists := PositionKeepers[positionKeeperType]; exists {
		log.Fatalln(positionKeeperType, "Position Keeper already registered")
	}

	PositionKeepers[positionKeeperType] = positionKeeper

}
