package reader_tailfile

import (
	"bufio"
	"fmt"
	"github.com/dcodix/grigori/lib/domain/message"
	"github.com/dcodix/grigori/lib/domain/positionkeeper"
	"github.com/dcodix/grigori/lib/domain/reader"
	"github.com/dcodix/grigori/lib/infrastructure/register_modules"
	"io"
	"log"
	"os"
	"time"
)

type ReaderTailFile struct {
	Resource                         string
	MaxLines                         int
	PositionKeeperComunicator        positionkeeper.PositionKeeperComunicator
	PositionKeeperMessageConstructor positionkeeper.PositionKeeperMessageConstructor
	position                         int64
}

// init registers with the program.
func init() {
	reader := new(ReaderTailFile)
	register_modules.RegisterReader("tail", reader)
}

func (r *ReaderTailFile) Config(config map[string]interface{}) {
	r.Resource = config["resource"].(string)
	if maxlines, ok := config["maxlines"]; ok {
		r.MaxLines = int(maxlines.(float64))
	} else {
		r.MaxLines = 50
	}
	log.Println("Monitoring: " + r.Resource)
}

func (r *ReaderTailFile) SetPPositionKeeperMessageConstructor(constructor positionkeeper.PositionKeeperMessageConstructor) {
	r.PositionKeeperMessageConstructor = constructor
}

func (r *ReaderTailFile) SetPositionKeeperComunicator(communicator positionkeeper.PositionKeeperComunicator) {
	r.PositionKeeperComunicator = communicator
}

func (r *ReaderTailFile) Clone() reader.Reader {
	clonedReader := new(ReaderTailFile)
	return clonedReader
}

func openFile(filePath string, quit chan bool) (tailFile *os.File, ok bool) {
BREAKFOR:
	for {
		select {
		case <-quit:
			break BREAKFOR
		default:

		}
		tailFile, err := os.Open(filePath)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Duration(5) * time.Second)
		} else {
			return tailFile, true
		}
	}
	return tailFile, false
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (r *ReaderTailFile) positionKeeperIsConfigured() bool {
	positionKeeperIsConfigured := false
	if r.PositionKeeperComunicator != nil && r.PositionKeeperMessageConstructor != nil {
		positionKeeperIsConfigured = true
	}
	return positionKeeperIsConfigured
}

func (r *ReaderTailFile) getCurrentPosition(positionKeeperIsConfigured bool) (currentPosition int64) {
	currentPosition = r.position
	if r.position == 0 && positionKeeperIsConfigured {
		//Load Position
		loadMessage := r.PositionKeeperMessageConstructor.ConstructLoadPostitionMessage(r.Resource, r.PositionKeeperComunicator.GetResponseChannel())
		r.PositionKeeperComunicator.SendMessage(loadMessage, r.PositionKeeperComunicator.GetMainChannel())

		positionResponseMessage, notFound := r.PositionKeeperComunicator.GetMessage()
		if !notFound {
			newPosition := r.PositionKeeperMessageConstructor.ExtractPosition(positionResponseMessage)
			r.position = newPosition
			currentPosition = newPosition
		}
	}
	return currentPosition
}

func (r *ReaderTailFile) savePosition(position int64) {
	saveMessage := r.PositionKeeperMessageConstructor.ConstructSavePostitionMessage(r.Resource, position, r.PositionKeeperComunicator.GetResponseChannel())
	r.PositionKeeperComunicator.SendMessage(saveMessage, r.PositionKeeperComunicator.GetMainChannel())
}

func fileIsTruncated(input io.ReadSeeker, position int64) (trunkedFile bool, err error) {
	trunkedFile = false
	scanLines := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanLines(data, atEOF)
		position += int64(advance)
		return
	}

	if _, err := input.Seek(position-1, 0); err != nil {
		log.Println("Error trying to seek a position in file.")
	}
	scanner := bufio.NewScanner(input)
	scanner.Split(scanLines)

	if !scanner.Scan() && scanner.Err() == nil {
		trunkedFile = true
	}
	return trunkedFile, err
}

func (r *ReaderTailFile) readLines(input io.ReadSeeker, messages chan message.Message) (reachedMaxLines bool, err error) {
	//Create ticker to save at intervals
	quitSaveTickerGouroutine := make(chan bool)
	defer close(quitSaveTickerGouroutine)
	ticker := time.NewTicker(time.Duration(2) * time.Second)
	defer ticker.Stop()

	positionKeeperIsConfigured := r.positionKeeperIsConfigured()
	position := r.getCurrentPosition(positionKeeperIsConfigured)

	// Seek to the desidered position
	if _, err := input.Seek(position, 0); err != nil {
		check(err)
	}
	scanner := bufio.NewScanner(input)

	// Define tokens to be used for line reading
	scanLines := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanLines(data, atEOF)
		position += int64(advance)
		return
	}
	scanner.Split(scanLines)

	//Save position at intervals given by ticker
	if positionKeeperIsConfigured {
		go func(r *ReaderTailFile) {
			select {
			case <-ticker.C:
				r.savePosition(position)
				r.position = position
			case <-quitSaveTickerGouroutine:
				break
			}
		}(r)
	}

	// Scan lines unil EOF or until MaxLines
	numLines := 0
	var message message.Message
	for scanner.Scan() {
		message.Message = scanner.Text()
		messages <- message

		r.position = position
		numLines = numLines + 1
		if numLines >= r.MaxLines {
			reachedMaxLines = true
			break
		}
	}
	if isTruncated, _ := fileIsTruncated(input, position); isTruncated == true {
		position = 0
	}
	r.position = position
	if positionKeeperIsConfigured {
		r.savePosition(position)
		//Close timed save position goroutine
		quitSaveTickerGouroutine <- true
	}

	return reachedMaxLines, scanner.Err()
}

func (r *ReaderTailFile) Read(messages chan message.Message, quit chan bool) {
	if r.Resource == "" {
		fmt.Println("Resource not specified.")
		return
	}
	tailFile, ok := openFile(r.Resource, quit)
	if ok {
		defer tailFile.Close()
	} else {
		return
	}

	for {
		select {
		case <-quit:
			return
		default:

		}
		reachedMaxLines, _ := r.readLines(tailFile, messages)
		//time.Sleep(time.Duration(5000000) * time.Nanosecond)
		if reachedMaxLines {
			time.Sleep(time.Duration(5000) * time.Microsecond)
		} else {
			time.Sleep(time.Duration(500000) * time.Microsecond)
		}

	}
}

func (r *ReaderTailFile) GetResource() string {
	return r.Resource
}
