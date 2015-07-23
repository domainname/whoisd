package service

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/domainname/daemon"
	"github.com/domainname/whoisd/client"
	"github.com/domainname/whoisd/config"
	"github.com/domainname/whoisd/storage"
)

// Version of the Whois Daemon
// Date of current version release
const (
	Version = "0.1.0"
	Date    = "2014-08-02T12:01:33Z"
)

// Record - standard record (struct) for service package
type Record struct {
	Name   string
	Config *config.Record
	daemon.Daemon
}

// New - Create a new service record
func New(name, description string) (*Record, error) {
	daemonInstance, err := daemon.New(name, description)
	if err != nil {
		return nil, err
	}

	return &Record{name, config.New(), daemonInstance}, nil
}

// Run or manage the service
func (service *Record) Run() (string, error) {

	// if received any kind of command, do it
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		}
	}

	// Load configuration and get mapping
	mapp, err := service.Config.Load()
	if err != nil {
		return "Loading mapping file was unsuccessful", err
	}

	// Logs for what is host&port used
	serviceHostPort := fmt.Sprintf("%s:%d", service.Config.Host, service.Config.Port)
	log.Printf("%s started on %s\n", service.Name, serviceHostPort)
	log.Printf("Used storage %s on %s:%d\n",
		service.Config.Storage.StorageType,
		service.Config.Storage.Host,
		service.Config.Storage.Port,
	)

	// Set up listener for defined host and port
	listener, err := net.Listen("tcp", serviceHostPort)
	if err != nil {
		return "Possibly was a problem with the port binding", err
	}

	// set up channel to collect client queries
	channel := make(chan client.Record, service.Config.Connections)

	// set up current storage
	repository := storage.New(service.Config, mapp)

	// init workers
	for i := 0; i < service.Config.Workers; i++ {
		go client.ProcessClient(channel, repository)
	}

	// This block is for testing purpose only
	if service.Config.TestMode == true {
		// make pipe connections for testing
		// connIn will ready to write into by function ProcessClient
		connIn, connOut := net.Pipe()
		defer connIn.Close()
		defer connOut.Close()
		newClient := client.Record{Conn: connIn}

		// prepare query for ProcessClient
		newClient.Query = []byte(service.Config.TestQuery)

		// send it into channel
		channel <- newClient
		// just read answer from channel pipe
		buffer := make([]byte, 4096)
		numBytes, err := connOut.Read(buffer)
		log.Println("Read bytes:", numBytes)
		return string(buffer), err
	}

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	// set up channel on which to send accepted connections
	listen := make(chan net.Conn, service.Config.Connections)
	go acceptConnection(listener, listen)

	// loop work cycle with accept connections or interrupt
	// by system signal
	for {
		select {
		case conn := <-listen:
			newClient := client.Record{Conn: conn}
			go newClient.HandleClient(channel)
		case killSignal := <-interrupt:
			log.Println("Got signal:", killSignal)
			log.Println("Stoping listening on ", listener.Addr())
			listener.Close()
			if killSignal == os.Interrupt {
				return "Daemon was interruped by system signal", nil
			}
			return "Daemon was killed", nil
		}
	}

	// never happen, but need to complete code
	return "If you see that, you are lucky bastard", nil
}

// Accept a client connection and collect it in a channel
func acceptConnection(listener net.Listener, listen chan<- net.Conn) {
	defer func() {
		if recovery := recover(); recovery != nil {
			log.Println("Recovered in ListenConnection:", recovery)
		}
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		listen <- conn
	}
}
