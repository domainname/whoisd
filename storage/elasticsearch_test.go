package storage

import (
	"testing"
)

func TestElasticSearchEmpty(t *testing.T) {
	dummy := ElasticsearchRecord{"localhost", 9200, "whois", "domain"}
	var emptyResult map[string][]string
	var err error
	emptyResult, err = dummy.Search("name", "")
	if err == nil {
		t.Error("Expected error for empty query, got", err)
	}
	emptyResult, err = dummy.Search("name", "aaa")
	if len(emptyResult) != 0 {
		t.Error("Expected len of empty query", 0, ", got", len(emptyResult))
	}
	emptyResult, err = dummy.SearchRelated("customer", "", "")
	if err == nil {
		t.Error("Expected error for empty query, got", err)
	}
	emptyResult, err = dummy.SearchRelated("customer", "handle", "AA-BB")
	if len(emptyResult) != 0 {
		t.Error("Expected len of empty query", 0, ", got", len(emptyResult))
	}
	emptyResult, err = dummy.SearchMultiple("nameserver", "", "")
	if err == nil {
		t.Error("Expected error for empty query, got", err)
	}
	emptyResult, err = dummy.SearchMultiple("nameserver", "nsgroupId", "0")
	if len(emptyResult) != 0 {
		t.Error("Expected len of empty query", 0, ", got", len(emptyResult))
	}
}
