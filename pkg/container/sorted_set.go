package container

const (
	maxLevel    = 32
	scaleFactor = 4
)

// SortedSetContainer is the sorted set data structure interface
type SortedSetContainer interface {
	ContainerObject

	Add([]float64, []*StringContainer) error
	Del([]*StringContainer)

	Count(float64, float64) int
	Score(*StringContainer) (float64, error)

	IncreaseBy(*StringContainer, float64) (float64, error)

	PopMin() (*StringContainer, error)
	PopMax() (*StringContainer, error)

	Rank(*StringContainer) (int, error)

	RangeByRank(int, int) ([]*StringContainer, error)
	RangeByScore(float64, float64) ([]*StringContainer, error)

	DelRangeByRank(int, int) error
	DelRangeByScore(float64, float64) error

	Len() int
}

type skipListNode struct {
	score float64
	data  *StringContainer
	next  []*skipListNode
	prev  []*skipListNode
	level int
}

type skipList struct {
	key  string
	head *skipListNode
	tail *skipListNode
	set  map[string]*skipListNode
	size int
}

func (*skipList) insert(score float64, entry *StringContainer) {

}

func NewSkipListSortedSetContainer(key string) SortedSetContainer {
	s := &skipList{key: key}

	return s
}

func (s *skipList) isContainer() {}

func (s *skipList) Key() string {
	return s.key
}

func (s *skipList) Type() ContainerType {
	return SortedSetType
}

func (s *skipList) Add(scores []float64, entries []*StringContainer) error {
	panic("implement me")
}

func (s *skipList) Del(entries []*StringContainer) {
	panic("implement me")
}

func (s *skipList) Count(min, max float64) int {
	panic("implement me")
}

func (s *skipList) Score(entry *StringContainer) (float64, error) {
	panic("implement me")
}

func (s *skipList) IncreaseBy(entry *StringContainer, increment float64) (float64, error) {
	panic("implement me")
}

func (s *skipList) PopMin() (*StringContainer, error) {
	panic("implement me")
}

func (s *skipList) PopMax() (*StringContainer, error) {
	panic("implement me")
}

func (s *skipList) Rank(entry *StringContainer) (int, error) {
	panic("implement me")
}

func (s *skipList) RangeByRank(start, end int) ([]*StringContainer, error) {
	panic("implement me")
}

func (s *skipList) RangeByScore(min, max float64) ([]*StringContainer, error) {
	panic("implement me")
}

func (s *skipList) DelRangeByRank(start, end int) error {
	panic("implement me")
}

func (s *skipList) DelRangeByScore(min, max float64) error {
	panic("implement me")
}

func (s *skipList) Len() int {
	return s.size
}
