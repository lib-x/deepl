package deepl

import (
	"github.com/lib-x/deepl/internal/deeplx"
	"github.com/lib-x/deepl/internal/option"
	"github.com/lib-x/deepl/internal/rpc"
)

// Translate request data. support lang code using deepL api
// DE: German
// EN: English
// ES: Spanish
// FR: French
// IT: Italian
// JA: Japanese
// NL: Dutch
// PL: Polish
// PT: Portuguese
// RU: Russian
// ZH: Chinese
// BG: Bulgarian
// CS: Czech
// DA: Danish
// EL: Greek
// ET: Estonian
// FI: Finnish
// HU: Hungarian
// LT: Lithuanian
// LV: Latvian
// RO: Romanian
// SK: Slovakian
// SL: Slovenian
// SV: Swedish
func Translate(sourceLanguage, targetLanguage, textToTranslate string, options ...Option) (jsonRpcResponse *rpc.JsonRpcResponse, err error) {
	opt := &option.DeepLClientOption{}
	for _, apply := range options {
		apply(opt)
	}
	if !opt.UseDeepLPro {
		cli := deeplx.New(opt)
		return cli.Translate(sourceLanguage, targetLanguage, textToTranslate)
	}
	return nil, err
}
