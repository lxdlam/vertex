package util

// Index is an index container to help resolving a specific index of
type Index struct {
	index int
}

// NewIndex is a shorthand for construct new index instance
func NewIndex(index int) *Index {
	return &Index{
		index: index,
	}
}

// ResolveRaw will resolve the index into the given range [0, size), the right side is open.
func (i *Index) ResolveRaw(size int) int {
	if i.index < 0 {
		return i.index + size
	}

	return i.index
}

// Resolve will using the result of ResolveRaw, but -1 will be returned if
// the resolved index is not landing in the given range.
func (i *Index) Resolve(size int) int {
	index := i.ResolveRaw(size)

	if index < 0 || index >= size {
		return -1
	}

	return index
}

// Slice will include two Index to resolve a given segment on the range
type Slice struct {
	left, right *Index
}

// NewSlice is a shorthand for construct new slice instance
func NewSlice(left, right int) *Slice {
	return &Slice{
		left:  NewIndex(left),
		right: NewIndex(right),
	}
}

// ResolveRaw will just return results of Index.ResolveRaw on the left and right segment
func (s *Slice) ResolveRaw(size int) (int, int) {
	return s.left.ResolveRaw(size), s.right.ResolveRaw(size)
}

// Resolve will report the correct segment when given the real length of the whole range.
//
// After resolving, below situation will be considered invalid:
//    - The left is larger than right, e.g., [10, 5]
//    - One of the index is out of range, e.g., the length is 10 but the segment is [20, 30]
// If the slice is considered invalid, (-1, -1) will be returned.
func (s *Slice) Resolve(size int) (int, int) {
	left, right := s.ResolveRaw(size)

	if left < 0 || right < 0 || left >= size || right >= size || left > right {
		return -1, -1
	}

	return left, right
}
