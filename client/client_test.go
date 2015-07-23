package client

import (
	"net"
	"strings"
	"testing"

	"github.com/domainname/whoisd/config"
	"github.com/domainname/whoisd/storage"
)

func TestClientHandling(t *testing.T) {
	conf := config.New()
	conf.ConfigPath = "../test/testconfig.conf"
	conf.MappingPath = "../test/testmapping.json"
	mapp, err := conf.Load()
	if err != nil {
		t.Error("Expected config loading without error, got", err.Error())
	}
	channel := make(chan Record, conf.Connections)
	repository := storage.New(conf, mapp)
	go ProcessClient(channel, repository)

	// make pipe connections for testing
	// connIn will ready to write into by function ProcessClient
	connIn, connOut := net.Pipe()
	defer connIn.Close()
	defer connOut.Close()
	newClient := Record{Conn: connIn}

	// prepare query for ProcessClient
	newClient.Query = []byte("google.com")

	// send it into channel
	channel <- newClient

	// just read answer from channel pipe
	buffer := make([]byte, 256)
	numBytes, err := connOut.Read(buffer)
	if err != nil {
		t.Error("Network communication error", err.Error())
	}
	if numBytes == 0 || len(buffer) == 0 {
		t.Error("Expexted some data read, got", string(buffer))
	}
	partAnswer := "Updated Date: 2014-05-19T04:00:17Z"
	if !strings.Contains(string(buffer), partAnswer) {
		t.Error("Expexted that contains", partAnswer, ", got", string(buffer))
	}
}
