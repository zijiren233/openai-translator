package openaitranslator

import (
	"fmt"
	"log"

	"github.com/sashabaranov/go-openai"
)

func Translate(text, To, Token string, opt ...Option) (string, error) {
	cfg := DefaultConfig()
	for _, v := range opt {
		v(cfg)
	}
	return TranslateWithConfig(text, To, Token, cfg)
}

func TranslateWithConfig(text, To, Token string, cfg *TranslationConfig) (string, error) {
	url, err := parseOpenaiAPIURLv1(cfg.Url)
	if err != nil {
		return "", err
	}
	cfg.correct()
	openaiConf := openai.DefaultConfig(Token)
	if url != "" {
		openaiConf.BaseURL = url
	}
	resp, err := openai.NewClientWithConfig(openaiConf).CreateChatCompletion(cfg.Ctx, openai.ChatCompletionRequest{
		Model:            cfg.Model,
		MaxTokens:        cfg.MaxTokens,
		Temperature:      cfg.Temperature,
		TopP:             cfg.TopP,
		PresencePenalty:  cfg.PresencePenalty,
		FrequencyPenalty: cfg.FrequencyPenalty,

		Messages: generateChat(text, To, cfg),
	})
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

func generateChat(text, To string, params *TranslationConfig) []openai.ChatCompletionMessage {
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
	chat := []openai.ChatCompletionMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: assistantPrompt},
		{Role: "user", Content: text},
	}
	if params.Debug {
		log.Println(chat)
	}
	return chat
}
