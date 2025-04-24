package domain

import "time"

type Question struct {
	ID        uint      `gorm:"primaryKey"`
	Question  string    `gorm:"column:definition"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (Question) TableName() string {
	return "questions"
}
