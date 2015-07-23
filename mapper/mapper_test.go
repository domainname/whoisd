package mapper

import (
	"testing"
)

func TestMapper(t *testing.T) {

	mapp := Record{}
	if mapp.Count() != 0 {
		t.Error("Expected empty mapper record, got", mapp.Count())
	}
}
