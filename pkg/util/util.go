package util

// LexicalCompare takes two string, `lhs` and `rhs`, to compare them. If `lhs` < `rhs`, it will return `-idx`,
// else if `lhs` > `rhs`, it will return `idx`. If the two string are equal, it will return 0.
// `idx` here is the first words occurs not equal and where we can decide the order, note the
// `idx` will be 1-indexed, to distinguish from the result when the strings are equal.
func LexicalCompare(lhs, rhs string) int {
	lhsLen := len(lhs)
	rhsLen := len(rhs)

	length := lhsLen
	if lhsLen > rhsLen {
		length = rhsLen
	}

	for idx := 0; idx < length; idx++ {
		if lhs[idx] < rhs[idx] {
			return -idx - 1
		} else if lhs[idx] > rhs[idx] {
			return idx + 1
		}
	}

	if lhsLen == rhsLen {
		return 0
	} else if lhsLen < rhsLen {
		return -lhsLen
	}

	return rhsLen
}
