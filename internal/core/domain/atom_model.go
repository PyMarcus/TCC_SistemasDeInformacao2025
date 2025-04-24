package domain

import "time"

type Atom struct {
	ID                   uint   `gorm:"primaryKey"`
	QuestionID           int    `gorm:"column:question_id"`
	Question             string `gorm:"type:text"`
	AgentOneAnswer       string `gorm:"type:text"`
	AgentTwoAnswer       string `gorm:"type:text"`
	DatasetID            int    `gorm:"column:dataset_id"`
	AtomSearched         string `gorm:"column:atom_searched"`
	AtomFindedByAgentOne string `gorm:"column:atom_finded_by_agent_one"`
	AtomFindedByAgentTwo string `gorm:"column:atom_finded_by_agent_two"`
	AgentOneIsCorrect    bool   `gorm:"column:agent_one_is_correct"`
	AgentTwoIsCorrect    bool   `gorm:"column:agent_two_is_correct"`
	Failed               bool
	ErrorID              int       `gorm:"column:error_id"`
	CreatedAt            time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt            time.Time `gorm:"column:updated_at"`
}

func (Atom) TableName() string {
	return "atoms"
}
