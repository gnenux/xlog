package xlog

import (
	"os"
	"testing"
	"time"
)

func TestA(t *testing.T) {
	logger := NewLogger(os.Stdout, "", 1)
	// logger := NewLoggerFromFile("log", "", 1)
	logger.Info("hello world")
	time.Sleep(1 * time.Second)
}
