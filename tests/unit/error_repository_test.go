package unit

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/adapters/repository"
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T){
	db, mock, cleanup := setupMockDB(t)

	defer cleanup()

	repo := repository.NewErrorRepository(db)

	errorModel := &domain.Error{
		Definition: "TEST",
		CreatedAt: time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "errors"`).WithArgs(errorModel.Definition, errorModel.CreatedAt).
	WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	id, err := repo.Create(errorModel)
	
	assert.NoError(t, err)

	assert.Equal(t, int(id), 1)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)


}