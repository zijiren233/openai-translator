package openaitranslator

import (
	"os"
	"testing"
)

func TestTranslate(t *testing.T) {
	for _, unit := range []struct {
		text, from, to string
	}{
		{`Oh yeah! I'm a translator!`, "", "zh"},
		{`Oh yeah! I'm a translator!`, "", "wyw"},
		{`Oh yeah! I'm a translator!`, "", "ja"},
		{`Oh yeah! I'm a translator!`, "", "de"},
		{`Oh yeah! I'm a translator!`, "", "fr"},
	} {
		result, err := Translate(unit.text, unit.to, os.Getenv("OPENAI_API_KEY"), WithFrom(unit.from), WithDebug())
		if err != nil {
			t.Fatal(err)
		}
		t.Log(result)
	}
}

func TestRegistLanguage(t *testing.T) {
	t.Log(GetLangMap())
	RegistLanguage("zh-CN", "简体中文")
	t.Log(GetLangMap())
}
