package repositories

import (
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/domain"
)

type DatasetRepository interface{
	FindAll() ([]*domain.Datasets, error)
}