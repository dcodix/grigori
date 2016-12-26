package register_modules

import (
	"github.com/dcodix/grigori/lib/domain/reader"
	"log"
)

var Readers = make(map[string]reader.Reader)

// Register is called to register a reader for use by the program.
func RegisterReader(readerType string, reader reader.Reader) {
	if _, exists := Readers[readerType]; exists {
		log.Fatalln(readerType, "Reader already registered")
	}
	Readers[readerType] = reader
}
