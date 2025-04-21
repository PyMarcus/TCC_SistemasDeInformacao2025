package usecase

import (
	"log"
	"os"
	"path/filepath"

	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/constants"
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/ports/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LoggerUsecase struct{
	jsonLogger logger.Logger
}

var LoggerConfig *zap.Logger

func init(){
	writerSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join(constants.LOG_DIR, constants.LOG_NAME),
		MaxSize:    300, // megabytes
		MaxBackups: 5,
		MaxAge:     30,   // days
		Compress:   true, // gzip 
	})

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		writerSyncer,
		zap.InfoLevel,
	)

	LoggerConfig = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

// NewLoggerUsecase: receives the global var LoggerConfig *zap.Logger
func NewLoggerUsecase(logger logger.Logger) *LoggerUsecase{

	if err := os.MkdirAll(constants.LOG_DIR, os.ModePerm); err != nil{
		log.Fatal("[-] Fail to create log dir: " + err.Error())
		return nil
	}

	return &LoggerUsecase{
		jsonLogger: logger,
	}
}

func (lu LoggerUsecase) Info(message string, fields ...zap.Field){
	lu.jsonLogger.Info(message, fields...)
}

func (lu LoggerUsecase) Error(message string, fields ...zap.Field){
	lu.jsonLogger.Error(message, fields...)
}