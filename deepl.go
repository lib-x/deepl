package deepl

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"github.com/abadojack/whatlanggo"
	"github.com/andybalholm/brotli"
	"golang.org/x/net/context"
	"golang.org/x/net/proxy"
	"io"
	"math/big"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	deepLRpcServer  = "https://api.deepl.com/jsonrpc"
	deepLXRpcServer = "https://www2.deepl.com/jsonrpc"
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
	var (
		requestServer = deepLXRpcServer
		headerIter    = deeplXHeaderIter()
	)

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
	if clientOpt.tagHandling != "" {
		postData.Params.TagHandling = clientOpt.tagHandling
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
	if clientOpt.useDeepLPro {
		requestServer = deepLRpcServer
		headerIter = deeplProHeaderIter(clientOpt.dlSession)
	}
	request, err := http.NewRequest("POST", requestServer, reader)
	if err != nil {
		return nil, err
	}
	// Set Headers

	for headerKey, headValue := range headerIter {
		request.Header.Set(headerKey, headValue)
	}

	client := &http.Client{}
	if transport := createProxyTransportWith(clientOpt); transport != nil {
		client.Transport = transport
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

func createProxyTransportWith(clientOpt *deepLClientOption) *http.Transport {
	var transport *http.Transport
	if clientOpt.httpProxy != "" {
		httpProxy, _ := url.Parse(clientOpt.httpProxy)
		if httpProxy != nil {
			transport = &http.Transport{Proxy: http.ProxyURL(httpProxy)}
		}
	}
	if clientOpt.socket5Proxy != "" {
		var auth *proxy.Auth
		if clientOpt.socket5ProxyUser != "" || clientOpt.socket5proxyPassword != "" {
			auth = &proxy.Auth{User: clientOpt.socket5ProxyUser, Password: clientOpt.socket5proxyPassword}
		}
		dialer, err := proxy.SOCKS5("tcp", clientOpt.socket5Proxy, auth, proxy.Direct)
		if err == nil {
			dialContext := func(ctx context.Context, network, address string) (net.Conn, error) {
				return dialer.Dial(network, address)
			}
			transport = &http.Transport{DialContext: dialContext}
		}
	}
	if clientOpt.ignoreSSLVerification && transport != nil {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	return nil
}
