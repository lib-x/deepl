package translator

import "github.com/lib-x/deepl/internal/rpc"

type Translator interface {
	Translate(sourceLanguage, targetLanguage, textToTranslate string) (jsonRpcResponse *rpc.JsonRpcResponse, err error)
}
