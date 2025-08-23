package core

import (
	"fmt"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/models"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/components"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/flows"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/texts"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *Service) arithlandConstitutionFlow() flows.Flow {
	return flows.Flow{
		KeyboardButtonActions: map[string]components.Action{
			texts.ArithlandConstitution_Show: func(user *models.User, incomingMessage *tgbotapi.Message) error {
				if err := s.telegram.SendOrEdit(
					components.Message{Text: texts.ArithlandConstitution}, incomingMessage.Chat.ID, nil,
				); err != nil {
					return fmt.Errorf("failed to send arithland constitution message: %w", err)
				}

				return nil
			},
		},
	}
}
