package repository

import (
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/domain"
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/ports/repositories"
	"gorm.io/gorm"
)

type ErrorRepository struct{
	db *gorm.DB
}

func NewErrorRepository(db *gorm.DB) repositories.Error{
	return &ErrorRepository{
		db: db,
	}
}

func (er *ErrorRepository) Create(errorModel *domain.Error) (uint, error){
	err := er.db.Save(errorModel).Error
	if err != nil{
		return 0, err 
	}

	return errorModel.ID, nil
}