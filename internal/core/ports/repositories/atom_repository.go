package repositories

import "github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/domain"

type Atom interface{
	Create(atom *domain.Atom) (uint, error)
}