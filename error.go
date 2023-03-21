package deepl

import "errors"

var (
	ErrorNoTranslateTextFound = errors.New("no Translate Text Found")
	ErrorInvalidTargetLang    = errors.New("invalid Target Lang")
	ErrorTooManyRequests      = errors.New("too Many Requests")
)
