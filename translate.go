package openaitranslate

import (
	"errors"

	gpt3 "github.com/sashabaranov/go-openai"
)

var errTokenIsNone = errors.New("token is none")

func Translate(text, To, Token string, opt ...Option) (string, error) {
	cfg := DefaultConfig()
	for _, v := range opt {
		v(cfg)
	}
	return TranslateWithConfig(text, To, Token, cfg)
}

func TranslateWithConfig(text, To, Token string, cfg *TranslationConfig) (string, error) {
	if Token == "" {
		return "", errTokenIsNone
	}
	cfg.correct()
	resp, err := gpt3.NewClient(Token).CreateChatCompletion(cfg.Ctx, gpt3.ChatCompletionRequest{
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
