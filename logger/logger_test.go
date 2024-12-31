package logger

import (
	"testing"
)

func TestXxx(t *testing.T) {
	//	Init("./", slog.LevelDebug)
	log := GetLog("testlog", "", "./")
	log.Info("test init", "basedir", "./")
	log.Error("test init", "basedir", "./")

	Close()
}
