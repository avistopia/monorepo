package models

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type UserQuestion struct {
	gorm.Model

	UserID     uint      `gorm:"type:uuid;index:user_question_user_id_idx;uniqueIndex:user_question_unique_idx"`
	QuestionID uint      `gorm:"type:uuid;index;uniqueIndex:user_question_unique_idx"`
	Answered   bool      `gorm:"type:boolean"`
	AnsweredAt time.Time `gorm:"type:timestamp;default:null"`

	User     User     `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Question Question `gorm:"foreignKey:QuestionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type UserQuestionFilter struct {
	UserID   *uint
	Answered *bool
}

func (f *UserQuestionFilter) Where(q *gorm.DB) *gorm.DB {
	if f.UserID != nil {
		q = q.Where("user_id = ?", f.UserID)
	}

	if f.Answered != nil {
		q = q.Where("answered = ?", *f.Answered)
	}

	return q
}

type UserQuestionRepo struct {
	db *gorm.DB
}

func NewUserQuestionRepo(db *gorm.DB) (*UserQuestionRepo, error) {
	if err := db.AutoMigrate(&UserQuestion{}); err != nil {
		return nil, fmt.Errorf("failed to migrate UserQuestion: %w", err)
	}

	return &UserQuestionRepo{db: db}, nil
}

func (r *UserQuestionRepo) Count(filter UserQuestionFilter) (int64, error) {
	query := filter.Where(r.db.Model(&UserQuestion{}))

	var count int64

	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count user questions: %w", err)
	}

	return count, nil
}

func (r *UserQuestionRepo) GetByID(id uint) (*UserQuestion, error) {
	userQuestion := &UserQuestion{}

	if err := r.WithPreloads(UserQuestionFilter{}).Where("id = ?", id).First(userQuestion).Error; err != nil {
		return nil, fmt.Errorf("failed to get user questions: %w", err)
	}

	return userQuestion, nil
}

func (r *UserQuestionRepo) Last(filter UserQuestionFilter) (*UserQuestion, error) {
	userQuestion := &UserQuestion{}

	if err := r.WithPreloads(filter).Order("id DESC").First(userQuestion).Error; err != nil {
		return nil, fmt.Errorf("failed to get last user questions: %w", err)
	}

	return userQuestion, nil
}

func (r *UserQuestionRepo) Save(userQuestion *UserQuestion) error {
	if err := r.db.Save(userQuestion).Error; err != nil {
		return fmt.Errorf("failed to create user question: %w", err)
	}

	return nil
}

func (r *UserQuestionRepo) Navigate(id uint, filter UserQuestionFilter, direction Direction) (*UserQuestion, error) {
	userQuestion := &UserQuestion{}

	if err := direction.Navigate(r.WithPreloads(filter), id).First(userQuestion).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = direction.Cycle(r.WithPreloads(filter)).First(userQuestion).Error; err != nil {
				return nil, fmt.Errorf("failed to get first user question: %w", err)
			}

			return userQuestion, nil
		}

		return nil, fmt.Errorf("failed to get user question: %w", err)
	}

	return userQuestion, nil
}

func (r *UserQuestionRepo) WithPreloads(filter UserQuestionFilter) *gorm.DB {
	return filter.Where(r.db.Model(&UserQuestion{})).Preload("User").Preload("Question")
}
