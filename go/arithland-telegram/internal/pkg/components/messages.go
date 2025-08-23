package components

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Message struct {
	Text           string
	PhotoID        string
	InlineKeyboard InlineKeyboard
	Keyboard       Keyboard
}

func (m Message) Render(chatID int64, originalMessageID *int) (tgbotapi.Chattable, error) {
	if originalMessageID != nil {
		return m.renderEdit(chatID, originalMessageID)
	}

	return m.renderSend(chatID)
}

func (m Message) renderEdit(chatID int64, originalMessageID *int) (tgbotapi.Chattable, error) {
	var replyMarkup *tgbotapi.InlineKeyboardMarkup

	switch {
	case m.Keyboard != nil:
		return nil, fmt.Errorf("cannot update keyboard in edit message")
	case m.InlineKeyboard != nil:
		replyMarkup = m.InlineKeyboard.Render()
	}

	if m.PhotoID != "" {
		return tgbotapi.EditMessageMediaConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:      chatID,
				MessageID:   *originalMessageID,
				ReplyMarkup: replyMarkup,
			},
			Media: tgbotapi.InputMediaPhoto{
				BaseInputMedia: tgbotapi.BaseInputMedia{
					Type:    "photo",
					Media:   tgbotapi.FileID(m.PhotoID),
					Caption: m.Text,
				},
			},
		}, nil
	}

	return tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      chatID,
			MessageID:   *originalMessageID,
			ReplyMarkup: replyMarkup,
		},
		Text: m.Text,
	}, nil
}

func (m Message) renderSend(chatID int64) (tgbotapi.Chattable, error) {
	var replyMarkup any

	switch {
	case m.Keyboard != nil && m.InlineKeyboard != nil:
		return nil, fmt.Errorf("cannot have both keyboard and inline keyboard defined on a message")
	case m.Keyboard != nil:
		replyMarkup = m.Keyboard.Render()
	case m.InlineKeyboard != nil:
		replyMarkup = m.InlineKeyboard.Render()
	}

	if m.PhotoID != "" {
		return tgbotapi.PhotoConfig{
			BaseFile: tgbotapi.BaseFile{
				BaseChat: tgbotapi.BaseChat{ChatID: chatID, ReplyMarkup: replyMarkup},
				File:     tgbotapi.FileID(m.PhotoID),
			},
			Caption: m.Text,
		}, nil
	}

	return tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{ChatID: chatID, ReplyMarkup: replyMarkup},
		Text:     m.Text,
	}, nil
}
