package models

import (
	"fmt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	TelegramUserID int64   `gorm:"type:bigint;unique"`
	DisplayName    string  `gorm:"type:varchar(255)"`
	Balance        float64 `gorm:"type:decimal(10,2)"`
	State          State   `gorm:"type:json"`
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) (*UserRepo, error) {
	if err := db.AutoMigrate(&User{}); err != nil {
		return nil, fmt.Errorf("failed to migrate User: %w", err)
	}

	return &UserRepo{db: db}, nil
}

func (r *UserRepo) GetOrCreateUserByTelegramUserID(telegramUserID int64) (*User, error) {
	user := &User{TelegramUserID: telegramUserID, Balance: 900}
	if err := r.db.Where(User{TelegramUserID: telegramUserID}).FirstOrCreate(user).Error; err != nil {
		return nil, fmt.Errorf("failed to get or create user: %w", err)
	}

	return user, nil
}

func (r *UserRepo) GetByID(id uint) (*User, error) {
	user := &User{}

	err := r.db.Where("id = ?", id).First(user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *UserRepo) Save(user *User) error {
	if err := r.db.Save(user).Error; err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}
