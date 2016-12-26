package monitoring

import (
	"github.com/dcodix/grigori/lib/domain/config"
	"github.com/dcodix/grigori/lib/domain/message"
	"log"
	"net"
	"net/http"
	"strconv"
)

type MonitoringElements struct {
	Channels map[string]chan message.Message
	Config   *config.Config
}

var monitoringStruct MonitoringElements

func Monitor(monitoringElements MonitoringElements, port int, quit chan bool) {
	monitoringStruct = monitoringElements
	router := NewRouter()

	//Default to port
	if port == 0 {
		port = 8080
	}

	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("Can't listen: %s", err)
	}

	go func() {
		_ = <-quit
		ln.Close()
	}()

	http.Serve(ln, router)
}
