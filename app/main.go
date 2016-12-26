package main

import (
	"github.com/dcodix/grigori/lib"
	"github.com/dcodix/grigori/lib/infrastructure/config/read_config_file"
	"github.com/droundy/goopt"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
)

func main() {
	useProfile := false

	if useProfile {
		f, err := os.Create("proffilename.prof")
		if err != nil {
			panic(err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	//Read command line arguments
	goopt.Description = func() string {
		return "Generic log shipper"
	}
	goopt.Version = "0.1"
	goopt.Summary = "This program will read data from some source, process it and write it to a destination. As an example, it can read logs, transform them to jsn and wtite them to redis."
	var config_file = goopt.String([]string{"-c", "--config"}, "/etc/grigori.cfg", "Config file path.")
	goopt.Parse(nil)

	//Read config file
	configReadedFile := read_config_file.ReadConfig(*config_file)

	//Grigori full configuration
	con := read_config_file.GetConfig(configReadedFile)
	con.ConfigReadedFile = configReadedFile

	//Listen for signals to stop app
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	//Run grigori
	grigori.Run(con, sigs)

}