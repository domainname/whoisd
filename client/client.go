package client

import (
	"bytes"
	"log"
	"net"

	"github.com/takama/whoisd/storage"
)

const (
	queryBufferSize = 256
)

// Record - standard record (struct) for client package
type Record struct {
	Conn  net.Conn
	Query []byte
}

// HandleClient - Sends a client data into the channel
func (client *Record) HandleClient(channel chan<- Record) {
	defer func() {
		if recovery := recover(); recovery != nil {
			log.Println("Recovered in HandleClient:", recovery)
			channel <- *client
		}
	}()
	buffer := make([]byte, queryBufferSize)
	numBytes, err := client.Conn.Read(buffer)
	if numBytes == 0 || err != nil {
		return
	}
	client.Query = bytes.ToLower(bytes.Trim(buffer, "\u0000\u000a\u000d"))
	channel <- *client
}

// ProcessClient - Asynchronous a client handling
func ProcessClient(channel <-chan Record, repository *storage.Record) {
	message := Record{}
	defer func() {
		if recovery := recover(); recovery != nil {
			log.Println("Recovered in ProcessClient:", recovery)
			if message.Conn != nil {
				message.Conn.Close()
			}
		}
	}()
	for {
		message = <-channel
		query := string(message.Query)
		data, ok := repository.Search(query)
		message.Conn.Write([]byte(data))
		log.Println(message.Conn.RemoteAddr().String(), query, ok)
		message.Conn.Close()
	}
}
