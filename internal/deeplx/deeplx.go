package deeplx

import (
	"bytes"
	"encoding/json"
	"github.com/abadojack/whatlanggo"
	"github.com/andybalholm/brotli"
	"github.com/lib-x/deepl/internal/headers"
	"github.com/lib-x/deepl/internal/innererrors"
	"github.com/lib-x/deepl/internal/kits"
	"github.com/lib-x/deepl/internal/option"
	"github.com/lib-x/deepl/internal/postprocess"
	"github.com/lib-x/deepl/internal/rpc"
	"github.com/lib-x/deepl/internal/translator"
	"io"
	"net/http"
	"strings"
)

var _ translator.Translator = (*Client)(nil)

const (
	deepLXRpcServer = "https://www2.deepl.com/jsonrpc"
)

type Client struct {
	opt *option.DeepLClientOption
}

func New(opt *option.DeepLClientOption) *Client {
	return &Client{
		opt: opt,
	}
}

func (d *Client) Translate(sourceLanguage, targetLanguage, textToTranslate string) (jsonRpcResponse *rpc.JsonRpcResponse, err error) {
	var (
		requestServer = deepLXRpcServer
		headerIter    = headers.DeeplXHeaderIter()
	)
	if sourceLanguage == "" {
		lang := whatlanggo.DetectLang(textToTranslate)
		deepLLang := strings.ToUpper(lang.Iso6391())
		sourceLanguage = deepLLang
	}
	if targetLanguage == "" {
		targetLanguage = "EN"
	}
	if textToTranslate == "" {
		return nil, innererrors.ErrorNoTranslateTextFound
	}
	postData := rpc.New(sourceLanguage, targetLanguage)
	text := rpc.TextParam{
		Text:                rpc.Text{Text: textToTranslate},
		RequestAlternatives: 3,
	}
	if d.opt.TagHandling != "" {
		postData.Params.TagHandling = d.opt.TagHandling
	}
	// set id
	id := postprocess.GenerateNextId() + 1
	postData.Id = id
	// set text
	postData.Params.Texts = append(postData.Params.Texts, text)
	// set timestamp
	postData.Params.Timestamp = postprocess.GenerateTimestampForContent(textToTranslate)
	postByte, _ := json.Marshal(postData)
	postByte = postprocess.AdjustContent(id, postByte)
	reader := bytes.NewReader(postByte)

	request, err := http.NewRequest("POST", requestServer, reader)
	if err != nil {
		return nil, err
	}
	// Set Headers

	for headerKey, headValue := range headerIter {
		request.Header.Set(headerKey, headValue)
	}

	client := &http.Client{}
	if transport := kits.BuildHttpTransportWith(d.opt); transport != nil {
		client.Transport = transport
	}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, innererrors.ErrorInvalidResponse
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
		return jsonRpcResponse, innererrors.ErrorTooManyRequests
	}

	jsonRpcResponse = &rpc.JsonRpcResponse{}
	if err := json.NewDecoder(bodyReader).Decode(jsonRpcResponse); err != nil {
		return nil, err
	}
	if jsonRpcResponse.ErrorInfo != nil {
		if jsonRpcResponse.ErrorInfo.ErrorCode == -32600 {
			return nil, innererrors.ErrorInvalidTargetLang
		}
	}
	return jsonRpcResponse, nil
}
