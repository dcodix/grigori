package processor_plain_to_json

import (
	"encoding/json"
	"github.com/dcodix/grigori/lib/domain/message"
	"github.com/dcodix/grigori/lib/domain/processor"
	"github.com/dcodix/grigori/lib/infrastructure/register_modules"
	"os"
)

type JsonMessage struct {
	Source_host string   `json:"source_host,interface{}"`
	File        string   `json:"file,interface{}"`
	Message     string   `json:"message,intefcae{}"`
	Tags        []string `json:"tags,[]string"`
	Type        string   `json:"type,inerface{}"`
}

type ProcessorPlainToJson struct {
	Tags              []string
	File              string
	Type              string
	hostname          string
	messageProcessed  message.Message
	jsonMessageObject JsonMessage
}

// init registers with the program.
func init() {
	processor := new(ProcessorPlainToJson)
	register_modules.RegisterProcessor("plaintojson", processor)
}

func (p *ProcessorPlainToJson) Process(messagesIn chan message.Message, messagesOut chan message.Message) {
	for messageIn := range messagesIn {
		p.jsonMessageObject.Source_host = p.hostname
		p.jsonMessageObject.File = p.File
		p.jsonMessageObject.Message = messageIn.Message
		p.jsonMessageObject.Tags = p.Tags
		p.jsonMessageObject.Type = p.Type
		JsonMessageJsonFormat, _ := json.Marshal(p.jsonMessageObject)
		p.messageProcessed.Message = string(JsonMessageJsonFormat)

		messagesOut <- p.messageProcessed
	}
}

func getTags(config map[string]interface{}) []string {
	var configuredTags []string
	if tags, ok := config["tags"]; ok {
		for _, tag := range tags.([]interface{}) {
			configuredTags = append(configuredTags, tag.(string))
		}
	}
	return configuredTags
}

func getFile(config map[string]interface{}) string {
	var configuredFile string
	if resource, ok := config["resource"]; ok {
		configuredFile = resource.(string)
	}
	return configuredFile
}

func getType(config map[string]interface{}) string {
	var readedType string
	if log_type, ok := config["type"]; ok {
		readedType = log_type.(string)
	}
	return readedType
}

func (p *ProcessorPlainToJson) Config(config map[string]interface{}) {
	p.Tags = getTags(config)
	p.File = getFile(config)
	p.Type = getType(config)
	p.hostname, _ = os.Hostname()
}

func (p *ProcessorPlainToJson) Clone() processor.Processor {
	clonedProcessor := new(ProcessorPlainToJson)
	return clonedProcessor
}
