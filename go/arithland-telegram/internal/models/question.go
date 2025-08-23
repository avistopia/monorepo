package models

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type QuestionLevel string

func (v QuestionLevel) Render() string {
	switch v {
	case QuestionLevel_Easy:
		return "ساده"
	case QuestionLevel_Medium:
		return "متوسط"
	case QuestionLevel_Hard:
		return "سخت"
	default:
		return "نامشخص"
	}
}

func (v QuestionLevel) Emoji() string {
	switch v {
	case QuestionLevel_Easy:
		return "🟢"
	case QuestionLevel_Medium:
		return "🟡"
	case QuestionLevel_Hard:
		return "🔴"
	default:
		return "⚫"
	}
}

const (
	QuestionLevel_Easy   QuestionLevel = "Easy"
	QuestionLevel_Medium QuestionLevel = "Medium"
	QuestionLevel_Hard   QuestionLevel = "Hard"
)

var (
	QuestionLevels = []QuestionLevel{QuestionLevel_Easy, QuestionLevel_Medium, QuestionLevel_Hard}

	QuestionPrices = map[QuestionLevel]int64{
		QuestionLevel_Easy:   1,
		QuestionLevel_Medium: 3,
		QuestionLevel_Hard:   5,
	}
)

type Question struct {
	gorm.Model

	Text     string        `gorm:"type:text"`
	PhotoID  string        `gorm:"type:text"`
	Answer   string        `gorm:"type:text"`
	Level    QuestionLevel `gorm:"type:varchar(20)"`
	Inactive bool          `gorm:"type:bool"`
}

type QuestionFilter struct {
	Level             *QuestionLevel
	Inactive          *bool
	UserIDDoesNotHave *uint
}

func (f *QuestionFilter) Where(q *gorm.DB) *gorm.DB {
	if f.Inactive != nil {
		q = q.Where("inactive = ?", f.Inactive)
	}

	if f.Level != nil {
		q = q.Where("level = ?", string(*f.Level))
	}

	if f.UserIDDoesNotHave != nil {
		q = q.
			Joins(
				"LEFT JOIN user_questions uq ON uq.question_id = questions.id AND uq.user_id = ?", *f.UserIDDoesNotHave,
			).
			Where("uq.id IS NULL")
	}

	return q
}

type QuestionRepo struct {
	db *gorm.DB
}

func NewQuestionRepo(db *gorm.DB) (*QuestionRepo, error) {
	if err := db.AutoMigrate(&Question{}); err != nil {
		return nil, fmt.Errorf("failed to migrate Question: %w", err)
	}

	return &QuestionRepo{db: db}, nil
}

func (r *QuestionRepo) GetByID(id uint) (*Question, error) {
	question := &Question{}

	err := r.db.Where("id = ?", id).First(question).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get question: %w", err)
	}

	return question, nil
}

func (r *QuestionRepo) GetByIDOrNext(id uint) (*Question, error) {
	question := &Question{}

	err := r.db.Where("id >= ?", id).Order("id ASC").First(question).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = r.db.Order("id ASC").First(question).Error; err != nil {
				return nil, fmt.Errorf("failed to get first question: %w", err)
			}

			return question, nil
		}

		return nil, fmt.Errorf("failed to get question: %w", err)
	}

	return question, nil
}

func (r *QuestionRepo) Navigate(id uint, direction Direction) (*Question, error) {
	question := &Question{}

	if err := direction.Navigate(r.db.Model(&Question{}), id).First(question).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = direction.Cycle(r.db.Model(&Question{})).First(question).Error; err != nil {
				return nil, fmt.Errorf("failed to get first question: %w", err)
			}

			return question, nil
		}

		return nil, fmt.Errorf("failed to get question: %w", err)
	}

	return question, nil
}

func (r *QuestionRepo) GetRandom(filter QuestionFilter) (*Question, error) {
	question := &Question{}

	if err := filter.Where(r.db).Order("RANDOM()").First(question).Error; err != nil {
		return nil, fmt.Errorf("failed to get question: %w", err)
	}

	return question, nil
}

func (r *QuestionRepo) Save(question *Question) error {
	if err := r.db.Save(question).Error; err != nil {
		return fmt.Errorf("failed to create question: %w", err)
	}

	return nil
}

func (r *QuestionRepo) Count() (int64, error) {
	var count int64
	if err := r.db.Model(&Question{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count questions: %w", err)
	}

	return count, nil
}
