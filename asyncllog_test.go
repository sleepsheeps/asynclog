package asynclog

import (
	"log/slog"
	"testing"
)

func TestAsyncLog(t *testing.T) {
	log := NewAsyncLog(WithFile("logs", "test"), WithMaxSize(1024*10))
	defer log.Close()

	for i := 0; i < 512; i++ {
		slog.Info("test", "test", "test")
	}
}
