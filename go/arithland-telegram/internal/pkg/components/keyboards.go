package components

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Keyboard [][]KeyboardButton

func (k Keyboard) Render() *tgbotapi.ReplyKeyboardMarkup {
	return &tgbotapi.ReplyKeyboardMarkup{
		ResizeKeyboard: true,
		Keyboard:       renderTable[tgbotapi.KeyboardButton](k),
	}
}

type InlineKeyboard [][]InlineButton

func (k InlineKeyboard) Render() *tgbotapi.InlineKeyboardMarkup {
	return &tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: renderTable[tgbotapi.InlineKeyboardButton](k),
	}
}
