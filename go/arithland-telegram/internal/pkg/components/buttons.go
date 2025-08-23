package components

import (
	"fmt"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/values"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type KeyboardButton = Renderer[tgbotapi.KeyboardButton]

type keyboardButton struct {
	Text string
}

func NewKeyboardButton(text string) KeyboardButton {
	return &keyboardButton{Text: text}
}

func (b *keyboardButton) Render() tgbotapi.KeyboardButton {
	return tgbotapi.KeyboardButton{Text: b.Text}
}

type InlineButton = Renderer[tgbotapi.InlineKeyboardButton]

type inlineButton struct {
	Text string

	ActionName InlineButtonActionName

	// Data can be maximum of 64 bytes
	Data string
}

func NewInlineButton(text string, actionName InlineButtonActionName, data string) InlineButton {
	return &inlineButton{
		Text:       text,
		ActionName: actionName,
		Data:       data,
	}
}

func (b *inlineButton) Render() tgbotapi.InlineKeyboardButton {
	return tgbotapi.InlineKeyboardButton{
		Text:         b.Text,
		CallbackData: values.Ptr(fmt.Sprintf("%s:%s", b.ActionName, b.Data)),
	}
}
