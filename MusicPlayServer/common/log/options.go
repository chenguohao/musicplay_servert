package log

import (
	"MusicPlayServer/common"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
)

type logConf struct {
	Env         string
	ProcessName string
	Level       zapcore.Level
	LogPath     string
	WithPid     bool
}

var DefaultLogOptions = logConf{
	Env:         common.ModeDebug,
	ProcessName: logName(),
	Level:       zapcore.DebugLevel,
}

type LogOptions func(conf *logConf)

func logName() string {
	return path.Base(os.Args[0])
}
