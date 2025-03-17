package asynclog

import (
	"context"
	"log/slog"
)

type AsyncLogOpt struct {
	FilePath string
	Filename string
	Level    slog.Level
	isFile   bool
	MaxSize  int
}

type AsyncLogOption func(*AsyncLogOpt)

type AsyncLog struct {
	receiver Receiver
}

func WithFile(filePath, filename string) AsyncLogOption {
	return func(log *AsyncLogOpt) {
		log.FilePath = filePath
		log.Filename = filename
		log.isFile = true
	}
}

func WithLevel(level slog.Level) AsyncLogOption {
	return func(log *AsyncLogOpt) {
		log.Level = level
	}
}

func WithMaxSize(maxSize int) AsyncLogOption {
	return func(log *AsyncLogOpt) {
		log.MaxSize = maxSize
	}
}

func NewAsyncLog(opts ...AsyncLogOption) *AsyncLog {
	asyncLogOpt := &AsyncLogOpt{
		Level:  slog.LevelInfo,
		isFile: false,
	}
	for _, opt := range opts {
		opt(asyncLogOpt)
	}

	var writer WriterSync

	if asyncLogOpt.isFile {
		writer = NewFileWriter(asyncLogOpt.FilePath, asyncLogOpt.Filename, asyncLogOpt.MaxSize)
	} else {
		writer = NewStdWriter()
	}

	return &AsyncLog{
		receiver: NewNormalReceiver(writer, asyncLogOpt.Level),
	}
}

func (a *AsyncLog) Close() error {
	return a.receiver.Close()
}

// Enabled implements slog.Handler interface
func (a *AsyncLog) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

// Handle implements slog.Handler interface
func (a *AsyncLog) Handle(ctx context.Context, r slog.Record) error {
	// 创建属性切片来存储所有属性
	attrs := make([]slog.Attr, 0, r.NumAttrs())

	// 添加记录中的属性
	r.Attrs(func(attr slog.Attr) bool {
		attrs = append(attrs, attr)
		return true
	})

	// 创建日志事件
	event := LogEntry{
		Time:    r.Time,
		Level:   r.Level,
		Message: r.Message,
		Attrs:   attrs,
	}

	// 发送到接收器
	return a.receiver.Put(event)
}

// WithAttrs implements slog.Handler interface
func (a *AsyncLog) WithAttrs(attrs []slog.Attr) slog.Handler {
	return a
}

// WithGroup implements slog.Handler interface
func (a *AsyncLog) WithGroup(name string) slog.Handler {
	return a
}
