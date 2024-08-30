package log

import (
	"MusicPlayServer/common"
	"go.uber.org/zap"
)

func fileLogIniter(config *logConf) (*zap.Logger, error) {
	var (
		zapConfig zap.Config
		zaplog    *zap.Logger
		err       error
	)

	zapConfig = zap.NewProductionConfig()

	zapConfig.DisableStacktrace = true
	zapConfig.EncoderConfig.EncodeTime = epochFullTimeEncoder
	zapConfig.Level = zap.NewAtomicLevelAt(config.Level)
	if config.Env != common.ModeProd {
		zapConfig.Development = true
	}

	if len(config.LogPath) > 0 {
		logPath := config.LogPath + "/info.log"
		errLogPath := config.LogPath + "/error.log"

		zapConfig.OutputPaths = []string{logPath}
		zapConfig.ErrorOutputPaths = []string{errLogPath}
	}

	zaplog, err = zapConfig.Build()

	return zaplog, err
}
