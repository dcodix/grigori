package writer_stdout

import (
	"github.com/dcodix/grigori/lib/domain/message"
	"github.com/dcodix/grigori/lib/domain/writer"
	"github.com/dcodix/grigori/lib/infrastructure/register_modules"
	"log"
)

type WriterStdout struct{}

// init registers with the program.
func init() {
	writer := new(WriterStdout)
	register_modules.RegisterWriter("stdout", writer)
}

func (w *WriterStdout) Config(config map[string]interface{}) {
	log.Println("Writing to: STDOUT")
}

func (w *WriterStdout) Clone() writer.Writer {
	cloneWriter := new(WriterStdout)
	return cloneWriter
}

func (w *WriterStdout) Write(messages chan message.Message) {
	for message := range messages {
		log.Println(message)
	}
}
