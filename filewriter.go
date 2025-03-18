package asynclog

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"
)

const (
	DefaultFileSize = 512 * 1024 * 1024
)

type FileWriter struct {
	f          *os.File
	filePath   string
	filename   string
	dir        string
	writer     *bufio.Writer
	lastRotate time.Time
	maxSize    int
	fileSize   int64
}

func NewFileWriter(filePath, filename string, maxSize int) WriterSync {

	fw := &FileWriter{
		filePath: filePath,
		filename: filename,
		maxSize:  maxSize,
	}

	if maxSize <= 0 {
		fw.maxSize = DefaultFileSize
	}

	fw.lastRotate = time.Now()

	if err := fw.InitFile(); err != nil {
		panic(err)
	}

	return fw
}

func (w *FileWriter) Write(p []byte) (n int, err error) {
	if err := w.Rotate(); err != nil {
		return 0, err
	}
	w.fileSize += int64(len(p))
	return w.writer.Write(p)
}

func (w *FileWriter) Sync() error {
	if w.f == nil {
		return nil
	}
	if err := w.writer.Flush(); err != nil {
		return err
	}
	return w.f.Sync()
}

func (w *FileWriter) Close() error {
	if w.f == nil {
		return nil
	}
	return w.f.Close()
}

func (w *FileWriter) Rotate() error {
	// 检测日期，如果不是同年同月同日，则创建新文件
	// 是否是同一天
	now := time.Now()
	if !checkSameDay(now, w.lastRotate) {
		w.lastRotate = now
		return w.InitFile()
	}

	// 大小分割
	if !w.checkSize() {
		return nil
	}

	w.Sync()

	// Close the current file before renaming
	if err := w.f.Close(); err != nil {
		return err
	}

	fullPath := filepath.Join(w.dir, fmt.Sprintf("%s.log", w.filename))

	// Check if the file exists before renaming
	if _, err := os.Stat(fullPath); err == nil {
		newName := filepath.Join(w.dir, fmt.Sprintf("%s.%s.log", w.filename, now.Format("150405.000")))
		if err := os.Rename(fullPath, newName); err != nil {
			return err
		}
	}

	// 创建新文件
	file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	w.f = file
	w.fileSize = 0
	w.writer = bufio.NewWriterSize(w.f, 256*1024)
	return nil
}

func (w *FileWriter) InitFile() error {
	now := time.Now()

	// 新目录
	dir := path.Join(w.filePath, now.Format("2006-01-02"))

	// 创建目录
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	w.dir = dir

	fileName := fmt.Sprintf("%s.log", w.filename)

	fullPath := filepath.Join(dir, fileName)

	// 创建文件
	file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	w.f = file
	w.fileSize = stat.Size()
	w.writer = bufio.NewWriterSize(w.f, 256*1024) // 256KB

	return nil
}

func (w *FileWriter) checkSize() bool {
	return w.fileSize >= int64(w.maxSize)
}

func checkSameDay(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}
