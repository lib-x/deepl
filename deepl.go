package deepl

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/abadojack/whatlanggo"
	"io"
	"math/big"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const (
	deepLRpcServer   = "https://www2.deepl.com/jsonrpc"
	deepLReferServer = "https://www.deepl.com/"
)

var (
	ErrorNoTranslateTextFound = errors.New("no Translate Text Found")
	ErrorInvalidTargetLang    = errors.New("invalid Target Lang")
	ErrorTooManyRequests      = errors.New("too Many Requests")
)

type Lang struct {
	SourceLangUserSelected string `json:"source_lang_user_selected"`
	TargetLang             string `json:"target_lang"`
}

type CommonJobParameters struct {
	WasSpoken       bool   `json:"wasSpoken,omitempty"`
	TranscribeAS    string `json:"transcribe_as,omitempty"`
	RegionalVariant string `json:"regionalVariant,omitempty"`
}

type TextTranslationJsonRpcRequestParams struct {
	Texts           []TextParam         `json:"texts"`
	Splitting       string              `json:"splitting"`
	Lang            Lang                `json:"lang"`
	Timestamp       int64               `json:"timestamp"`
	CommonJobParams CommonJobParameters `json:"commonJobParams"`
}
type TextTranslationJsonResponse struct {
	Texts               []TextWithAlternatives `json:"texts"`
	LanguageCode        string                 `json:"lang"`
	LanguageIsConfident bool                   `json:"lang_is_confident"`
	DetectedLanguages   map[string]float64     `json:"detectedLanguages"`
}
type Text struct {
	Text string `json:"text"`
}

type ErrorInfo struct {
	ErrorCode int `json:"code"`
}

type TextWithAlternatives struct {
	Text
	Alternatives []Text `json:"alternatives,omitempty"`
}
type TextParam struct {
	Text
	RequestAlternatives int `json:"requestAlternatives,omitempty"`
}

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
				//WasSpoken:    false,
				//TranscribeAS: "",
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
		jsonContent = bytes.ReplaceAll(jsonContent, []byte("\"method\":\""), []byte("\"method\" : \""))
	} else {
		jsonContent = bytes.ReplaceAll(jsonContent, []byte("\"method\":\""), []byte("\"method\": \""))
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
func Translate(sourceLanguage, targetLanguage, textToTranslate string) (jsonRpcResponse *JsonRpcResponse, err error) {
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
	request.Header.Set("User-Agent", "DeepL-iOS/2.6.0 iOS 16.3.0 (iPhone13,2)")
	request.Header.Set("x-app-build", "353933")
	request.Header.Set("x-app-version", "2.6")
	request.Header.Set("Connection", "keep-alive")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return jsonRpcResponse, ErrorTooManyRequests
	}

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	jsonRpcResponse = &JsonRpcResponse{}
	if err = json.Unmarshal(body, jsonRpcResponse); err != nil {
		return nil, err
	}
	if jsonRpcResponse.ErrorInfo != nil {
		if jsonRpcResponse.ErrorInfo.ErrorCode == -32600 {
			return nil, ErrorInvalidTargetLang
		}
	}
	return jsonRpcResponse, nil
}
