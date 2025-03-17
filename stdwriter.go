package asynclog

import (
	"io"
	"os"
)

type WriterSync interface {
	Write(p []byte) (n int, err error)
	Sync() error
	Close() error
}

type stdWriter struct {
	io.WriteCloser
}

func NewStdWriter() WriterSync {
	return &stdWriter{
		WriteCloser: os.Stdout,
	}
}

func (w *stdWriter) Write(p []byte) (n int, err error) {
	return w.WriteCloser.Write(p)
}

func (w *stdWriter) Sync() error {
	return nil
}

func (w *stdWriter) Close() error {
	return nil
}
