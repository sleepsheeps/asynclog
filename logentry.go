package asynclog

import (
	"log/slog"
	"time"
)

type LogEntry struct {
	Level   slog.Level
	Message string
	Time    time.Time
	Attrs   []slog.Attr
}

func (l *LogEntry) String() string {
	buf := GetBuffer()
	defer PutBuffer(buf)

	buf.WriteString(l.Time.Format("2006-01-02 15:04:05.000"))
	buf.WriteString("|")
	buf.WriteString(l.Level.String())
	buf.WriteString("|")
	buf.WriteString(l.Message)

	for _, attr := range l.Attrs {
		buf.WriteString("|")
		buf.WriteString(attr.Key)
		buf.WriteString("=")
		switch attr.Value.Kind() {
		default:
			buf.WriteString(attr.Value.String())
		}
	}
	buf.WriteString("\n")

	return buf.String()
}
