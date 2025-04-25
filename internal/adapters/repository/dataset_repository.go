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

	if err := dr.db.Where("status_code = ?", "200 OK").Find(&dataset).Error; err != nil {
		return nil, err
	}
	return dataset, nil
}

func (dr *DatasetRepository) UpdateMarkedByAgent(agent, id int) error{
	if agent == 1{
		return dr.db.Model(&domain.Datasets{}).Where("id = ?", id).Update("marked_by_agent_one", true).Error
	}
	return dr.db.Model(&domain.Datasets{}).Where("id = ?", id).Update("marked_by_agent_two", true).Error
}