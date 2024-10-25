package postprocess

import "strings"

func CountAlphaI(translateText string) int64 {
	return int64(strings.Count(translateText, "i"))
}
