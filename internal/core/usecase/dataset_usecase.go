package usecase

import (
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/domain"
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/ports/repositories"
)

type DatasetUsecase struct {
	Repo repositories.DatasetRepository
}

func NewDatasetUsecase(repo repositories.DatasetRepository) *DatasetUsecase {
	return &DatasetUsecase{
		Repo: repo,
	}
}

func (du DatasetUsecase) FindAll() ([]*domain.DatasetModel, error) {
	return du.Repo.FindAll()
}
