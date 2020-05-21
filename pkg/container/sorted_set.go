package container

import (
	"errors"
	"math"

	"github.com/lxdlam/vertex/pkg/util"
)

const (
	maxLevel    = 32
	scaleFactor = 4
)

var (
	// ErrSortedSetLengthNotMatch will be raised if the size of scores and entries in ZADD are not equal.
	ErrSortedSetLengthNotMatch = errors.New("sorted_set_container: the size of scores and entries are not equal")

	// ErrEntryNotFound will be raised in all the entry related operations when the entry is not found
	ErrEntryNotFound = errors.New("sorted_set_container: entry not found")

	// ErrSortedSetEmpty will be raised in PopMin/Max
	ErrSortedSetEmpty = errors.New("sorted_set_container: empty sorted set")
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

type skipListLevel struct {
	node *skipListNode
	span int // span is the skipped node from the predecessor
}

type skipListNode struct {
	score float64
	data  *StringContainer
	next  []skipListLevel
	prev  *skipListNode
}

func (sln *skipListNode) release() {
	sln.data = nil
	sln.prev = nil

	for _, item := range sln.next {
		item.node = nil
	}
}

type skipList struct {
	key   string
	head  *skipListNode
	tail  *skipListNode
	level int
	set   map[string]*skipListNode
}

func getRandomLevel() int {
	level := 1
	for util.GetGlobalRandom().Int()%scaleFactor == 0 {
		level++
	}

	if level > maxLevel {
		level = maxLevel
	}

	return level
}

func (sl *skipList) insert(score float64, entry *StringContainer) {
	if _, exist := sl.set[entry.String()]; exist {
		return
	}

	// update[level] is the node that is the predecessor of cur at level
	// rank[level] is the total skipped node at level
	// so rank[0] is the total skipped node, which is the rank of the new node
	update := make([]*skipListNode, maxLevel)
	rank := make([]int, maxLevel)

	// first walk the skip list and record the predecessors and the rank
	cur := sl.head
	for level := sl.level - 1; level >= 0; level-- {
		if level == sl.level-1 {
			rank[level] = 0
		} else {
			rank[level] = rank[level+1]
		}

		for cur.next[level].node != sl.tail &&
			(cur.next[level].node.score < score ||
				(cur.next[level].node.score == score && cur.next[level].node.data.CompareTo(entry) < 0)) {
			rank[level] += cur.next[level].span
			cur = cur.next[level].node
		}

		update[level] = cur
	}

	// gen a new level, if the new level is higher than current, just using head as predecessor
	newLevel := getRandomLevel()
	if newLevel > sl.level {
		for level := sl.level; level < newLevel; level++ {
			rank[level] = 0
			update[level] = sl.head
			update[level].next[level].span = sl.Len()
		}

		sl.level = newLevel
	}

	cur = &skipListNode{
		score: score,
		data:  entry,
		next:  nil,
	}

	// add nodes, do below two things:
	// 1. add the node as we do in linked list
	// 2. adjust the correct span, pictures below
	//         +----------------------------------------------------------------------------------------------------------------------+
	//         |                                           update[level].next[level].span                                             |
	//         |                                                                                                                      |
	// +-------+-------+                           +---------------+                                                          +-------v-------+
	// |               |                           |               |                                                          |               |
	// | update[level] |                           |      cur      |                                                          | update[level] |
	// |               |                           |               |                                                          | .next[level]  |
	// |               |                           |               |                                                          | .node         |
	// |               |                           |               |                                                          |               |
	// |      +------------------+-------------+------------>   +-----------------------------------+----------+---------------------->       |
	// |               | rank[0] - rank[level] + 1 |               | update[level].next[level].span - (rank[0] - rank[level]) |               |
	// |               |                           |               |                                                          |               |
	// +---------------+                           +---------------+                                                          +---------------+
	for level := 0; level < sl.level; level++ {
		cur.next = append(cur.next, skipListLevel{})
		cur.next[level].node = update[level].next[level].node
		update[level].next[level].node = cur

		cur.next[level].span = update[level].next[level].span - (rank[0] - rank[level])
		update[level].next[level].span = rank[0] - rank[level] + 1
	}

	// those level are not touched, they just skipped one more node
	// so we increase each span by 1
	for level := newLevel; level < sl.level; level++ {
		update[level].next[level].span++
	}

	cur.prev = update[0]
	cur.next[0].node.prev = cur

	sl.set[entry.String()] = cur
}

// the prev node is only point to the 0 layer, but we need to adjust all node in all occured layer
// so a update slice is needed, resolve it by different method
func (sl *skipList) delete(node *skipListNode, update []*skipListNode) {
	for level := 0; level < sl.level; level++ {
		if update[level].next[level].node == node {
			update[level].next[level].node = node.next[level].node
			update[level].next[level].span += node.next[level].span - 1
		} else {
			update[level].next[level].span--
		}
	}

	for sl.level > 1 && sl.head.next[sl.level-1].node == sl.tail {
		sl.level--
	}

	delete(sl.set, node.data.String())

	node.release()
}

func (sl *skipList) deleteNode(node *skipListNode) {
	update := make([]*skipListNode, maxLevel)

	cur := sl.head
	for level := sl.level - 1; level >= 0; level-- {
		for cur.next[level].node != sl.tail && cur.next[level].node != node {
			cur = cur.next[level].node
		}

		update[level] = cur
	}

	sl.delete(node, update)
}

func (sl *skipList) getRank(node *skipListNode) int {
	rank := 0
	cur := sl.head
	for level := sl.level - 1; level >= 0; level-- {
		for cur.next[level].node != sl.tail || cur.next[level].node != node {
			rank += cur.next[level].span
			cur = cur.next[level].node
		}
	}

	return rank
}

// NewSortedSetContainer returns a new SortedSetContainer which implementation is
// skip list.
func NewSortedSetContainer(key string) SortedSetContainer {
	s := &skipList{key: key, level: 1}

	head := &skipListNode{
		score: -math.MaxFloat64,
		prev:  nil,
	}

	tail := &skipListNode{
		score: math.MaxFloat64,
		next:  nil,
	}

	for idx := 0; idx < maxLevel; idx++ {
		head.next = append(head.next, skipListLevel{
			node: tail,
			span: 0,
		})
	}
	tail.prev = head

	s.head = head
	s.tail = tail
	s.set = make(map[string]*skipListNode)

	return s
}

func (sl *skipList) isContainer() {}

func (sl *skipList) Key() string {
	return sl.key
}

func (sl *skipList) Type() ContainerType {
	return SortedSetType
}

func (sl *skipList) Add(scores []float64, entries []*StringContainer) error {
	l := len(scores)
	if l != len(entries) {
		return ErrSortedSetLengthNotMatch
	}

	for idx := 0; idx < l; idx++ {
		sl.insert(scores[idx], entries[idx])
	}

	return nil
}

func (sl *skipList) Del(entries []*StringContainer) {
	for _, entry := range entries {
		if node, ok := sl.set[entry.String()]; ok {
			sl.deleteNode(node)
		}
	}
}

func (sl *skipList) Count(min, max float64) int {
	panic("implement me")
}

func (sl *skipList) Score(entry *StringContainer) (float64, error) {
	if entryNode, ok := sl.set[entry.String()]; !ok {
		return 0, ErrEntryNotFound
	} else {
		return entryNode.score, nil
	}
}

func (sl *skipList) IncreaseBy(entry *StringContainer, increment float64) (float64, error) {
	node, ok := sl.set[entry.String()]

	if !ok {
		return 0, ErrEntryNotFound
	}

	newScore := node.score + increment
	data := node.data
	sl.deleteNode(node)
	sl.insert(newScore, data)

	return newScore, nil
}

func (sl *skipList) PopMin() (*StringContainer, error) {
	if len(sl.set) == 0 {
		return nil, ErrSortedSetEmpty
	}

	data := sl.head.next[0].node.data
	sl.deleteNode(sl.head.next[0].node)

	return data, nil
}

func (sl *skipList) PopMax() (*StringContainer, error) {
	if len(sl.set) == 0 {
		return nil, ErrSortedSetEmpty
	}

	data := sl.tail.prev.data
	sl.deleteNode(sl.tail.prev)

	return data, nil
}

func (sl *skipList) Rank(entry *StringContainer) (int, error) {
	if node, ok := sl.set[entry.String()]; !ok {
		return -1, ErrEntryNotFound
	} else {
		return sl.getRank(node), nil
	}
}

func (sl *skipList) RangeByRank(start, end int) ([]*StringContainer, error) {
	panic("implement me")
}

func (sl *skipList) RangeByScore(min, max float64) ([]*StringContainer, error) {
	panic("implement me")
}

func (sl *skipList) DelRangeByRank(start, end int) error {
	panic("implement me")
}

func (sl *skipList) DelRangeByScore(min, max float64) error {
	panic("implement me")
}

func (sl *skipList) Len() int {
	return len(sl.set)
}
