package openaitranslate

import (
	"context"

	gpt3 "github.com/sashabaranov/go-openai"
)

type Translated struct {
	// Detected      Detected `json:"detected"`
	Text string `json:"text"` // translated text
	// Pronunciation string   `json:"pronunciation"` // pronunciation of translated text
}

// Detected represents language detection result
type Detected struct {
	Lang       string  `json:"lang"`       // detected language
	Confidence float64 `json:"confidence"` // the confidence of detection result (0.00 to 1.00)
}

type TranslationConfig struct {
	Ctx                 context.Context
	MaxTokens           int
	Temperature         float32 // 0-2, 越高越随机
	TopP                float32 // 0.1表示仅考虑包含最高前10%概率质量的令牌,推荐1.0
	PresencePenalty     float32 // 介于-2.0和2.0之间的数字。正值会根据新标记到目前为止是否出现在文本中来惩罚它们，从而增加模型谈论新主题的可能性。
	FrequencyPenalty    float32 // 介于-2.0和2.0之间的数字。正值会根据新符号在文本中的现有频率来惩罚它们，从而降低模型逐字重复同一行的可能性。
	From, SelectedWords string
}

const (
	DefaultMaxTokens        = 1000
	DefaultTemperature      = 0.0
	DefaultTopP             = 1.0
	DefaultPresencePenalty  = 1.0
	DefaultFrequencyPenalty = 1.0
)

func (cfg *TranslationConfig) correct() {
	if cfg.Ctx == nil {
		cfg.Ctx = context.Background()
	}
	if cfg.MaxTokens < 0 || cfg.MaxTokens > 4096 {
		cfg.MaxTokens = DefaultMaxTokens
	}
	if cfg.Temperature < 0 || cfg.Temperature > 2 {
		cfg.Temperature = DefaultTemperature
	}
	if cfg.PresencePenalty < -2 || cfg.PresencePenalty > 2 {
		cfg.PresencePenalty = DefaultPresencePenalty
	}
	if cfg.FrequencyPenalty < -2 || cfg.FrequencyPenalty > 2 {
		cfg.FrequencyPenalty = DefaultFrequencyPenalty
	}
}

func Translate(text, To, Token string, cfg TranslationConfig) (string, error) {
	c := gpt3.NewClient(Token)
	cfg.correct()
	resp, err := c.CreateChatCompletion(cfg.Ctx, gpt3.ChatCompletionRequest{
		Model:            gpt3.GPT3Dot5Turbo0301,
		MaxTokens:        1000,
		Temperature:      0,
		TopP:             1,
		PresencePenalty:  1,
		FrequencyPenalty: 1,

		Messages: generateChat(text, To, &cfg),
	})
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}
