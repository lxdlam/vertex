package command

// Command should provide some basic
type Command interface {
	// Key should return the command's key in client, e.g., SET or GET.
	Key() string
}

// StatefullCommand will have a context to do something necessary
type StatefullCommand interface {
	Command

	// Just implement to point out what it is
	statefullCommand()
}
