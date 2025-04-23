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

	expected := []*domain.Datasets{
		{ID: 1, Class: "java", Atom: "example", Snippet: "1 + 2"},
	}

	mockRepo.EXPECT().FindAll().Return(expected, nil)

	dtUsecase := usecase.NewDatasetUsecase(mockRepo)

	result, err := dtUsecase.FindAll()

	assert.NoError(t, err)
	assert.Equal(t, result, expected)
}

func Test_UpdateMarkedByAgent(t *testing.T){
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRepo := mockports.NewMockDatasetRepository(ctrl)

	mockRepo.EXPECT().UpdateMarkedByAgent(1, 1).Times(1)

	dtUsecase := usecase.NewDatasetUsecase(mockRepo)

	err := dtUsecase.UpdateMarkedByAgent(1, 1)

	assert.NoError(t, err)
}