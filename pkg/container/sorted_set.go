package container

//import (
//	"github.com/lxdlam/vertex/pkg/util"
//	"math/rand"
//)
//
//const (
//	MaxLevel    = 32
//	ScaleFactor = 4
//)
//
//// SortedSetContainer is a set which ensures uniqueness of member and sorted by the score.
//// Multiple score are allowed.
//type SortedSetContainer interface {
//}
//
//type entryNode struct {
//	score float64
//	data  *StringContainer
//	prev  *entryNode
//	next  *entryNode
//}
//
//func (en *entryNode) delete() (*StringContainer, float64) {
//	data := en.data
//	score := en.score
//
//	en.prev.next = en.next
//	en.next.prev = en.prev
//	en.data = nil
//	en.next = nil
//	en.prev = nil
//
//	return data, score
//}
//
//type entryList struct {
//	head *entryNode
//	tail *entryNode
//	size int
//}
//
//func (el *entryList) len() int {
//	return el.size
//}
//
//// always add at tail
//func addEntry(data *StringContainer, score float64) *entryNode {
//
//}
//
//func newEntryList() *entryList {
//	el := &entryList{
//	}
//}
//
//type skipListNode struct {
//	score     float64
//	entryList *entryList
//	maxLevel  int
//	next      [32]*skipListNode
//	prev      [32]*skipListNode
//}
//
//type skipList struct {
//	key  string
//	head *skipListNode
//	tail *skipListNode
//	size int
//	rand *rand.Rand
//}
//
//// NewSkipListSortedSet returns a new sorted set which implementation is skip list
//func NewSkipListSortedSet(key string) SortedSetContainer {
//	s := &skipList{
//		key:  key,
//		size: 0,
//		rand: util.GetNewRandom(),
//	}
//
//	head := &skipListNode{
//		score:    -1,
//		entries:  nil,
//		maxLevel: 0,
//	}
//
//	tail := &skipListNode{
//		score:    -1,
//		entries:  nil,
//		maxLevel: 0,
//	}
//
//	head.next[0] = tail
//	tail.prev[0] = head
//	s.head = head
//	s.tail = tail
//
//	return s
//}
//
////// lowerBound will found the first node that the score is larger or equal to the node
////func (sl *skipList) lowerBound(score float64) *skipListNode {
////
////}
////
////// upperBound will found the first node that the score is larger to the node
////func (sl *skipList) upperBound(score float64) *skipListNode {
////
////}
