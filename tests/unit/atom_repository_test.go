package unit

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/adapters/repository"
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestAtomCreate(t *testing.T){
	db, mock, cleanup := setupMockDB(t)

	defer cleanup()

	repo := repository.NewAtomRepository(db)

	atomModel := &domain.Atom{
		QuestionID:            0,
		Question:              "TEST",
		Answer:        "",
		DatasetID:             0,
		AtomSearched:          "",
		AtomFinded:  "",
		IsCorrect:     false,
		Failed:                false,
		ErrorID:               0,
		UpdatedAt:             time.Now(),
		CreatedAt:             time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "atoms"`).WithArgs(
        atomModel.QuestionID,
        atomModel.Question,
        atomModel.Answer,
        atomModel.DatasetID,
        atomModel.AtomSearched,
        atomModel.AtomFinded,
        atomModel.IsCorrect,
        atomModel.Failed,
        atomModel.ErrorID,
        atomModel.UpdatedAt,
        atomModel.CreatedAt,
    ).
	WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	id, err := repo.Create(atomModel)
	
	assert.NoError(t, err)

	assert.Equal(t, int(id), 1)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

}