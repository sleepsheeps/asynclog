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
	level    slog.Level
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

	log := &AsyncLog{
		receiver: NewNormalReceiver(writer),
		level:    asyncLogOpt.Level,
	}

	logger := slog.New(log)
	slog.SetDefault(logger)

	return log
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
	if r.Level < a.level {
		return nil
	}

	buf := GetBytesBuffer()
	buf.SetLevel(r.Level)
	buf.WriteString(r.Time.Format(TimeFormat))
	buf.WriteString(Split)
	buf.WriteString(r.Level.String())
	buf.WriteString(Split)
	buf.WriteString(r.Message)

	r.Attrs(func(attr slog.Attr) bool {
		buf.WriteString(Split)
		buf.WriteString(attr.Key)
		buf.WriteString(ValueSplit)
		buf.WriteString(attr.Value.String())
		return true
	})
	buf.WriteString(NewLine)

	// 发送到接收器
	return a.receiver.Put(buf)
}

// WithAttrs implements slog.Handler interface
func (a *AsyncLog) WithAttrs(attrs []slog.Attr) slog.Handler {
	return a
}

// WithGroup implements slog.Handler interface
func (a *AsyncLog) WithGroup(name string) slog.Handler {
	return a
}
