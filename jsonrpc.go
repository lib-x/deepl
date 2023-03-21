package deepl

type JsonRpcResponse struct {
	Jsonrpc   string                      `json:"jsonrpc"`
	Id        int64                       `json:"id"`
	Result    TextTranslationJsonResponse `json:"result"`
	ErrorInfo *ErrorInfo                  `json:"error,omitempty"`
}
type JsonRpcRequest struct {
	Jsonrpc string                              `json:"jsonrpc"`
	Method  string                              `json:"method"`
	Id      int64                               `json:"id"`
	Params  TextTranslationJsonRpcRequestParams `json:"params"`
}
