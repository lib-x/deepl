package deepl

import "testing"

func TestDeeplProHeader(t *testing.T) {
	headerIter := deeplProHeaderIter("test")
	for s, s2 := range headerIter {
		t.Log(s, "=>", s2)
	}
}
