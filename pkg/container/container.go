package container

type ContainerType int

const (
	_ ContainerType = iota
	StringType
	GlobalType
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
	Global() StringMap

	GetList(string) ListContainer
	GetOrCreateList(string) ListContainer

	GetHash(string) HashContainer
	GetOrCreateHash(string) HashContainer

	GetSet(string) SetContainer
	GetOrCreateSet(string) SetContainer

	//GetSortedSet(string) SortedSetContainer
	//GetOrCreateSortedSet(string) SortedSetContainer
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

func (c *containers) Global() StringMap {
	return c.global
}

func (c *containers) GetList(key string) ListContainer {
	l, ok := c.lists[key]

	if !ok {
		return nil
	}

	return l
}

func (c *containers) GetOrCreateList(key string) ListContainer {
	l, ok := c.lists[key]

	if !ok {
		l = NewLinkedListContainer(key)
		c.lists[key] = l
	}

	return l
}

func (c *containers) GetHash(key string) HashContainer {
	h, ok := c.hashes[key]

	if !ok {
		return nil
	}

	return h
}

func (c *containers) GetOrCreateHash(key string) HashContainer {
	h, ok := c.hashes[key]

	if !ok {
		h = NewHashContainer(key)
		c.hashes[key] = h
	}

	return h
}

func (c *containers) GetSet(key string) SetContainer {
	s, ok := c.sets[key]

	if !ok {
		return nil
	}

	return s
}

func (c *containers) GetOrCreateSet(key string) SetContainer {
	s, ok := c.sets[key]

	if !ok {
		s = NewSetContainer(key)
		c.sets[key] = s
	}

	return s
}

//func (c *containers) GetSortedSet(key string) SortedSetContainer {
//	s, ok := c.sets[key]
//
//	if !ok {
//		s = NewSetContainer(key)
//		c.sets[key] = s
//	}
//
//	return s
//}
//
//func (c *containers) GetOrCreateSortedSet(key string) SortedSetContainer {
//	s, ok := c.sets[key]
//
//	if !ok {
//		s = NewSortedSetContainer(key)
//		c.sets[key] = s
//	}
//
//	return s
//}
