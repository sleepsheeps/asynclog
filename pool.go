package asynclog

import (
	"bytes"
	"sync"
)

var (
	bufferPool = NewBufferPool()
)

func GetBuffer() *bytes.Buffer {
	return bufferPool.Get()
}

func PutBuffer(b *bytes.Buffer) {
	bufferPool.Put(b)
}

type BufferPool struct {
	sync.Pool
}

func NewBufferPool() *BufferPool {
	return &BufferPool{
		Pool: sync.Pool{
			New: func() any { return bytes.NewBuffer(nil) },
		},
	}
}

func (p *BufferPool) Get() *bytes.Buffer {
	return p.Pool.Get().(*bytes.Buffer)
}

func (p *BufferPool) Put(b *bytes.Buffer) {
	b.Reset()
	p.Pool.Put(b)
}
