package unit

import (
	"testing"

	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/domain"
	mockports "github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/ports/mocks"
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/usecase"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_Create(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mockports.NewMockError(ctrl)

	erroR := domain.Error{Definition: "ErroR"}

	mockRepo.EXPECT().Create(&erroR).Times(1)

	_, err := usecase.NewErrorUsecase(mockRepo).Create(&erroR)

	assert.NoError(t, err)
	
}
