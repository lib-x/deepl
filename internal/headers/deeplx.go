package headers

import (
	"iter"
)

var (
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

func DeeplXHeaderIter() iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		for headerKey, HeaderValue := range deeplXHeaders {
			if !yield(headerKey, HeaderValue) {
				return
			}
		}

	}
}
