package writer

import (
	"github.com/dcodix/grigori/lib/domain/message"
)

type Writer interface {
	Config(config map[string]interface{})
	Clone() (writer Writer)
	Write(Messages chan message.Message)
}
