package clients

import (
	"fmt"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/components"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type Telegram interface {
	SendOrEdit(message components.Message, chatID int64, originalMessageID *int) error
	AnswerCallback(text string, callbackQueryID string) error

	GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel
	StopReceivingUpdates()
}

type telegram struct {
	bot *tgbotapi.BotAPI
}

func NewTelegram(token string) (Telegram, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize telegram bot: %w", err)
	}

	return &telegram{bot: bot}, nil
}

func (c *telegram) SendOrEdit(message components.Message, chatID int64, originalMessageID *int) error {
	response, err := message.Render(chatID, originalMessageID)
	if err != nil {
		return fmt.Errorf("failed to render: %w", err)
	}

	if _, err = c.bot.Send(response); err != nil {
		if originalMessageID != nil && strings.Contains(err.Error(), "message is not modified") {
			return nil
		}

		return fmt.Errorf("failed to send: %w", err)
	}

	return nil
}

func (c *telegram) AnswerCallback(text string, callbackQueryID string) error {
	if _, err := c.bot.Send(&tgbotapi.CallbackConfig{
		Text:            text,
		CallbackQueryID: callbackQueryID,
	}); err != nil {
		return fmt.Errorf("failed to send: %w", err)
	}

	return nil
}

func (c *telegram) GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel {
	return c.bot.GetUpdatesChan(config)
}

func (c *telegram) StopReceivingUpdates() {
	c.bot.StopReceivingUpdates()
}
