package log

import (
	"go.uber.org/zap"
)

func consoleLogIniter(config *logConf) (*zap.Logger, error) {
	var (
		zapConfig zap.Config
		zaplog    *zap.Logger
		err       error
	)

	zapConfig = zap.NewDevelopmentConfig()

	zapConfig.DisableStacktrace = true
	zapConfig.EncoderConfig.EncodeTime = epochFullTimeEncoder

	zaplog, err = zapConfig.Build()

	return zaplog, err
}
