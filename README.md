# deepl
## free deepl translate api in go
[![Go Report Card](https://goreportcard.com/badge/github.com/czyt/deepl)](https://goreportcard.com/report/github.com/czyt/deepl)

reference [deepLX](https://github.com/OwO-Network/DeepLX) ,Reverse-engineering the DeepL Windows client has improved several details.

deepL windows Client download Url： https://appdownload.deepl.com/windows/0install/deepl.xml

Example:
```go
import "github.com/tiny-lib/deepl"
translateResp, err := Translate("", "zh", "I love Go programming language")
	if err != nil {
		t.Fatal(err)
	}
log.Println(translateResp.Result.Texts[0].Text)
```
