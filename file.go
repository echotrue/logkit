package logkit

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileTarget struct {
	FilePath string

	lock      sync.Mutex
	fd        *os.File
	errWriter io.Writer
	writer    *bufio.Writer
	tick      *time.Ticker
	close     chan bool
}

func NewFileTarget() *FileTarget {
	return &FileTarget{
		FilePath: os.TempDir(),
		close:    make(chan bool, 0),
	}
}

func (t *FileTarget) Open(w io.Writer) error {
	path, err := buildPath(t.FilePath)
	if err != nil {
		return err
	}
	t.FilePath = filepath.Join(path, time.Now().Format("20060102")+".log")

	fd, err := os.OpenFile(t.FilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("FileTarget was unable to create a log file: %v\n", err)
	}
	t.fd = fd
	t.writer = bufio.NewWriter(t.fd)
	t.errWriter = w

	t.tick = time.NewTicker(2 * time.Second)

	go t.clearBuffer()

	return nil
}

func (t *FileTarget) Process(entry *Entry) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if entry == nil {
		t.close <- true
		return
	}
	if t.fd != nil && t.writer != nil {
		content := entry.String() + "\n"
		_, _ = t.writer.Write([]byte(content))
	}
}

// clearBuffer clear buffer every second
func (t *FileTarget) clearBuffer() {
	for {
		select {
		case <-t.tick.C:
			_ = t.writer.Flush()
		}
	}
}

func (t *FileTarget) Close() {
	<-t.close

	// Stop turns off a ticker
	t.tick.Stop()
	// Write any buffered data to the underlying io.Writer before close
	_ = t.writer.Flush()
}

func buildPath(p string) (path string, err error) {
	path, err = filepath.Abs(p)
	if err != nil {
		return
	}
	if _, err := os.Stat(p); err != nil {
		if os.IsNotExist(err) {
			_ = os.MkdirAll(p, 755)
		}
	}
	return
}
