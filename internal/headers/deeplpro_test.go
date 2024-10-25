package headers

import "testing"

func TestDeeplProHeader(t *testing.T) {
	headerIter := DeeplProHeaderIter("test")
	for s, s2 := range headerIter {
		t.Log(s, "=>", s2)
	}
}
