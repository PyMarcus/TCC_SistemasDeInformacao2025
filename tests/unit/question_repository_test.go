package unit

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/adapters/repository"
	"github.com/stretchr/testify/assert"
)

func TestQuestionFindAll(t *testing.T){
	db, mock, cleanup := setupMockDB(t)

	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "questions"`)).
	WillReturnRows(
		sqlmock.NewRows(
			[]string{
				"id", "question", "create_at",
			}).AddRow(
				1, "", time.Now(),
			),
	)

	repo := repository.NewQuestionRepository(db)

	result, err := repo.FindAll()

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, 1, int(result[0].ID))
}
