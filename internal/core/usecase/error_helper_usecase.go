package usecase

import (
	"fmt"

	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/constants"
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/adapters/http/dto"
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/domain"
	"go.uber.org/zap"
)

func HandleAgentError(
	logger *LoggerUsecase,
	errorUsecase *ErrorUsecase,
	err error,
	apiName string,
	url string,
	status string,
	agentChan chan<- dto.ClientResponseDTO,
) int {
	logger.Error(
		fmt.Sprintf("[-] Fail to request in %s", apiName),
		zap.String(constants.URL, url),
		zap.String(constants.ERROR_DESCRIPTION, err.Error()),
		zap.String(constants.STATUS_CODE_STR, status),
	)

	id, dbErr := errorUsecase.Create(&domain.Error{Definition: err.Error()})
	if dbErr != nil {
		logger.Error(
			fmt.Sprintf("[-] Fail to insert error into database in %s", apiName),
			zap.String(constants.ERROR_DESCRIPTION, dbErr.Error()),
		)
	}
	agentChan <- dto.ClientResponseDTO{
		Message: err.Error(),
		Api: apiName,
	}

	return int(id)
}
