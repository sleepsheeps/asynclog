package asynclog

import (
	"log/slog"
	"sync"
	"time"
)

type Receiver interface {
	Put(LogEntry) error
	Close() error
}

type normalReceiver struct {
	WriterSync
	level        slog.Level
	entries      chan LogEntry
	cacheEntries []LogEntry
	wg           sync.WaitGroup
}

func NewNormalReceiver(writer WriterSync, level slog.Level) Receiver {
	r := &normalReceiver{
		WriterSync:   writer,
		entries:      make(chan LogEntry, 1024),
		cacheEntries: make([]LogEntry, 0, 128),
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
			case entry, ok := <-r.entries:
				if !ok {
					r.flush()
					return
				}
				r.cacheEntries = append(r.cacheEntries, entry)
				if entry.Level == slog.LevelError {
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
		if entry.Level < r.level {
			continue
		}
		r.Write([]byte(entry.String()))
	}
	r.Sync()
	if cap(r.cacheEntries) >= 256 {
		r.cacheEntries = make([]LogEntry, 0, 128)
	} else {
		r.cacheEntries = r.cacheEntries[:0]
	}
}

func (r *normalReceiver) Put(entry LogEntry) error {
	r.entries <- entry
	return nil
}

func (r *normalReceiver) Close() error {
	close(r.entries)
	r.wg.Wait()
	r.WriterSync.Close()
	return nil
}
