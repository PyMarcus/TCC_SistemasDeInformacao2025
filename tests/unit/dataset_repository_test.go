package unit

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/adapters/repository"
	"github.com/stretchr/testify/assert"
)

func TestFindAll(t *testing.T){
	db, mock, cleanup := setupMockDB(t)

	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "datasets"`)).
	WillReturnRows(
		sqlmock.NewRows(
			[]string{
				"id", "class", "atom", "snippet", "line", "github_link", "status_code",
				"marked_by_agent_one", "marked_by_agent_two",
			}).AddRow(
				1, "security", "atom-1", "snippet-1", "line-1", "http://github.com", "200 OK", true, false,
			),
	)

	repo := repository.NewDatasetRepository(db)

	result, err := repo.FindAll()

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, "security", result[0].Class)
}