package core_test

import (
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/models"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/clients/clientsmocks"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/core"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/test"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite
}

func TestServiceSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ServiceSuite))
}

func (s *ServiceSuite) TestFlow() {
	db, err := test.DB()
	s.Require().NoError(err)

	telegram := clientsmocks.NewMockTelegram(s.T())

	userRepo, err := models.NewUserRepo(db)
	s.Require().NoError(err)

	userQuestionRepo, err := models.NewUserQuestionRepo(db)
	s.Require().NoError(err)

	questionRepo, err := models.NewQuestionRepo(db)
	s.Require().NoError(err)

	_, err = core.NewService(telegram, map[string]struct{}{}, "", db, userRepo, questionRepo, userQuestionRepo).Flow()
	s.Require().NoError(err)
}
