package deepl

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
