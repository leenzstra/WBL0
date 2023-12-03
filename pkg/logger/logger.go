package logger

import "go.uber.org/zap"

func New(logFile string, debug bool) (*zap.Logger, error) {
	var config zap.Config
	if debug {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	config.OutputPaths = []string{"stdout", logFile}
	logger, err := config.Build(zap.AddCaller())
	if err != nil {
		return nil, err
	}

	return logger, nil
}
