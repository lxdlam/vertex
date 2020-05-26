package log

import (
	"bufio"
	"errors"
	"io"
	"os"
	"sync"
)

// ErrTargetFileInvalid will be raise when give PersistentFile.SaveTo a nil file.
var ErrTargetFileInvalid = errors.New("persistent: target file invalid")

// PersistentFile is a simple buffered file for write,
// all operations can be done safely across goroutines.
// It warps all bufio.Writer methods with a RWMutex.
type PersistentFile struct {
	mutex  sync.RWMutex
	file   *os.File
	writer *bufio.Writer
}

// NewPersistentFile will returns a new PersistentFile instance warps the given file.
func NewPersistentFile(file *os.File) *PersistentFile {
	return &PersistentFile{
		mutex:  sync.RWMutex{},
		file:   file,
		writer: bufio.NewWriter(file),
	}
}

func (p *PersistentFile) Available() int {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return p.writer.Available()
}

func (p *PersistentFile) Buffered() int {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return p.writer.Buffered()
}

func (p *PersistentFile) Flush() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.writer.Flush()
}

// ReadFrom do what?
// Not to proxy it here.

// Reset unused here. Not to proxy as well.

func (p *PersistentFile) Size() int {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return p.writer.Size()
}

func (p *PersistentFile) Write(b []byte) (int, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.writer.Write(b)
}

func (p *PersistentFile) WriteByte(c byte) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.writer.WriteByte(c)
}

func (p *PersistentFile) WriteString(s string) (int, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.writer.WriteString(s)
}


// SaveTo will copy the current log file into another file.
// A Flush() will be performed in advance to correctly save the current file.
func (p *PersistentFile) SaveTo(file *os.File) error {
	if file == nil {
		return ErrTargetFileInvalid
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	// force save file
	if err := p.writer.Flush(); err != nil {
		return err
	}

	r := bufio.NewReader(p.file)
	w := bufio.NewWriter(file)

	if _, err := io.Copy(w, r); err != nil {
		return err
	}

	err := w.Flush()
	return err
}
