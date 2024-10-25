package headers

import (
	"iter"
	"strings"
)

var (
	deeplProHeaders = map[string]string{
		"Referer":         deepLReferServer,
		"Origin":          deepLReferServer,
		"Content-Type":    "application/json",
		"Accept":          "*/*",
		"Accept-Language": "en-US,en;q=0.9",
		"Accept-Encoding": "gzip, deflate, br",
		"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
		"Connection":      "keep-alive",
		"Cookie":          "",
	}
)

func DeeplProHeaderIter(dlSession string) iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		for headerKey, HeaderValue := range deeplProHeaders {
			if strings.EqualFold(headerKey, "Cookie") {
				HeaderValue = "dl_session=" + dlSession
			}
			if !yield(headerKey, HeaderValue) {
				return
			}
		}
	}
}
