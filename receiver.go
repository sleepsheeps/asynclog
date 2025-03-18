package asynclog

import (
	"log/slog"
	"sync"
	"time"
)

type Receiver interface {
	Put(*BytesBuffer) error
	Close() error
}

type normalReceiver struct {
	WriterSync
	ch           chan *BytesBuffer
	cacheEntries []*BytesBuffer
	wg           sync.WaitGroup
}

func NewNormalReceiver(writer WriterSync) Receiver {
	r := &normalReceiver{
		WriterSync:   writer,
		ch:           make(chan *BytesBuffer, 1024),
		cacheEntries: make([]*BytesBuffer, 0, 128),
	}

	r.run()
	return r
}

func (r *normalReceiver) run() {
	if r.WriterSync == nil {
		panic("WriterSync is nil")
	}
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		ticker := time.NewTicker(time.Second * 3)
		defer ticker.Stop()
		for {
			select {
			case entry, ok := <-r.ch:
				if !ok {
					r.flush()
					return
				}
				r.cacheEntries = append(r.cacheEntries, entry)
				if entry.Level() == slog.LevelError {
					r.flush()
				}
			case <-ticker.C:
				r.flush()
			}
		}
	}()
}

func (r *normalReceiver) flush() {
	for _, entry := range r.cacheEntries {
		r.Write(entry.Bytes())
		PutBytesBuffer(entry)
	}
	r.Sync()
	if cap(r.cacheEntries) >= 256 {
		r.cacheEntries = make([]*BytesBuffer, 0, 128)
	} else {
		r.cacheEntries = r.cacheEntries[:0]
	}
}

func (r *normalReceiver) Put(entry *BytesBuffer) error {
	r.ch <- entry
	return nil
}

func (r *normalReceiver) Close() error {
	close(r.ch)
	r.wg.Wait()
	r.WriterSync.Close()
	return nil
}
