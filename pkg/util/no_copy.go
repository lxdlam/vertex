package util

// NoCopy is a simple compile time tag that to imply that a struct should not be copied.
type NoCopy struct {
}

// Lock will do nothing, just implement sync.Locker
func (*NoCopy) Lock() {}

// Unlock does the same thing like Lock
func (*NoCopy) Unlock() {}
