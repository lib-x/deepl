package rpc

func New(sourceLang string, targetLang string) *JsonRpcRequest {
	return &JsonRpcRequest{
		Jsonrpc: "2.0",
		Method:  "LMT_handle_texts",
		Params: TextTranslationJsonRpcRequestParams{
			Splitting: "newlines",
			Lang: Lang{
				SourceLangUserSelected: sourceLang,
				TargetLang:             targetLang,
			},
			CommonJobParams: CommonJobParameters{
				WasSpoken:    false,
				TranscribeAS: "",
			},
		},
	}
}
