package unit

import (
	"testing"

	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/domain"
	mockports "github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/ports/mocks"
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/usecase"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_FindAll(t *testing.T){
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRepo := mockports.NewMockDatasetRepository(ctrl)

	expected := []*domain.DatasetModel{
		{ID: 1, Class: "java", Atom: "example", Snippet: "1 + 2"},
	}

	mockRepo.EXPECT().FindAll().Return(expected, nil)

	dtUsecase := usecase.NewDatasetUsecase(mockRepo)

	result, err := dtUsecase.FindAll()

	assert.NoError(t, err)
	assert.Equal(t, result, expected)
}