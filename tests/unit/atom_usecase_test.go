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

func TestAtomUsecaseCreate(t *testing.T){
	ctrl := gomock.NewController(t)

	mockRepo := mockports.NewMockAtom(ctrl)

	atom := domain.Atom{
		QuestionID:            0,
		Question:              "TEST",
		AgentOneAnswer:        "",
		AgentTwoAnswer:        "",
		DatasetID:             0,
		AtomSearched:          "",
		AtomFindedByAgentOne:  "",
		AtomFindedByAgentTwo:  "",
		AgentOneIsCorrect:     false,
		AgentTwoIsCorrect:     false,
		Failed:                false,
		ErrorID:               0,
		UpdatedAt:             time.Now(),
		CreatedAt:             time.Now(),
	}

	mockRepo.EXPECT().Create(&atom).Times(1)

	_, err := usecase.NewAtomUsecase(mockRepo).Create(&atom)

	assert.NoError(t, err)
}