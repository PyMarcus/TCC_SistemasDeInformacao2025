package repository

import (
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/domain"
	"gorm.io/gorm"
)

type AtomRepository struct{
	db *gorm.DB
}

func NewAtomRepository(db *gorm.DB) *AtomRepository{
	return &AtomRepository{
		db: db,
	}
}

func (a *AtomRepository) Create(atom *domain.Atom) (uint, error){
	if err := a.db.Save(atom).Error; err != nil{
		return 0, err 
	}

	return atom.ID, nil
}