package processor_passthrough

import (
	"github.com/dcodix/grigori/lib/domain/message"
	"github.com/dcodix/grigori/lib/domain/processor"
	"github.com/dcodix/grigori/lib/infrastructure/register_modules"
)

type ProcessorPassthrough struct {
}

// init registers the matcher with the program.
func init() {
	processor := new(ProcessorPassthrough)
	register_modules.RegisterProcessor("passthrough", processor)
}

func (p *ProcessorPassthrough) Process(messagesIn chan message.Message, messagesOut chan message.Message) {
	for message := range messagesIn {
		messagesOut <- message
	}
}

func (p *ProcessorPassthrough) Config(config map[string]interface{}) {

}

func (p *ProcessorPassthrough) Clone() processor.Processor {
	clonedProcessor := new(ProcessorPassthrough)
	return clonedProcessor

}
