package core

import (
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/components"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/texts"
)

func (s *Service) isAdmin(username string) bool {
	_, isAdmin := s.adminUsernames[username]
	return isAdmin
}
func (s *Service) forbiddenMessage() components.Message {
	return components.Message{Text: texts.Err_Forbidden}
}
