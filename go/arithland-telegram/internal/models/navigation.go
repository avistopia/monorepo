package models

import "gorm.io/gorm"

type Direction string

const (
	Direction_Newer Direction = "Newer"
	Direction_Older Direction = "Older"
)

func (d Direction) Navigate(q *gorm.DB, id uint) *gorm.DB {
	switch d {
	case Direction_Older:
		return q.Where("id < ?", id).Order("id DESC")
	case Direction_Newer:
		return q.Where("id > ?", id).Order("id ASC")
	default:
		return q
	}
}

func (d Direction) Cycle(q *gorm.DB) *gorm.DB {
	switch d {
	case Direction_Older:
		return q.Order("id DESC")
	case Direction_Newer:
		return q.Order("id ASC")
	default:
		return q
	}
}
