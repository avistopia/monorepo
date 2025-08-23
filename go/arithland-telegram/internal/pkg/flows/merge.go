package flows

import (
	"fmt"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/models"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/components"
)

func MergeFlows(flows []Flow) (*Flow, error) {
	flow := &Flow{
		CommandActions:        map[string]components.Action{},
		MessageActions:        map[models.StateName]components.Action{},
		InlineButtonActions:   map[components.InlineButtonActionName]components.InlineButtonAction{},
		KeyboardButtonActions: map[string]components.Action{},
	}

	for i := range flows {
		err := mergeMapsWithNoDuplicate(flow.CommandActions, flows[i].CommandActions)
		if err != nil {
			return nil, fmt.Errorf("failed to aggregate command actions: %w", err)
		}

		err = mergeMapsWithNoDuplicate(flow.MessageActions, flows[i].MessageActions)
		if err != nil {
			return nil, fmt.Errorf("failed to aggregate message actions: %w", err)
		}

		err = mergeMapsWithNoDuplicate(flow.InlineButtonActions, flows[i].InlineButtonActions)
		if err != nil {
			return nil, fmt.Errorf("failed to aggregate inline button actions: %w", err)
		}

		err = mergeMapsWithNoDuplicate(flow.KeyboardButtonActions, flows[i].KeyboardButtonActions)
		if err != nil {
			return nil, fmt.Errorf("failed to aggregate keyboard button actions: %w", err)
		}
	}

	return flow, nil
}

func mergeMapsWithNoDuplicate[K comparable, V any](destination map[K]V, source map[K]V) error {
	for k, v := range source {
		if _, exists := destination[k]; exists {
			return fmt.Errorf("duplicate key '%v'", k)
		}

		destination[k] = v
	}

	return nil
}
