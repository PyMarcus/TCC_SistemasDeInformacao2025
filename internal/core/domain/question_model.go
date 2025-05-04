package domain

import "time"

type Question struct {
	ID        uint      `gorm:"primaryKey"`
	Question  string    `gorm:"column:question"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (Question) TableName() string {
	return "questions"
}
