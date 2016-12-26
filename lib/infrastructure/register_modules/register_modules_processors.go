package register_modules

import (
	"github.com/dcodix/grigori/lib/domain/processor"
	"log"
)

var Processors = make(map[string]processor.Processor)

// Register is called to register a processor for use by the program.
func RegisterProcessor(processorType string, processor processor.Processor) {
	if _, exists := Processors[processorType]; exists {
		log.Fatalln(processorType, "processor already registered")
	}
	Processors[processorType] = processor
}
