package deepl

import (
	"bytes"
	"encoding/json"
	"github.com/abadojack/whatlanggo"
	"github.com/andybalholm/brotli"
	"io"
	"math/big"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	deepLRpcServer   = "https://www2.deepl.com/jsonrpc"
	deepLReferServer = "https://www.deepl.com/"
)

var (
	methodPartNormal       = []byte("\"method\":\"")
	methodPartWithOneSpace = []byte("\"method\": \"")
	methodPartWithTwoSpace = []byte("\"method\" : \"")
)

func newJsonRpcRequest(sourceLang string, targetLang string) *JsonRpcRequest {
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
				// RegionalVariant: "en-US",
			},
		},
	}
}

func countIs(translateText string) int64 {
	return int64(strings.Count(translateText, "i"))
}

func generateNextId() int64 {
	randProvider := rand.NewSource(time.Now().Unix())
	nextId := big.NewInt(0).Sqrt(big.NewInt(10000))
	random := big.NewInt(randProvider.Int63())
	nextId = nextId.Mul(nextId, random)
	return nextId.Int64()
}

func generateTimestamp(iCount int64) int64 {
	ts := time.Now().UnixMilli()
	if iCount != 0 {
		iCount = iCount + 1
		return ts - ts%iCount + iCount
	} else {
		return ts
	}
}

func adjustJsonContent(id int64, jsonContent []byte) []byte {
	// add space if necessary
	if (id+5)%29 == 0 || (id+3)%13 == 0 {
		jsonContent = bytes.ReplaceAll(jsonContent, methodPartNormal, methodPartWithTwoSpace)
	} else {
		jsonContent = bytes.ReplaceAll(jsonContent, methodPartNormal, methodPartWithOneSpace)
	}
	return jsonContent
}

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
func Translate(sourceLanguage, targetLanguage, textToTranslate string, options ...Option) (jsonRpcResponse *JsonRpcResponse, err error) {

	clientOpt := &deepLClientOption{}
	if len(options) > 0 {
		for _, optFunc := range options {
			optFunc(clientOpt)
		}
	}
	if sourceLanguage == "" {
		lang := whatlanggo.DetectLang(textToTranslate)
		deepLLang := strings.ToUpper(lang.Iso6391())
		sourceLanguage = deepLLang
	}
	if targetLanguage == "" {
		targetLanguage = "EN"
	}
	if textToTranslate == "" {
		return nil, ErrorNoTranslateTextFound
	}
	postData := newJsonRpcRequest(sourceLanguage, targetLanguage)
	text := TextParam{
		Text:                Text{Text: textToTranslate},
		RequestAlternatives: 3,
	}
	// set id
	id := generateNextId() + 1
	postData.Id = id
	// set text
	postData.Params.Texts = append(postData.Params.Texts, text)
	// set timestamp
	postData.Params.Timestamp = generateTimestamp(countIs(textToTranslate))
	postByte, _ := json.Marshal(postData)
	postByte = adjustJsonContent(id, postByte)
	reader := bytes.NewReader(postByte)
	request, err := http.NewRequest("POST", deepLRpcServer, reader)
	if err != nil {
		return nil, err
	}

	// Set Headers
	request.Header.Set("Referer", deepLReferServer)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "*/*")
	request.Header.Set("x-app-os-name", "iOS")
	request.Header.Set("x-app-os-version", "16.3.0")
	request.Header.Set("Accept-Language", "en-US,en;q=0.9")
	request.Header.Set("Accept-Encoding", "gzip, deflate, br")
	request.Header.Set("x-app-device", "iPhone13,2")
	request.Header.Set("User-Agent", "DeepL-iOS/2.9.1 iOS 16.3.0 (iPhone13,2)")
	request.Header.Set("x-app-build", "510265")
	request.Header.Set("x-app-version", "2.9.1")
	request.Header.Set("Connection", "keep-alive")

	client := &http.Client{}
	if clientOpt.httpProxy != "" {
		httpProxy, _ := url.Parse(clientOpt.httpProxy)
		if httpProxy != nil {
			client.Transport = &http.Transport{Proxy: http.ProxyURL(httpProxy)}
		}
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, ErrorInvalidResponse
	}
	var bodyReader io.Reader
	switch resp.Header.Get("Content-Encoding") {
	case "br":
		bodyReader = brotli.NewReader(resp.Body)
	default:
		bodyReader = resp.Body
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return jsonRpcResponse, ErrorTooManyRequests
	}

	jsonRpcResponse = &JsonRpcResponse{}
	if err := json.NewDecoder(bodyReader).Decode(jsonRpcResponse); err != nil {
		return nil, err
	}
	if jsonRpcResponse.ErrorInfo != nil {
		if jsonRpcResponse.ErrorInfo.ErrorCode == -32600 {
			return nil, ErrorInvalidTargetLang
		}
	}
	return jsonRpcResponse, nil
}
