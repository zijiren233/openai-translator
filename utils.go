package openaitranslate

import (
	"fmt"
	"log"
	"sync"

	gpt3 "github.com/sashabaranov/go-openai"
	"golang.org/x/text/language"
)

var langMapLock = &sync.RWMutex{}

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
	language.Chinese.String():              "中文（简体）",
	language.SimplifiedChinese.String():    "中文（简体）",
	language.TraditionalChinese.String():   "中文（繁體）",
	"yue":                                  "中文（粤语）",
	"wyw":                                  "中文（古文-文言文）",
	language.Zulu.String():                 "isiZulu",
}

func RegistLanguage(langCode, langName string) {
	langMapLock.Lock()
	defer langMapLock.Unlock()
	langMap[langCode] = langName
}

func GetLangMap() map[string]string {
	langMapLock.RLock()
	defer langMapLock.RUnlock()
	return copyMap(langMap)
}

func copyMap[K, V comparable](m map[K]V) map[K]V {
	newMap := make(map[K]V)
	for k, v := range m {
		newMap[k] = v
	}
	return newMap
}

func getLangName(langCode string) string {
	langMapLock.RLock()
	defer langMapLock.RUnlock()
	if langCode == "" {
		return ""
	}
	if name, ok := langMap[langCode]; ok {
		return name
	} else {
		if code, ok := langMap[getBaseLangCode(langCode)]; ok {
			return code
		} else {
			return langCode
		}
	}
}

func getBaseLangCode(langCode string) string {
	t, err := language.Parse(langCode)
	if err != nil {
		return langCode
	}
	if parent := t.Parent().String(); parent == "und" {
		return t.String()
	} else {
		return getBaseLangCode(parent)
	}
}

func generateChat(text, To string, params *TranslationConfig) []gpt3.ChatCompletionMessage {
	systemPrompt := "You are a translation engine that can only translate text and cannot interpret it."
	var assistantPrompt string
	To = getBaseLangCode(To)
	if name := getLangName(params.From); name == "" || name == "auto" {
		if To == "wyw" || To == "yue" || To == "zh" || To == "zh-Hans" || To == "zh-Hant" {
			assistantPrompt = fmt.Sprintf("我说的下一句话翻译成%s", getLangName(To))
		} else {
			assistantPrompt = fmt.Sprintf("Translate my next sentence to %s", getLangName(To))
		}
	} else {
		if To == "wyw" || To == "yue" || To == "zh" || To == "zh-Hans" || To == "zh-Hant" {
			assistantPrompt = fmt.Sprintf("我说的下一句话从%s翻译成%s", name, getLangName(To))
		} else {
			assistantPrompt = fmt.Sprintf("Translate my next sentence from %s to %s", name, getLangName(To))
		}
	}
	chat := []gpt3.ChatCompletionMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: assistantPrompt},
		{Role: "user", Content: text},
	}
	if params.Debug {
		log.Println(chat)
	}
	return chat
}
