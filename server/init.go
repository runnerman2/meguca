// Package server handles client requests, both for HTML page rendering and
// websocket connections.
package server

import (
	"fmt"
	"github.com/sevlyar/go-daemon"
	"log"
	"os"
	"syscall"
	"time"
)

// Start parses command line arguments and initializes the server.
func Start() {
	loadConfig()

	// Parse command line arguments
	if len(os.Args) < 2 {
		printUsage()
	}
	arg := os.Args[1]
	switch arg {
	case "init":
		os.Exit(0)
	case "debug":
		config.Hard.Debug = true
	case "start":
	case "stop":
		killDaemon()
		os.Exit(0)
	case "restart":
		killDaemon()
	default:
		printUsage()
	}

	if !config.Hard.Debug {
		daemonise()
	} else {
		startServer()
	}
}

func printUsage() {
	fmt.Print(`usage: ./meguca [ start | stop | debug | help ]
	start   - start the meguca server
	stop    - stop a running daemonised meguca server
	restart - combination of stop + start
	init    - create intial files and folders and exit
	debug   - force debug mode
	help    - print this help text
`)
	os.Exit(1)
}

func startServer() {
	loadLanguagePacks()
	compileTemplates()
	loadDB()
	startWebServer()
}

// Configuration variables for handling daemons
var daemonContext = &daemon.Context{
	PidFileName: "./.pid",
	LogFileName: "./error.log",
}

// Spawn a detached process to work in the background
func daemonise() {
	child, err := daemonContext.Reborn()
	if err != nil {
		if err.Error() == "resource temporarily unavailable" {
			fmt.Println("Error: Server already running")
			os.Exit(1)
		}
		panic(err)
	}
	if child != nil {
		return
	}
	defer daemonContext.Release()
	log.Println("Server started ------------------------------------")
	go startServer()
	throw(daemon.ServeSignals())
	log.Println("Server terminated")
}

// Terminate the running meguca server daemon
func killDaemon() {
	proc, err := daemonContext.Search()
	if err != nil && (!os.IsNotExist(err) && err.Error() != "EOF") {
		panic(err)
	}
	if proc != nil {
		throw(proc.Signal(syscall.SIGTERM))

		// Assertain process has exited
		for {
			if err := proc.Signal(syscall.Signal(0)); err != nil {
				if err.Error() == "os: process already finished" {
					break
				}
				panic(err)
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}
