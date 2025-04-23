package repository

import (
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/domain"
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/ports/repositories"
	"gorm.io/gorm"
)

type DatasetRepository struct {
	db *gorm.DB
}

func NewDatasetRepository(db *gorm.DB) repositories.DatasetRepository {
	return &DatasetRepository{
		db: db,
	}
}

func (dr *DatasetRepository) FindAll() ([]*domain.Datasets, error) {
	var dataset []*domain.Datasets

	if err := dr.db.Find(&dataset).Error; err != nil {
		return nil, err
	}
	return dataset, nil
}
