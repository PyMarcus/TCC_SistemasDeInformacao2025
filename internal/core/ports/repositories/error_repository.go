package repositories

import "github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/domain"

type Error interface{
	Create(error_model *domain.Error) (uint, error)
}