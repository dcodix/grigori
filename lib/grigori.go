package grigori

import (
	"fmt"
	"github.com/dcodix/grigori/lib/domain/config"
	"github.com/dcodix/grigori/lib/domain/message"
	"github.com/dcodix/grigori/lib/domain/processor"
	"github.com/dcodix/grigori/lib/domain/reader"
	"github.com/dcodix/grigori/lib/domain/writer"
	"github.com/dcodix/grigori/lib/infrastructure/monitoring"
	"os"
	"runtime"
	"sync"
)

func runReader(reader reader.Reader, wg *sync.WaitGroup, writers_wg *sync.WaitGroup, messagesReaded chan message.Message, quit chan bool) {
	go func() {
		defer wg.Done()
		defer writers_wg.Done()
		reader.Read(messagesReaded, quit)
	}()
}

func runProcessors(p processor.Processor, messagesReaded chan message.Message, messages chan message.Message, wg *sync.WaitGroup, processors_wg *sync.WaitGroup) {
	wg.Add(1)
	processors_wg.Add(1)
	go func() {
		defer wg.Done()
		defer processors_wg.Done()
		p.Process(messagesReaded, messages)
	}()
}

func runMessagesReadedQueuesCloser(messagesReadedQueues []chan message.Message, writers_wg *sync.WaitGroup, quitPositionKeeper chan bool) {
	go func() {
		writers_wg.Wait()
		for _, messagesReadedQueue := range messagesReadedQueues {
			close(messagesReadedQueue)
		}
		quitPositionKeeper <- true
	}()
}

func runWriters(w writer.Writer, messages chan message.Message, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		w.Write(messages)
	}()
}

func GetGoMaxProcs(configGoMaxProgs int) (goMaxProcs int) {
	if configGoMaxProgs == 0 {
		goMaxProcs = 1
	} else {
		goMaxProcs = configGoMaxProgs
	}
	return goMaxProcs
}

func GetQueueLength(configQueueLength int) (queueLength int) {
	if configQueueLength == 0 {
		queueLength = 1000
	} else {
		queueLength = configQueueLength
	}
	return queueLength
}

func Run(c *config.Config, sigs chan os.Signal) {
	runtime.GOMAXPROCS(GetGoMaxProcs(c.ConfigReadedFile.Limits.GoMaxProcs))

	queueLength := GetQueueLength(c.ConfigReadedFile.Limits.QueueLength)

	channelLib := make(map[string]chan message.Message)
	var monitoringElements monitoring.MonitoringElements
	quitMonitoring := make(chan bool)
	defer close(quitMonitoring)

	var wg sync.WaitGroup
	var writers_wg sync.WaitGroup
	var processors_wg sync.WaitGroup

	quitPositionKeeper := make(chan bool)

	wg.Add(1)
	go func() {
		defer wg.Done()
		c.PositionKeeperRepository.Run(quitPositionKeeper)
	}()

	messages := make(chan message.Message, queueLength)
	channelLib["writer"] = messages
	var quitReaders []chan bool
	var messagesReadedQueues []chan message.Message
	for _, r := range c.Resources {

		quit := make(chan bool)
		defer close(quit)
		quitReaders = append(quitReaders, quit)

		messagesReaded := make(chan message.Message, queueLength)
		messagesReadedQueues = append(messagesReadedQueues, messagesReaded)
		if c.ConfigReadedFile.Monitoring.Enabled {
			channelLib[r.Reader.GetResource()] = messagesReaded
		}
		wg.Add(1)
		writers_wg.Add(1)
		runReader(r.Reader, &wg, &writers_wg, messagesReaded, quit)
		for _, p := range r.Processors {
			runProcessors(p, messagesReaded, messages, &wg, &processors_wg)
		}

	}
	runMessagesReadedQueuesCloser(messagesReadedQueues, &writers_wg, quitPositionKeeper)

	go func() {
		processors_wg.Wait()
		close(messages)
	}()

	if c.ConfigReadedFile.Monitoring.Enabled {
		//Launch monitoring interface
		monitoringElements.Channels = channelLib
		monitoringElements.Config = c
		go func() {
			monitoring.Monitor(monitoringElements, c.ConfigReadedFile.Monitoring.Port, quitMonitoring)
		}()
	}
	for _, w := range c.Writers {
		runWriters(w, messages, &wg)
	}

	//GOROUTINE that will listen to signals to properly kill application
	go func() {
		sig := <-sigs
		fmt.Println(sig)
		fmt.Println("Shutting down Grigori.")
		for _, quitReader := range quitReaders {
			quitReader <- true
		}
		quitPositionKeeper <- true
		close(quitPositionKeeper)
	}()

	wg.Wait()
	if c.ConfigReadedFile.Monitoring.Enabled {
		quitMonitoring <- true
	}

}
