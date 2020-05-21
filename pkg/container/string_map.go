package container

import (
	"errors"
	"fmt"
)

var (
	// ErrGlobalLengthNotMatch will be raised if the keys and values in SET are not equal
	ErrGlobalLengthNotMatch = errors.New("global: the size of keys and values are not equal")

	// ErrKeyNotFound will be raised in any situation if the operation is related to the key
	ErrKeyNotFound = errors.New("global: key not found")
)

// StringMap is the global string container data structure interface
type StringMap interface {
	Set([]*StringContainer, []*StringContainer) error
	Get([]*StringContainer) []*StringContainer

	StringLen(*StringContainer) (int, error)
	Append(*StringContainer, *StringContainer) error

	Increase(*StringContainer, int64) (int64, error)
	Decrease(*StringContainer, int64) (int64, error)

	GetRange(*StringContainer, int, int) (*StringContainer, error)

	Len() int

	Exist(*StringContainer) bool
}

// NewStringMap will return a new global string map instance
func NewStringMap() StringMap {
	return &simpleStringMap{
		container: make(map[string]*StringContainer),
	}
}

type simpleStringMap struct {
	container map[string]*StringContainer
}

func (ssm *simpleStringMap) Set(keys, values []*StringContainer) error {
	l := len(keys)
	if l != len(values) {
		return ErrGlobalLengthNotMatch
	}

	for idx := 0; idx < l; idx++ {
		ssm.container[keys[idx].String()] = values[idx]
	}

	return nil
}

func (ssm *simpleStringMap) Get(keys []*StringContainer) []*StringContainer {
	l := len(keys)
	result := make([]*StringContainer, l)

	for idx := 0; idx < l; idx++ {
		if entry, ok := ssm.container[keys[idx].String()]; ok {
			result[idx] = entry
		} else {
			result[idx] = nil
		}
	}

	return result
}

func (ssm *simpleStringMap) StringLen(key *StringContainer) (int, error) {
	if entry, ok := ssm.container[key.String()]; ok {
		return entry.Len(), nil
	}

	return 0, ErrKeyNotFound
}

func (ssm *simpleStringMap) Append(key, str *StringContainer) error {
	panic("implement me")
}

func (ssm *simpleStringMap) Increase(key *StringContainer, increment int64) (int64, error) {
	if entry, ok := ssm.container[key.String()]; ok {
		if ret, err := entry.Increase(increment); err == nil {
			return ret, nil
		} else {
			return 0, fmt.Errorf("global: increase met an error. key=%s, err={%w}", key.String(), err)
		}
	}

	return 0, ErrKeyNotFound
}

func (ssm *simpleStringMap) Decrease(key *StringContainer, decrement int64) (int64, error) {
	if entry, ok := ssm.container[key.String()]; ok {
		if ret, err := entry.Decrease(decrement); err == nil {
			return ret, nil
		} else {
			return 0, fmt.Errorf("global: decrease met an error. key=%s, err={%w}", key.String(), err)
		}
	}

	return 0, ErrKeyNotFound
}

func (ssm *simpleStringMap) GetRange(key *StringContainer, start, end int) (*StringContainer, error) {
	if entry, ok := ssm.container[key.String()]; ok {
		if ret, err := entry.GetRange(start, end); err == nil {
			return ret, nil
		} else {
			return nil, fmt.Errorf("global: get range error. key=%s, range=[%d, %d], err={%w}", key.String(), start, end, err)
		}
	}

	return nil, ErrKeyNotFound
}

func (ssm *simpleStringMap) Len() int {
	return len(ssm.container)
}

func (ssm *simpleStringMap) Exist(key *StringContainer) bool {
	_, exist := ssm.container[key.String()]

	return exist
}
