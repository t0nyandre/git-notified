package logger

import (
	"log"

	"go.uber.org/zap"
)

func NewLogger() *zap.SugaredLogger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Could not create logger: %v", err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	return sugar
}
