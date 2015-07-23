package storage

import (
	"testing"
)

func TestDummySearch(t *testing.T) {

	type testData struct {
		query    string
		name     string
		ask      string
		expected []string
	}

	var tests = []testData{
		{"google.com", "name", "ownerHandle", []string{"MMR-2383"}},
		{"google.com", "name", "updatedDate", []string{"2014-05-19 04:00:17"}},
		{"google.com", "name", "dnssec", []string{"unsigned"}},
		{"example.tld", "name", "techHandle", []string{"5372811-ERL"}},
		{"example.tld", "name", "domainStatus",
			[]string{
				"clientDeleteProhibited",
				"clientRenewProhibited",
				"clientTransferProhibited",
			},
		},
		{"example.tld", "name", "dnssec", []string{"signedDelegation"}},
	}

	dummy := DummyRecord{"localhost", 9200, "whois", "domain"}
	for _, data := range tests {
		result, err := dummy.Search(data.name, data.query)
		if err != nil {
			t.Error(err.Error())
		} else {
			if len(result) == 0 {
				t.Error("Empty set for", data.query)
			}
			for index, item := range result[data.ask] {
				if item != data.expected[index] {
					t.Error("Expected", data.expected, ", got", item)
				}
			}
		}
	}
}

func TestDummySearchRelated(t *testing.T) {

	type testData struct {
		query     string
		name      string
		relatedTo string
		ask       string
		expected  []string
	}

	var tests = []testData{
		{"MMR-2383", "handle", "customer", "address.street", []string{"1600 Amphitheatre Parkway"}},
		{"MMR-2383", "handle", "customer", "email", []string{"dns-admin@google.com"}},
		{"MMR-2383", "handle", "customer", "name.lastName", []string{"Admin"}},
		{"MMA-2211", "handle", "customer", "address.street", []string{"2400 E. Bayshore Pkwy"}},
	}

	dummy := DummyRecord{"localhost", 9200, "whois", "domain"}
	for _, data := range tests {
		result, err := dummy.SearchRelated(data.relatedTo, data.name, data.query)
		if err != nil {
			t.Error(err.Error())
		} else {
			if len(result) == 0 {
				t.Error("Empty set for", data.query)
			}
			for index, item := range result[data.ask] {
				if item != data.expected[index] {
					t.Error("Expected", data.expected[index], ", got", item)
				}
			}
		}
	}
}

func TestDummySearchMultiple(t *testing.T) {

	type testData struct {
		query     string
		name      string
		relatedTo string
		ask       string
		expected  []string
	}

	var tests = []testData{
		{"1", "nsgroupId", "nameserver", "name",
			[]string{
				"NS01.EXAMPLE-REGISTRAR.TLD",
				"NS02.EXAMPLE-REGISTRAR.TLD",
			},
		},
		{"2", "nsgroupId", "nameserver", "name",
			[]string{
				"ns1.google.com",
				"ns2.google.com",
				"ns3.google.com",
				"ns4.google.com",
			},
		},
	}

	dummy := DummyRecord{"localhost", 9200, "whois", "domain"}
	for _, data := range tests {
		result, err := dummy.SearchMultiple(data.relatedTo, data.name, data.query)
		if err != nil {
			t.Error(err.Error())
		} else {
			if len(result) == 0 {
				t.Error("Empty set for", data.query)
			}
			if len(result[data.ask]) != len(data.expected) {
				t.Error("No multiple records, expected", len(data.expected), ", got", len(result[data.ask]))
			}
			for index, item := range result[data.ask] {
				if item != data.expected[index] {
					t.Error("Expected", data.expected[index], ", got", item)
				}
			}
		}
	}
}

func TestDummySearchEmpty(t *testing.T) {
	dummy := DummyRecord{"localhost", 9200, "whois", "domain"}
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
	emptyResult, err = dummy.SearchMultiple("nameserver", "nsgroupId", "7")
	if len(emptyResult) != 0 {
		t.Error("Expected len of empty query", 0, ", got", len(emptyResult))
	}
}
