package rpc

type Lang struct {
	SourceLangUserSelected string `json:"source_lang_user_selected"`
	TargetLang             string `json:"target_lang"`
}

type CommonJobParameters struct {
	WasSpoken       bool   `json:"wasSpoken,omitempty"`
	TranscribeAS    string `json:"transcribe_as,omitempty"`
	RegionalVariant string `json:"regionalVariant,omitempty"`
}
