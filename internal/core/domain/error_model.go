package domain

import "time"

type Error struct {
	ID         uint      `gorm:"primaryKey"`
	Definition string    `gorm:"column:definition"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}

func (Error) TableName() string{
	return "errors"
}