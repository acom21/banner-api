package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

// Init initializes the global logger
func Init() error {
	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}

	Log = logger
	return nil
}
