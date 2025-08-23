package core

import (
	"fmt"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/models"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/components"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/flows"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/texts"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *Service) mainMenuFlow() flows.Flow {
	return flows.Flow{
		CommandActions: map[string]components.Action{
			"start": func(user *models.User, message *tgbotapi.Message) error {
				if err := s.telegram.SendOrEdit(
					components.Message{Text: texts.MainMenu_WelcomeToArithland, Keyboard: s.defaultKeyboard(message)},
					message.Chat.ID,
					nil,
				); err != nil {
					return fmt.Errorf("failed to send welcome to arithland message: %w", err)
				}

				return nil
			},
		},
		MessageActions: map[models.StateName]components.Action{
			models.StateName_Default: func(user *models.User, message *tgbotapi.Message) error {
				if err := s.telegram.SendOrEdit(
					components.Message{Text: texts.MainMenu_DefaultResponse, Keyboard: s.defaultKeyboard(message)},
					message.Chat.ID,
					nil,
				); err != nil {
					return fmt.Errorf("failed to send default response message: %w", err)
				}

				return nil
			},
		},
	}
}

func (s *Service) defaultKeyboard(message *tgbotapi.Message) components.Keyboard {
	keyboard := components.Keyboard{
		{
			components.NewKeyboardButton(texts.ArithlandConstitution_Show),
		},
		{
			components.NewKeyboardButton(texts.ProfileManagement_Show),
			components.NewKeyboardButton(texts.QuestionsManagement_Show),
		},
	}

	if s.isAdmin(message.Chat.UserName) {
		keyboard = append(keyboard, []components.KeyboardButton{
			components.NewKeyboardButton(texts.AdminQuestions_Show),
		})
	}

	return keyboard
}
