package repository

import (
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/domain"
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/ports/repositories"
	"gorm.io/gorm"
)

type QuestionRepository struct{
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) repositories.Question{
	return &QuestionRepository{
		db: db,
	}
}

func (qr QuestionRepository) FindAll() ([]*domain.Question, error){
	var questions []*domain.Question

	if err := qr.db.Find(&questions).Error; err != nil{
		return nil, err 
	}

	return questions, nil
}