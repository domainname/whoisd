package service

import (
	"strings"
	"testing"
)

func TestService(t *testing.T) {
	daemonName, daemonDescription := "whoisd", "Whois Daemon"
	daemon, err := New(daemonName, daemonDescription)
	if err != nil {
		t.Error("Expected service create without error, got", err.Error())
	}
	if daemon.Name != daemonName {
		t.Error("Expected service name must be ", daemonName, ", got", daemon.Name)
	}
	if daemon.Config.Host != "0.0.0.0" {
		t.Error("Expected server host is 0.0.0.0, got", daemon.Config.Host)
	}
	daemon.Config.ConfigPath = "../test/testconfig.conf"
	daemon.Config.MappingPath = "../test/testmapping.json"
	answer, err := daemon.Run()
	if err != nil {
		t.Error("Expected service run without error, got", err.Error())
	}
	if daemon.Config.Host != "localhost" {
		t.Error("Expected server host is localhost, got", daemon.Config.Host)
	}
	partAnswer := "Updated Date: 2014-05-19T04:00:17Z"
	if !strings.Contains(string(answer), partAnswer) {
		t.Error("Expexted that contains", partAnswer, ", got", string(answer))
	}
}
