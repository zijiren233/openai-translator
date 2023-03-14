package openaitranslate

import (
	"os"
	"testing"
)

func TestTranslate(t *testing.T) {
	result, err := Translate("Go是一种语言层面支持并发（Go最大的特色、天生支持并发）\n内置runtime、支持垃圾回收（GC）、静态强类型，快速编译的语言", "en", os.Getenv(`OPENAI_APIKEY`))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Text: %s\n", result)
}

func TestRegistLanguage(t *testing.T) {
	t.Log(GetLangMap())
	RegistLanguage("zh-CN", "简体中文")
	t.Log(GetLangMap())
}
