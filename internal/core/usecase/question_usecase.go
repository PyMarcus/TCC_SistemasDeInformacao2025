package usecase

import (
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/domain"
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/ports/repositories"
)

type QuestionUsecase struct {
	Repo repositories.Question
}

func NewQuestionUsecase(repo repositories.Question) *QuestionUsecase {
	return &QuestionUsecase{
		Repo: repo,
	}
}

func (du QuestionUsecase) FindAll() ([]*domain.Question, error) {
	return du.Repo.FindAll()
}
