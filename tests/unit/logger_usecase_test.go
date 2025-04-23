package unit

import (
	"testing"

	mockports "github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/ports/mocks"
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/usecase"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)


func TestInfo(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockLogger := mockports.NewMockLogger(ctrl)

	mockLogger.EXPECT().Info("[!]Executing log test", zap.String("user", "marcus")).Times(1)

	usecase := usecase.NewLoggerUsecase(mockLogger)
	usecase.Info("[!]Executing log test", zap.String("user", "marcus"))

}

func TestError(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockLogger := mockports.NewMockLogger(ctrl)

	mockLogger.EXPECT().Error("[!]Executing error log test", zap.String("user", "marcus"))
	usecase := usecase.NewLoggerUsecase(mockLogger)
	usecase.Error("[!]Executing error log test", zap.String("user", "marcus"))
}
