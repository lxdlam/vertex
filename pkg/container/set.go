package container

// SetContainer is the set data structure interface
type SetContainer interface {
	ContainerObject

	Add([]*StringContainer) int
	Delete([]*StringContainer) int
	IsMember(*StringContainer) bool

	Members() []*StringContainer
	RandomMember(int) []*StringContainer
	Pop(int) []*StringContainer

	Diff([]SetContainer) SetContainer
	Intersect([]SetContainer) SetContainer
	Union([]SetContainer) SetContainer

	Len() int
}

type setContainer struct {
	key       string
	container map[string]*StringContainer
}

// NewSetContainer returns a new hash container
func NewSetContainer(key string) SetContainer {
	return &setContainer{
		key:       key,
		container: make(map[string]*StringContainer),
	}
}

func (sc *setContainer) isContainer() {}

func (sc *setContainer) Key() string {
	return sc.key
}

func (sc *setContainer) Type() ContainerType {
	return SetType
}

func (sc *setContainer) Add(s []*StringContainer) int {
	var added int

	for _, item := range s {
		if _, ok := sc.container[item.String()]; !ok {
			sc.container[item.String()] = item
			added++
		}
	}

	return added
}

func (sc *setContainer) Delete(s []*StringContainer) int {
	var removed int

	for _, item := range s {
		if _, ok := sc.container[item.String()]; ok {
			delete(sc.container, item.String())
			removed++
		}
	}

	return removed
}

func (sc *setContainer) IsMember(s *StringContainer) bool {
	_, ok := sc.container[s.String()]
	return ok
}

func (sc *setContainer) Members() []*StringContainer {
	var ret []*StringContainer

	for _, item := range sc.container {
		ret = append(ret, item)
	}

	return ret
}

func (sc *setContainer) RandomMember(count int) []*StringContainer {
	var ret []*StringContainer

	if count < 0 {
		count = -count
	}

	for _, item := range sc.container {
		if count--; count < 0 {
			break
		}
		ret = append(ret, item)
	}

	return ret
}

func (sc *setContainer) Pop(count int) []*StringContainer {
	ret := sc.RandomMember(count)
	sc.Delete(ret)

	return ret
}

func (sc *setContainer) Diff(cs []SetContainer) SetContainer {
	var candidate []*StringContainer

	ret := NewSetContainer("anonymous")
	ret.Add(sc.Members())

	for _, c := range cs {
		for _, item := range c.Members() {
			if ret.IsMember(item) {
				candidate = append(candidate, item)
			}
		}

		ret.Delete(candidate)
		candidate = nil
	}

	return ret
}

func (sc *setContainer) Intersect(cs []SetContainer) SetContainer {
	var candidate []*StringContainer

	ret := NewSetContainer("anonymous")
	ret.Add(sc.Members())

	for _, c := range cs {
		for _, item := range c.Members() {
			if ret.IsMember(item) {
				candidate = append(candidate, item)
			}
		}

		_ = ret.Pop(ret.Len())
		ret.Add(candidate)
		candidate = nil
	}

	return ret
}

func (sc *setContainer) Union(cs []SetContainer) SetContainer {
	ret := NewSetContainer("anonymous")

	for _, c := range cs {
		ret.Add(c.Members())
	}
	ret.Add(sc.Members())

	return ret
}

func (sc *setContainer) Len() int {
	return len(sc.container)
}
