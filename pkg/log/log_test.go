package log

import (
	"testing"

	"github.com/superOTAKU/onlineTexasPoker/pkg/config"
)

func TestLog(t *testing.T) {
	InitLogger(&config.LogConfig{
		FilePath: "./test.log",
	})
	Logger().Info("logger test")
}
