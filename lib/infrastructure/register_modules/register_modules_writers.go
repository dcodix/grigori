package register_modules

import (
	"github.com/dcodix/grigori/lib/domain/writer"
	"log"
)

var Writers = make(map[string]writer.Writer)

// Register is called to register a writer for use by the program.
func RegisterWriter(writerType string, writer writer.Writer) {
	if _, exists := Writers[writerType]; exists {
		log.Fatalln(writerType, "writer already registered")
	}
	Writers[writerType] = writer
}
