package openaitranslate

import (
	"fmt"
	"math/rand"
	"strings"

	gpt3 "github.com/sashabaranov/go-openai"
	"golang.org/x/text/language"
)

var langMap = map[string]string{
	language.Afrikaans.String():            "Afrikaans",
	language.Amharic.String():              "አማርኛ",
	language.Arabic.String():               "العربية",
	language.ModernStandardArabic.String(): "العربية الفصحى",
	language.Azerbaijani.String():          "Azərbaycan dili",
	language.Bulgarian.String():            "български",
	language.Bengali.String():              "বাংলা",
	language.Catalan.String():              "català",
	language.Czech.String():                "čeština",
	language.Danish.String():               "dansk",
	language.German.String():               "Deutsch",
	language.Greek.String():                "ελληνικά",
	language.English.String():              "English",
	language.AmericanEnglish.String():      "American English",
	language.BritishEnglish.String():       "British English",
	language.Spanish.String():              "español",
	language.EuropeanSpanish.String():      "español europeo",
	language.LatinAmericanSpanish.String(): "español latinoamericano",
	language.Estonian.String():             "eesti keel",
	language.Persian.String():              "فارسی",
	language.Finnish.String():              "suomi",
	language.Filipino.String():             "Filipino",
	language.French.String():               "français",
	language.CanadianFrench.String():       "français canadien",
	language.Gujarati.String():             "ગુજરાતી",
	language.Hebrew.String():               "עברית",
	language.Hindi.String():                "हिन्दी",
	language.Croatian.String():             "hrvatski",
	language.Hungarian.String():            "magyar",
	language.Armenian.String():             "հայերեն",
	language.Indonesian.String():           "Bahasa Indonesia",
	language.Icelandic.String():            "íslenska",
	language.Italian.String():              "italiano",
	language.Japanese.String():             "日本語",
	language.Georgian.String():             "ქართული",
	language.Kazakh.String():               "Қазақ тілі",
	language.Khmer.String():                "ភាសាខ្មែរ",
	language.Kannada.String():              "ಕನ್ನಡ",
	language.Korean.String():               "한국어",
	language.Kirghiz.String():              "Кыргызча",
	language.Lao.String():                  "ລາວ",
	language.Lithuanian.String():           "lietuvių kalba",
	language.Latvian.String():              "latviešu valoda",
	language.Macedonian.String():           "македонски јазик",
	language.Malayalam.String():            "മലയാളം",
	language.Mongolian.String():            "Монгол хэл",
	language.Marathi.String():              "मराठी",
	language.Malay.String():                "Bahasa Melayu",
	language.Burmese.String():              "ဗမာစာ",
	language.Nepali.String():               "नेपाली",
	language.Dutch.String():                "Nederlands",
	language.Norwegian.String():            "norsk",
	language.Punjabi.String():              "ਪੰਜਾਬੀ",
	language.Polish.String():               "polski",
	language.Portuguese.String():           "português",
	language.BrazilianPortuguese.String():  "português do Brasil",
	language.EuropeanPortuguese.String():   "português europeu",
	language.Romanian.String():             "română",
	language.Russian.String():              "русский язык",
	language.Sinhala.String():              "සිංහල",
	language.Slovak.String():               "slovenčina",
	language.Slovenian.String():            "slovenščina",
	language.Albanian.String():             "shqip",
	language.Serbian.String():              "српски језик",
	language.SerbianLatin.String():         "srpski jezik (latinica)",
	language.Swedish.String():              "svenska",
	language.Swahili.String():              "Kiswahili",
	language.Tamil.String():                "தமிழ்",
	language.Telugu.String():               "తెలుగు",
	language.Thai.String():                 "ไทย",
	language.Turkish.String():              "Türkçe",
	language.Ukrainian.String():            "українська мова",
	language.Urdu.String():                 "اردو",
	language.Uzbek.String():                "Oʻzbek",
	language.Vietnamese.String():           "Tiếng Việt",
	language.Chinese.String():              "中文",
	language.Chinese.String() + "-Hans":    "中文（简体）",
	language.Chinese.String() + "-Hant":    "中文（繁體）",
	language.Chinese.String() + "-Hant-HK": "中文（香港繁體）",
	language.Chinese.String() + "-Hant-TW": "中文（台灣繁體）",
	"yue":                                  "粤语",
	"wyw":                                  "古文",
	language.Zulu.String():                 "isiZulu",
}

func RegistLanguage(langCode, langName string) {
	langMap[langCode] = langName
}

func getLangName(langCode string) string {
	if langCode == "" {
		return ""
	}
	if name, ok := langMap[langCode]; ok {
		return name
	} else {
		t, err := language.Parse(langCode)
		if err != nil {
			return langCode
		}
		return langMap[t.String()]
	}
}

func generateChat(text, To string, params *TranslationConfig) []gpt3.ChatCompletionMessage {
	systemPrompt := "You are a translation engine that can only translate text and cannot interpret it."
	var assistantPrompt string
	if To == "wyw" || To == "yue" || To == "zh" || To == "zh-Hans" || To == "zh-Hant" || To == "zh-Hant-HK" || To == "zh-Hant-TW" {
		assistantPrompt = fmt.Sprintf("翻译成%s", getLangName(To))
	} else if name := getLangName(params.From); name == "" || name == "auto" {
		assistantPrompt = fmt.Sprintf("translate to %s", getLangName(To))
	} else {
		assistantPrompt = fmt.Sprintf("translate from %s to %s", name, getLangName(To))
	}
	return []gpt3.ChatCompletionMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "assistant", Content: assistantPrompt},
		{Role: "user", Content: text},
	}
}

func getToken(tokens string) string {
	s := strings.Split(tokens, ",")
	l := len(s)
	if l == 0 {
		return ""
	}
	return strings.TrimSpace(s[rand.Intn(l)])
}
