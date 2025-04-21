package main

import (
	"log"

	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/adapters/config"
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/usecase"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.LoadConfig(".env")

	logit := usecase.NewLoggerUsecase(usecase.LoggerConfig)

	if err != nil{
		logit.Error("[-] Fail to load .env", zap.String("Error", err.Error()))
		log.Fatal(err)
	}

	logit.Info(cfg.DatabaseUrl, zap.String("Function", "main"))
}