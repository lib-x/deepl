package deepl

import "testing"

func TestTranslate(t *testing.T) {
	translate, err := Translate("", "zh", "I love Go programming language")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(translate.Result.Texts[0].Text)
	t.Log(len(translate.Result.Texts[0].Alternatives))
}

func TestTranslateWithProSession(t *testing.T) {
	options := []Option{
		WithDeeplProSession("fa.e4xxxxxxxxxxxxxxxxxxx"),
	}
	translate, err := Translate("", "zh", "I am the apple of my father's eyes", options...)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(translate.Result.Texts[0].Text)
	t.Log(len(translate.Result.Texts[0].Alternatives))
}
func TestTranslateWithHttpProxy(t *testing.T) {
	options := []Option{
		WithHttpProxy("http://127.0.0.1:10808"),
	}
	translate, err := Translate("", "zh", "I love Go programming language", options...,
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(translate.Result.Texts[0].Text)
	t.Log(len(translate.Result.Texts[0].Alternatives))
}

func TestTranslateWithSocket5Proxy(t *testing.T) {
	options := []Option{
		WithSocket5Proxy("127.0.0.1:10808", "", ""),
	}
	translate, err := Translate("", "zh", "I love Go programming language", options...)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(translate.Result.Texts[0].Text)
	t.Log(len(translate.Result.Texts[0].Alternatives))
}
