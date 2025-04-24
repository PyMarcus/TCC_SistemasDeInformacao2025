package repositories

import "github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/domain"

type Question interface{
	FindAll() ([]*domain.Question, error)
}