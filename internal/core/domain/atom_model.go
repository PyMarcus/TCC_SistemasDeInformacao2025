package domain

import "time"

type Atom struct {
	ID                   uint   `gorm:"primaryKey"`
	QuestionID           int    `gorm:"column:question_id"`
	Question             string `gorm:"type:text"`
	Answer       		 string `gorm:"type:text"`
	DatasetID            int    `gorm:"column:datasets_id"`
	AtomSearched         string `gorm:"column:atom_searched"`
	AtomFinded			 string `gorm:"column:atom_finded"`
	IsCorrect    		 bool   `gorm:"column:is_correct"`
	Failed               bool   `gorm:"column:failed"`
	ErrorID              int       `gorm:"column:error_id"`
	CreatedAt            time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt            time.Time `gorm:"column:updated_at"`
}

func (Atom) TableName() string {
	return "atoms"
}
