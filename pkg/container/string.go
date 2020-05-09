package container

// String for containers
// A simple type alias
type String string

const (
	DUMMY String = "DUMMY"
)

func (s String) AsRune() []rune {
	return []rune(s)
}

func (s String) AsByte() []byte {
	return []byte(s)
}
