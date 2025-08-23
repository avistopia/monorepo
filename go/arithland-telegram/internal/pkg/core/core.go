package core

import (
	"fmt"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/models"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/clients"
	"github.com/avistopia/monorepo/go/arithland-telegram/internal/pkg/flows"
	"gorm.io/gorm"
)

type Service struct {
	telegram         clients.Telegram
	adminUsernames   map[string]struct{}
	defaultPhotoID   string
	db               *gorm.DB
	userRepo         *models.UserRepo
	questionRepo     *models.QuestionRepo
	userQuestionRepo *models.UserQuestionRepo
}

func NewService(
	telegram clients.Telegram,
	adminUsernames map[string]struct{},
	defaultPhotoID string,
	db *gorm.DB,
	userRepo *models.UserRepo,
	questionRepo *models.QuestionRepo,
	userQuestionRepo *models.UserQuestionRepo,
) *Service {
	return &Service{
		telegram:         telegram,
		adminUsernames:   adminUsernames,
		defaultPhotoID:   defaultPhotoID,
		db:               db,
		userRepo:         userRepo,
		questionRepo:     questionRepo,
		userQuestionRepo: userQuestionRepo,
	}
}

func (s *Service) Flow() (*flows.Flow, error) {
	flow, err := flows.MergeFlows([]flows.Flow{
		s.mainMenuFlow(),
		s.arithlandConstitutionFlow(),
		s.profileManagementFlow(),
		s.questionsManagementFlow(),
		s.AdminQuestionsFlow(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to merge flows: %w", err)
	}

	return flow, nil
}
