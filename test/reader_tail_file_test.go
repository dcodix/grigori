package main

import (
	"github.com/dcodix/grigori/lib/domain/message"
	"github.com/dcodix/grigori/lib/infrastructure/reader/reader_tailfile"
	"os"
	"sync"
	"testing"
	"time"
)

func CreateTestFile(testFileName string, max int) {
	if !Exists(testFileName) {
		var file, err = os.Create(testFileName)
		check(err)
		file.Close()
	}
	f, err := os.OpenFile(testFileName, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	for i := 0; i < max; i++ {

		_, err := f.WriteString("THIS IS A TEST\n")
		if err != nil {
			panic(err)
		}
	}
}

func trunkFileAddingLine(testFileName string) {
	if !Exists(testFileName) {
		var file, err = os.Create(testFileName)
		check(err)
		file.Close()
	}
	f, err := os.OpenFile(testFileName, os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = f.WriteString("FILE_TRUNKED\n")
	if err != nil {
		panic(err)
	}

}

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func DeleteTestFiles(testFile string) {
	os.Remove(testFile)
}

func countMessages(messagesReaded chan message.Message, quit_counting chan bool) (n int) {
	n = 0
FORBREAK:
	for {
		select {
		case <-messagesReaded:
			n += 1
		case <-quit_counting:
			break FORBREAK
		}
	}
	return n
}

func getOneMessage(messagesReaded chan message.Message) (message message.Message) {
	for i := 0; i < 10; i++ {
		select {
		case msg := <-messagesReaded:
			message = msg
		default:

		}
		time.Sleep(200 * time.Microsecond)
	}
	return message
}

func TestReaderTailFile(test *testing.T) {
	testFileName := "/tmp/testfilename828j74hfhnd9nd39nd38yn4938.log"
	n_lines_file := 5
	var reader_wg sync.WaitGroup
	quit := make(chan bool)
	quit_counting_messages := make(chan bool)
	r := &reader_tailfile.ReaderTailFile{Resource: testFileName, MaxLines: 20}

	DeleteTestFiles(testFileName)
	CreateTestFile(testFileName, n_lines_file)

	//GOROUTINE that will listen to signals to properly kill application
	go func() {
		time.Sleep(2 * time.Second)
		quit_counting_messages <- true
	}()

	messagesReaded := make(chan message.Message, 1000)
	reader_wg.Add(1)
	go func() {
		defer reader_wg.Done()
		r.Read(messagesReaded, quit)
	}()

	go func() {
		reader_wg.Wait()
		close(messagesReaded)
	}()

	n := countMessages(messagesReaded, quit_counting_messages)
	if n != n_lines_file {
		test.Error("Number of original messages and readed messages differ.")
	}

	//Test truncate file
	trunkFileAddingLine(testFileName)
	time.Sleep(2 * time.Second)
	message := getOneMessage(messagesReaded)
	if message.Message != "FILE_TRUNKED" {
		test.Error("Not getting expected results after truncating the file.")
		test.Error(message.Message)
	}
	DeleteTestFiles(testFileName)
}
