# deepl
## free deepl translate api in go
reference [deepLX](https://github.com/OwO-Network/DeepLX) ,Reverse-engineering the DeepL Windows client has improved several details.

deepL windows Client download Url： https://appdownload.deepl.com/windows/0install/deepl.xml

Example:
```go
import "github.com/czyt/deepl"
translateResp, err := Translate("", "zh", "I love Go programming language")
	if err != nil {
		t.Fatal(err)
	}
log.Println(translateResp.Result.Texts[0].Text)
```
