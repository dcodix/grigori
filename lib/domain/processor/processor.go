package processor

import (
	"github.com/dcodix/grigori/lib/domain/message"
)

type Processor interface {
	Process(messagesIn chan message.Message, messagesOut chan message.Message)
	Config(map[string]interface{})
	Clone() Processor
}
