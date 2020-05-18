package container

const (
	MaxLevel    = 32
	ScaleFactor = 4
)

// SortedSetContainer is the sorted set data structure interface
type SortedSetContainer interface {
}

type skipListNode struct {
}

type skipList struct {
	key string
}

func NewSkipListSortedSetContainer(key string) SortedSetContainer {
	s := &skipList{key: key}

	return s
}
