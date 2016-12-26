package config

import (
	"github.com/dcodix/grigori/lib/domain/positionkeeper"
	"github.com/dcodix/grigori/lib/domain/processor"
	"github.com/dcodix/grigori/lib/domain/reader"
	"github.com/dcodix/grigori/lib/domain/writer"
)

type Resource struct {
	Reader     reader.Reader
	Processors []processor.Processor
}

type Config struct {
	Resources                        []Resource
	Writers                          []writer.Writer
	PositionKeeperRepository         positionkeeper.PositionKeeperRepository
	PositionKeeperMessageConstructor positionkeeper.PositionKeeperMessageConstructor
	PositionKeeperComunicator        positionkeeper.PositionKeeperComunicator
	ConfigReadedFile                 ConfigReadedFile
}

type ConfigReadedFile struct {
	ConfigReaded `json:"config"`
	Monitoring   `json:"monitoring"`
	Limits       `json:"limits"`
}

type ConfigReaded struct {
	Writer         map[string]interface{}   `json:"writer"`
	Resources      []map[string]interface{} `json:"resources"`
	PositionKeeper map[string]string        `json:"position_keeper"`
}

type Monitoring struct {
	Enabled bool `json:"enabled"`
	Port    int  `json:"port"`
}

type Limits struct {
	GoMaxProcs  int `json:"gomaxprocs"`
	QueueLength int `json:"queue_length"`
}
