package container

type ContainerType int

const (
	_ ContainerType = iota
	StringType
	LinkedListType
	HashType
	SetType
	SortedSetType
	ContainerTypeLength
)

type ContainerObject interface {
	isContainer()
	Key() string
	Type() ContainerType
}

type Container interface {
}
