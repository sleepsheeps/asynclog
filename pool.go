package asynclog

import (
	"log/slog"
	"sync"
)

var (
	bytesBufferPool = sync.Pool{
		New: func() any { return NewBytesBuffer() },
	}
)

func GetBytesBuffer() *BytesBuffer {
	b := bytesBufferPool.Get().(*BytesBuffer)
	b.Reset()
	return b
}

func PutBytesBuffer(b *BytesBuffer) {
	bytesBufferPool.Put(b)
}

type BytesBuffer struct {
	bs    []byte
	level slog.Level
}

func NewBytesBuffer() *BytesBuffer {
	return &BytesBuffer{bs: make([]byte, 0, 1024)}
}

func (b *BytesBuffer) Write(p []byte) (n int, err error) {
	b.bs = append(b.bs, p...)
	return len(p), nil
}

func (b *BytesBuffer) WriteString(s string) (n int, err error) {
	b.bs = append(b.bs, s...)
	return len(s), nil
}

func (b *BytesBuffer) String() string {
	return string(b.bs)
}

func (b *BytesBuffer) Bytes() []byte {
	return b.bs
}

func (b *BytesBuffer) Reset() {
	b.bs = b.bs[:0]
}

func (b *BytesBuffer) Level() slog.Level {
	return b.level
}

func (b *BytesBuffer) SetLevel(level slog.Level) {
	b.level = level
}
