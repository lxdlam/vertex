package container

type ContainerType int

const (
	AnyType ContainerType = iota
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

type container struct {
	index      int
	lists      map[string]ListContainer
	hashes     map[string]HashContainer
	sets       map[string]SetContainer
	sortedSets map[string]SortedSetContainer
}
