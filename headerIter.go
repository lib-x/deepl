package deepl

import (
	"iter"
	"strings"
)

const (
	deepLReferServer = "https://www.deepl.com/"
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

	deeplXHeaders = map[string]string{
		"Referer":          deepLReferServer,
		"Content-Type":     "application/json",
		"Accept":           "*/*",
		"x-app-os-name":    "iOS",
		"x-app-os-version": "16.3.0",
		"Accept-Language":  "en-US,en;q=0.9",
		"Accept-Encoding":  "gzip, deflate, br",
		"x-app-device":     "iPhone13,2",
		"User-Agent":       "DeepL-iOS/2.9.1 iOS 16.3.0 (iPhone13,2)",
		"x-app-build":      "510265",
		"x-app-version":    "2.9.1",
		"Connection":       "keep-alive",
	}
)

func deeplProHeaderIter(dlSession string) iter.Seq2[string, string] {
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

func deeplXHeaderIter() iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		for headerKey, HeaderValue := range deeplXHeaders {
			if !yield(headerKey, HeaderValue) {
				return
			}
		}

	}
}
