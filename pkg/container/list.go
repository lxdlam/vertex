package container

import (
	"errors"

	"github.com/lxdlam/vertex/pkg/common"
	"github.com/lxdlam/vertex/pkg/util"
)

var (
	// ErrEmptyList will be raised when using PopHead/PopTail
	ErrEmptyList = errors.New("empty list")

	// ErrNoSuchPivot will be raised when using Insert
	ErrNoSuchPivot = errors.New("no such pivot")

	// ErrOutOfRange will be raised when using Set/Index
	ErrOutOfRange = errors.New("index out of range")
)

// LinkedList is the list container interface.
type LinkedList interface {
	PushHead([]String) (int, error)
	PushTail([]String) (int, error)

	PopHead() (String, error)
	PopTail() (String, error)

	Insert(String, String, bool) (int, error)
	Set(int, String) error

	Remove(int) error
	Trim(int, int) error

	Index(int) (String, error)
	Range(int, int) ([]String, error)

	Len() int
}

type listNode struct {
	data String
	next *listNode
	prev *listNode
}

type linkedList struct {
	key  string
	head *listNode
	tail *listNode
	size int
}

func insertBefore(n *listNode, data String) error {
	if n == nil {
		return common.Errorf("linked_list: insert a element before a nil listNode")
	}

	tmp := &listNode{
		data: data,
		prev: n.prev,
		next: n,
	}
	n.prev.next = tmp
	n.prev = tmp

	return nil
}

func insertAfter(n *listNode, data String) error {
	if n == nil {
		return common.Errorf("linked_list: insert a element after a nil listNode")
	}

	tmp := &listNode{
		data: data,
		prev: n,
		next: n.next,
	}
	n.next.prev = tmp
	n.next = tmp

	return nil
}

func delete(n *listNode) (String, error) {
	if n == nil {
		return DUMMY, common.Errorf("linked_list: delete a nil listNode")
	}

	data := n.data
	n.prev.next = n.next
	n.next.prev = n.prev
	n = nil

	return data, nil
}

func (l *linkedList) extract(index int) *listNode {
	normIndex := util.NewIndex(index).Resolve(l.size)

	if normIndex != -1 {
		cur := l.head
		for idx := -1; idx < normIndex; idx++ {
			cur = cur.next
		}

		return cur.next
	}

	return nil
}

// Avoid traverse twice
func (l *linkedList) extractSegment(left, right int) (*listNode, *listNode) {
	right = util.NewIndex(right).ResolveRaw(l.size)

	if right > l.size {
		right = l.size - 1
	}

	normLeft, normRight := util.NewSlice(left, right).Resolve(l.size)

	if normLeft != -1 && normRight != -1 {
		var leftNode, rightNode *listNode

		cur := l.head

		for idx := -1; idx < normRight; idx++ {
			if idx == normLeft {
				leftNode = cur
			} else if idx == normRight {
				rightNode = cur
			}

			cur = cur.next
		}

		return leftNode, rightNode
	}

	return nil, nil
}

// mark all pointer to nil, so the GC can properly release them.
// move from head to tail
func releaseList(n *listNode) int {
	if n == nil {
		return 0
	}

	length := releaseList(n.next)
	n.prev = nil
	n.next = nil

	return length + 1
}

// NewLinkedList will return a new linked list instance which is assigned of the give key
func NewLinkedList(key string) LinkedList {
	l := &linkedList{
		key:  key,
		head: nil,
		tail: nil,
		size: 0,
	}

	// Two dummy listNodes overhead here to simplify all process
	head := &listNode{data: DUMMY, prev: nil}
	tail := &listNode{data: DUMMY, next: nil}

	head.next = tail
	tail.prev = head
	l.head = head
	l.tail = tail

	return l
}

func (l *linkedList) PushHead(data []String) (int, error) {
	for _, item := range data {
		if err := insertAfter(l.head, item); err != nil {
			return -1, err
		}
		l.size++
	}

	return l.size, nil
}

func (l *linkedList) PushTail(data []String) (int, error) {
	for _, item := range data {
		if err := insertBefore(l.tail, item); err != nil {
			return -1, err
		}
		l.size++
	}

	return l.size, nil
}

func (l *linkedList) PopHead() (String, error) {
	if l.head.next == l.tail {
		return DUMMY, ErrEmptyList
	}

	return delete(l.head.next)
}

func (l *linkedList) PopTail() (String, error) {
	if l.tail.prev == l.head {
		return DUMMY, ErrEmptyList
	}

	return delete(l.tail.prev)
}

func (l *linkedList) Insert(pivot, data String, after bool) (int, error) {
	cur := l.head

	for cur != l.tail {
		if cur.data == pivot {
			if after {
				if err := insertAfter(cur, data); err != nil {
					return -2, err
				}
				l.size++
			} else {
				if err := insertBefore(cur, data); err != nil {
					return -2, err
				}
				l.size++
			}

			return l.size, nil
		}

		cur = cur.next
	}

	return -1, ErrNoSuchPivot
}

func (l *linkedList) Set(index int, data String) error {
	n := l.extract(index)

	if n == nil {
		return ErrOutOfRange
	}

	n.data = data
	return nil
}

func (l *linkedList) Remove(index int) error {
	n := l.extract(index)

	if n == nil {
		return ErrOutOfRange
	}

	n.prev.next = n.next
	n.next.prev = n.prev
	n.prev = nil
	n.next = nil
	l.size--

	return nil
}

func (l *linkedList) Index(index int) (String, error) {
	n := l.extract(index)

	if n == nil {
		return DUMMY, ErrOutOfRange
	}

	return n.data, nil
}

func (l *linkedList) Range(left, right int) ([]String, error) {
	leftNode, rightNode := l.extractSegment(left, right)
	if leftNode != nil && rightNode != nil {
		var result []String
		for leftNode != rightNode {
			result = append(result, leftNode.data)
			leftNode = leftNode.next
		}
	}

	return nil, ErrOutOfRange
}

func (l *linkedList) Trim(left, right int) error {
	leftNode, rightNode := l.extractSegment(left, right)
	if leftNode != nil && rightNode != nil {
		// release head.next->leftNode.prev
		leftNode.prev.next = nil
		l.size -= releaseList(l.head.next)

		// release rightNode.next->tail.prev
		l.tail.prev.next = nil
		l.size -= releaseList(rightNode.next)

		// link head->leftNode
		l.head.next = leftNode
		leftNode.prev = l.head

		// link rightNode->tail
		l.tail.prev = rightNode
		rightNode.next = l.tail

		return nil
	}

	return ErrOutOfRange
}

func (l *linkedList) Len() int {
	return l.size
}
