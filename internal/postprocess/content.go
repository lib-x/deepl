package postprocess

import (
	"bytes"
)

var (
	methodPartNormal       = []byte("\"method\":\"")
	methodPartWithOneSpace = []byte("\"method\": \"")
	methodPartWithTwoSpace = []byte("\"method\" : \"")
)

func AdjustContent(id int64, jsonContent []byte) []byte {
	// add space if necessary
	if (id+5)%29 == 0 || (id+3)%13 == 0 {
		jsonContent = bytes.ReplaceAll(jsonContent, methodPartNormal, methodPartWithTwoSpace)
	} else {
		jsonContent = bytes.ReplaceAll(jsonContent, methodPartNormal, methodPartWithOneSpace)
	}
	return jsonContent
}
