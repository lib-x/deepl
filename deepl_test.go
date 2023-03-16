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
