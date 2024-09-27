package log

import (
	"MusicPlayServer/common"
	"testing"
)

var cfg = common.Config{
	Mode: common.ModeDebug,
}

func TestLog(t *testing.T) {
	cfg.Mode = common.ModeTest
	cfg.Log.Output = "/data/logs"

	InitLog(&cfg)

	Debug("Debug log")
	Info("Info Log")
}
