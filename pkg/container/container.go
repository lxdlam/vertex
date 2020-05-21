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

type Containers interface {
}

type containers struct {
	global     StringMap
	lists      map[string]ListContainer
	hashes     map[string]HashContainer
	sets       map[string]SetContainer
	sortedSets map[string]SortedSetContainer
}

// NewContainers will return a new container that includes all (key, data structures) mapping
// for db to use
func NewContainers() Containers {
	return &containers{
		global:     NewStringMap(),
		lists:      make(map[string]ListContainer),
		hashes:     make(map[string]HashContainer),
		sets:       make(map[string]SetContainer),
		sortedSets: make(map[string]SortedSetContainer),
	}
}
