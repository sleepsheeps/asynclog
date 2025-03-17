package asynclog

import (
	"fmt"
	"log/slog"
	"testing"
	"time"
)

func TestLogEntry(t *testing.T) {
	entry := LogEntry{
		Level:   slog.LevelInfo,
		Message: "test",
		Attrs:   []slog.Attr{},
		Time:    time.Now(),
	}
	fmt.Println(entry.String())
}
func TestAsyncLog(t *testing.T) {
	log := NewAsyncLog(WithFile("logs", "test"), WithMaxSize(1024*10))
	defer log.Close()
	logger := slog.New(log)
	slog.SetDefault(logger)

	for i := 0; i < 10000; i++ {
		slog.Info("test", "test", "test")
	}
}
