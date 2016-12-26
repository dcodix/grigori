package read_config_file

import (
	"encoding/json"
	"fmt"
	"github.com/dcodix/grigori/lib/domain/config"
	"github.com/dcodix/grigori/lib/domain/positionkeepermessage"
	"github.com/dcodix/grigori/lib/domain/processor"
	"github.com/dcodix/grigori/lib/domain/writer"
	"github.com/dcodix/grigori/lib/infrastructure/positionkeeper/positionkeepercommunicator"
	"github.com/dcodix/grigori/lib/infrastructure/positionkeeper/positionkeepermessageconstructor"
	"github.com/dcodix/grigori/lib/infrastructure/positionkeeper/positionkeeperrepository"
	"github.com/dcodix/grigori/lib/infrastructure/register_modules"
	"io/ioutil"
	"os"
)

func ReadConfig(filepath string) (configReded config.ConfigReadedFile) {
	file, e := ioutil.ReadFile(filepath)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	err := json.Unmarshal(file, &configReded)
	if err != nil {
		fmt.Println(err)
	}
	return configReded
}

func getWriters(readedFileConfiguration config.ConfigReadedFile) []writer.Writer {
	var writers []writer.Writer
	for i := 1; i <= int(readedFileConfiguration.ConfigReaded.Writer["n_writers"].(float64)); i++ {
		writer := register_modules.Writers[readedFileConfiguration.ConfigReaded.Writer["type"].(string)].Clone()
		writer.Config(readedFileConfiguration.ConfigReaded.Writer)
		writers = append(writers, writer)
	}
	return writers
}

func getCommunicator(readedFileConfiguration config.ConfigReadedFile) *positionkeepercommunicator.PositionKeeperComunicator {
	communicator := new(positionkeepercommunicator.PositionKeeperComunicator)
	communicatorChannel := make(chan positionkeepermessage.PositionKeeperMessage, 10)
	communicator.SetMainChannel(communicatorChannel)
	communicator.SetResponseChannel(communicatorChannel)
	return communicator
}

func getResources(readedFileConfiguration config.ConfigReadedFile, communicator *positionkeepercommunicator.PositionKeeperComunicator, constructor *positionkeepermessageconstructor.PositionKeeperMessageConstructor) []config.Resource {
	var resources []config.Resource
	for _, resourceReaded := range readedFileConfiguration.ConfigReaded.Resources {
		var resource config.Resource
		//Reader
		resource.Reader = register_modules.Readers[resourceReaded["reader"].(string)].Clone()
		resource.Reader.Config(resourceReaded)
		resource.Reader.SetPositionKeeperComunicator(communicator.CreateResponseCommunicator())
		resource.Reader.SetPPositionKeeperMessageConstructor(constructor)

		//Processors
		var processors []processor.Processor
		for i := 1; i <= int(resourceReaded["n_processors"].(float64)); i++ {
			processor := register_modules.Processors[resourceReaded["processor"].(string)].Clone()
			processor.Config(resourceReaded)
			processors = append(processors, processor)
		}
		resource.Processors = processors

		//Resources
		resources = append(resources, resource)

	}
	return resources
}

func GetConfig(readedFileConfiguration config.ConfigReadedFile) (configToGrigori *config.Config) {
	//Writers
	writers := getWriters(readedFileConfiguration)

	//PositionKeeper
	positionkeeper := register_modules.PositionKeepers[readedFileConfiguration.PositionKeeper["type"]]
	positionkeeper.Config(readedFileConfiguration.PositionKeeper)

	//Communicator
	communicator := getCommunicator(readedFileConfiguration)

	//Constructor
	constructor := new(positionkeepermessageconstructor.PositionKeeperMessageConstructor)

	//Repository
	repository := new(positionkeeperrepository.PositionKeeperRepository)
	repository.SetCommunicator(communicator)
	repository.SetConstructor(constructor)
	repository.SetPositionKeeper(positionkeeper)

	//Resources
	resources := getResources(readedFileConfiguration, communicator, constructor)

	return &config.Config{resources, writers, repository, constructor, communicator, readedFileConfiguration}
}
