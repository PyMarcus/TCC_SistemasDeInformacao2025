package usecase

import (
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/domain"
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/ports/repositories"
)

type AtomUsecase struct{
	Repo repositories.Atom
}

func NewAtomUsecase(repo repositories.Atom) *AtomUsecase{
	return &AtomUsecase{
		Repo: repo,
	}
}

func (au *AtomUsecase) Create(atom *domain.Atom) (uint, error){
	return au.Repo.Create(atom)
}