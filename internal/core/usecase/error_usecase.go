package usecase

import (
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/domain"
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/ports/repositories"
)

type ErrorUsecase struct {
	repo repositories.Error
}

func NewErrorUsecase(repo repositories.Error) *ErrorUsecase {
	return &ErrorUsecase{
		repo: repo,
	}
}

func (eu *ErrorUsecase) Create(errorModel *domain.Error) (uint, error) {
	return eu.repo.Create(errorModel)
}
