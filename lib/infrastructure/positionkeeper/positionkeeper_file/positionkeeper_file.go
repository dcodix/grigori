package positionkeeper_file

import (
	"encoding/json"
	"fmt"
	"github.com/dcodix/grigori/lib/domain/positionkeepermessage"
	"github.com/dcodix/grigori/lib/infrastructure/register_modules"
	"io/ioutil"
	"os"
	"log"
)

type PositionKeeperFile struct {
	path      string
	positions positionkeepermessage.JsonPositions
}

// init registers with the program.
func init() {
	positionKeeper := new(PositionKeeperFile)
	register_modules.RegisterPositionKeeper("file", positionKeeper)
}

func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (p *PositionKeeperFile) prepare() {
	if !exists(p.path) {
		var file, err = os.Create(p.path)
		check(err)
		file.Close()
	}
}

func (p *PositionKeeperFile) Config(config map[string]string) {
	p.path = config["path"]
	p.positions.Object.Positions = map[string]positionkeepermessage.Position{} //Initialize NOT NIL
	p.prepare()
	log.Println("Using position keeper FILE: " + p.path)
}

func (p *PositionKeeperFile) writePositionFile() {
	positionJsonFormat, _ := json.Marshal(p.positions)
	err := ioutil.WriteFile(p.path, positionJsonFormat, 0644)
	if err != nil {
		fmt.Printf("WriteFileJson ERROR: %+v", err)
	}
}

func (p *PositionKeeperFile) Save(resource string, position int64, saveToBackend bool) {
	var newPosition positionkeepermessage.Position
	newPosition.Position = position
	p.positions.Object.Positions[resource] = newPosition
	if saveToBackend {
		p.writePositionFile()
	}
}

func (p *PositionKeeperFile) readPositionFile() {
	filepath := p.path
	file, e := ioutil.ReadFile(filepath)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	var fileReadedPositions positionkeepermessage.JsonPositions
	json.Unmarshal(file, &fileReadedPositions)
	if len(fileReadedPositions.Object.Positions) > 0 {
		p.positions = fileReadedPositions
	}
}

func (p *PositionKeeperFile) Load(resource string) (resourcePosition int64, notFound bool) {
	notFound = false

	p.readPositionFile()

	if position, ok := p.positions.Object.Positions[resource]; ok {
		resourcePosition = position.Position
	} else {
		notFound = true
	}

	return resourcePosition, notFound
}
