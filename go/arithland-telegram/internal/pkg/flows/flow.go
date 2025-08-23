package flows

import (
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/models"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/components"
)

type Flow struct {
	// CommandActions is a mapping from command name to actions, used to handle incoming commands.
	//   Commands are messages starting with slash (/)
	CommandActions map[string]components.Action

	// MessageActions is a mapping from state name to actions, used to handle incoming messages.
	//   Any message that is not a command.
	MessageActions map[models.StateName]components.Action

	// InlineButtonActions is a mapping from inline button action names to actions, used to handle incoming callbacks.
	InlineButtonActions map[components.InlineButtonActionName]components.InlineButtonAction

	// KeyboardButtonActions is a mapping from keyboard button texts to actions, used to handle incoming callbacks.
	KeyboardButtonActions map[string]components.Action
}
