package unit

import (
	"testing"
	"time"

	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/domain"
	mockports "github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/ports/mocks"
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/usecase"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestQuestionFindAllUsecase(t *testing.T){
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRepo := mockports.NewMockQuestion(ctrl)

	expected := []*domain.Question{
		{ID: 1, Question: "", CreatedAt:time.Now(),},
	}

	mockRepo.EXPECT().FindAll().Return(expected, nil)

	dtUsecase := usecase.NewQuestionUsecase(mockRepo)

	result, err := dtUsecase.FindAll()

	assert.NoError(t, err)

	assert.Equal(t, result, expected)
	
}