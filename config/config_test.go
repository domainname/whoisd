package config

import (
	"testing"
)

func TestConfig(t *testing.T) {

	conf := New()
	conf.ConfigPath = ""
	conf.MappingPath = ""
	mapp, err := conf.Load()
	if err == nil {
		t.Error("Expected error of loading mapping file, got nothing")
	}
	if conf.Connections != 1000 {
		t.Error("Expected 100 active connections, got", conf.Connections)
	}
	if conf.Workers != 1000 {
		t.Error("Expected 100 workers, got", conf.Workers)
	}
	if mapp != nil {
		t.Error("Expected nil mapper record, got not nil mapper record")
	}
	conf.ConfigPath = "../test/testconfig.conf"
	conf.MappingPath = "../test/testmapping.json"
	mapp, err = conf.Load()
	if err != nil {
		t.Error("Expected config loading without error, got", err.Error())
	}
	if conf.Connections != 100 {
		t.Error("Expected 100 active connections, got", conf.Connections)
	}
	if conf.Workers != 100 {
		t.Error("Expected 100 workers, got", conf.Workers)
	}
	if len(mapp.Fields) == 0 {
		t.Error("Expected loading of mapper, got empty mapper")
	}
	key := "01"
	expected := "Domain Name: "
	if mapp.Fields[key].Key != expected {
		t.Error("Expected", expected, ", got", mapp.Fields[key].Key)
	}
	key = "02"
	expected = "name"
	if mapp.Fields[key].Related != expected {
		t.Error("Expected", expected, ", got", mapp.Fields[key].Related)
	}
	key = "05"
	expected = "{date}"
	if mapp.Fields[key].Format != expected {
		t.Error("Expected", expected, ", got", mapp.Fields[key].Format)
	}
}
