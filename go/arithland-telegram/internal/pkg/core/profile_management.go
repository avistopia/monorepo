package core

import (
	"fmt"
	"log"

	"github.com/avistopia/monorepo/go/arithland-telegram/internal/models"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/clean"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/components"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/flows"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/texts"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/values"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *Service) profileManagementFlow() flows.Flow {
	return flows.Flow{
		KeyboardButtonActions: map[string]components.Action{
			texts.ProfileManagement_Show: func(user *models.User, message *tgbotapi.Message) error {
				if err := s.telegram.SendOrEdit(s.profileManagementMessage(user), message.Chat.ID, nil); err != nil {
					return fmt.Errorf("failed to send profile management message: %w", err)
				}

				return nil
			},
		},
		InlineButtonActions: map[components.InlineButtonActionName]components.InlineButtonAction{
			actionName_profileManagement_changeDisplayName: func(
				user *models.User, query *tgbotapi.CallbackQuery, data string,
			) (string, error) {
				user.State = models.NewWaitingForUserFieldState(models.UserFieldName_DisplayName)

				if err := s.userRepo.Save(user); err != nil {
					return "", fmt.Errorf("failed to save user: %w", err)
				}

				if err := s.telegram.SendOrEdit(
					s.enterDisplayNameMessage(), query.Message.Chat.ID, values.Ptr(query.Message.MessageID),
				); err != nil {
					log.Printf("failed to edit message: %v", err)
				}

				return texts.ProfileManagement_ChangeDisplayName, nil
			},
			actionName_profileManagement_backToShowProfileManagement: func(
				user *models.User, query *tgbotapi.CallbackQuery, data string,
			) (string, error) {
				user.State = models.NewDefaultState()

				if err := s.userRepo.Save(user); err != nil {
					return "", fmt.Errorf("failed to save user: %w", err)
				}

				if err := s.telegram.SendOrEdit(
					s.profileManagementMessage(user), query.Message.Chat.ID, values.Ptr(query.Message.MessageID),
				); err != nil {
					log.Printf("failed to edit message: %v", err)
				}

				return texts.Common_Cancelled, nil
			},
		},
		MessageActions: map[models.StateName]components.Action{
			models.StateName_WaitingForUserField: func(user *models.User, message *tgbotapi.Message) error {
				switch user.State.Data.WaitingForUserField.FieldName {
				case models.UserFieldName_DisplayName:
					displayName, validationError := clean.UserDisplayName(message.Text)
					if validationError != "" {
						if err := s.telegram.SendOrEdit(
							components.Message{Text: validationError}, message.Chat.ID, nil,
						); err != nil {
							return fmt.Errorf("failed to send validation error message: %w", err)
						}

						return nil
					}

					user.DisplayName = displayName
				default:
					return fmt.Errorf("unknown field name: %s", user.State.Data.WaitingForUserField.FieldName)
				}

				user.State = models.NewDefaultState()

				if err := s.userRepo.Save(user); err != nil {
					return fmt.Errorf("failed to save user: %w", err)
				}

				if err := s.telegram.SendOrEdit(s.profileManagementMessage(user), message.Chat.ID, nil); err != nil {
					return fmt.Errorf("failed to send profile management message: %w", err)
				}

				return nil
			},
		},
	}
}

func (s *Service) profileManagementMessage(user *models.User) components.Message {
	return components.Message{
		Text: texts.Format(texts.ProfileManagement, map[string]string{
			"displayName": user.DisplayName,
			"balance":     texts.FormatFloat(user.Balance),
		}),
		InlineKeyboard: components.InlineKeyboard{
			{
				components.NewInlineButton(
					texts.ProfileManagement_ChangeDisplayName, actionName_profileManagement_changeDisplayName, "",
				),
			},
		},
	}
}

func (s *Service) enterDisplayNameMessage() components.Message {
	return components.Message{
		Text: texts.ProfileManagement_EnterDisplayName,
		InlineKeyboard: components.InlineKeyboard{
			{
				components.NewInlineButton(
					texts.Common_Cancel, actionName_profileManagement_backToShowProfileManagement, "",
				),
			},
		},
	}
}
