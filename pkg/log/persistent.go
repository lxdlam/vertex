package log

import (
	"bufio"
	"errors"
	"os"
	"sync"
)

// ErrTargetFileInvalid will be raise when give PersistentFile.SaveTo a nil file.
var ErrTargetFileInvalid = errors.New("persistent: target file invalid")

// PersistentFile is a simple buffered file for write,
// all operations can be done safely across goroutines.
// It warps all bufio.Writer methods with a RWMutex.
type PersistentFile struct {
	mutex  *sync.RWMutex
	file   *os.File
	writer *bufio.Writer
}

// NewPersistentFile will returns a new PersistentFile instance warps the given file.
func NewPersistentFile(file *os.File) *PersistentFile {
	return &PersistentFile{
		mutex:  &sync.RWMutex{},
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
