package processor_plain_to_logstash

import (
	"encoding/json"
	"github.com/dcodix/grigori/lib/domain/message"
	"github.com/dcodix/grigori/lib/domain/processor"
	"github.com/dcodix/grigori/lib/infrastructure/register_modules"
	"os"
	"time"
)

type LogstashMessage struct {
	Version     int       `json:"@version,int"`
	Timestamp   time.Time `json:"@timestamp,time.Time"`
	Source_host string    `json:"source_host,interface{}"`
	File        string    `json:"file,interface{}"`
	Message     string    `json:"message,intefcae{}"`
	Tags        []string  `json:"tags,[]string"`
	Type        string    `json:"type,inerface{}"`
}

type ProcessorPlainToLogstash struct {
	Tags                  []string
	File                  string
	LogstashVersion       int
	Type                  string
	hostname              string
	messageProcessed      message.Message
	logstashMessageObject LogstashMessage
}

// init registers with the program.
func init() {
	processor := new(ProcessorPlainToLogstash)
	register_modules.RegisterProcessor("plaintologstash", processor)
}

func (p *ProcessorPlainToLogstash) Process(messagesIn chan message.Message, messagesOut chan message.Message) {
	for messageIn := range messagesIn {
		now := time.Now()
		p.logstashMessageObject.Version = p.LogstashVersion
		p.logstashMessageObject.Timestamp = now
		p.logstashMessageObject.Source_host = p.hostname
		p.logstashMessageObject.File = p.File
		p.logstashMessageObject.Message = messageIn.Message
		p.logstashMessageObject.Tags = p.Tags
		p.logstashMessageObject.Type = p.Type
		logstashMessageJsonFormat, _ := json.Marshal(p.logstashMessageObject)
		p.messageProcessed.Message = string(logstashMessageJsonFormat)

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

func getLogstashVersion(config map[string]interface{}) int {
	var readedLogstashVersion int
	if logstash_version, ok := config["version"]; ok {
		readedLogstashVersion = int(logstash_version.(float64))
	}
	return readedLogstashVersion
}

func getType(config map[string]interface{}) string {
	var readedType string
	if log_type, ok := config["type"]; ok {
		readedType = log_type.(string)
	}
	return readedType
}

func (p *ProcessorPlainToLogstash) Config(config map[string]interface{}) {
	p.Tags = getTags(config)
	p.File = getFile(config)
	p.LogstashVersion = getLogstashVersion(config)
	p.Type = getType(config)
	p.hostname, _ = os.Hostname()
}

func (p *ProcessorPlainToLogstash) Clone() processor.Processor {
	clonedProcessor := new(ProcessorPlainToLogstash)
	return clonedProcessor
}
