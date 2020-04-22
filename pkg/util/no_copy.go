package util

// NoCopy is a simple compile time tag that to avoid copy a struct
type NoCopy struct {
}

// Lock will do nothing, just implement sync.Locker
func (*NoCopy) Lock() {}

// Unlock does the same thing like Lock
func (*NoCopy) Unlock() {}
