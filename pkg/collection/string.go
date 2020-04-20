package collection

// String for containers
// A simple type alias
type String string

func (s String) AsRune() []rune {
	return []rune(s)
}

func (s String) AsByte() []byte {
	return []byte(s)
}

