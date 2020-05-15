package command

// Command should provide some basic
type Command interface {
	// Key should return the command's key in client, e.g., SET or GET.
	Key() string
}
