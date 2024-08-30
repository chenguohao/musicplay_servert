package log

import (
	"MusicPlayServer/common"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

func LogInit(config *logConf) (*zap.Logger, error) {
	var (
		log *zap.Logger
		err error
	)

	switch config.Env {
	case common.ModeTest, common.ModeProd:
		log, err = fileLogIniter(config)
	default:
		log, err = consoleLogIniter(config)
	}

	if err != nil {
		return log, err
	}

	if config.WithPid {
		log = log.With(zap.Int("pid", os.Getpid()))
	}

	return log, nil
}

func epochFullTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}
