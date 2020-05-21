package container

import (
	"errors"
	"fmt"
)

var (
	// ErrHashLengthNotMatch will be raised if the size of keys and values in SET not equal.
	ErrHashLengthNotMatch = errors.New("hash_container: the size of keys and values not equal")

	// ErrKeyNotExist will be raised if the key is not exist
	ErrKeyNotExist = errors.New("hash_container: the key is not exist")
)

// HashContainer is the hash data structure interface
type HashContainer interface {
	ContainerObject

	Set([]*StringContainer, []*StringContainer) error
	Get([]*StringContainer) []*StringContainer
	Exists(*StringContainer) bool
	Del(*StringContainer)

	Increase(*StringContainer, int64) (int64, error)

	Keys() []*StringContainer
	Values() []*StringContainer
	Entries() ([]*StringContainer, []*StringContainer)

	KeyLen(*StringContainer) (int, error)
	Len() int
}

type hashEntry struct {
	key   *StringContainer
	value *StringContainer
}

// So the real hashes is just a raw go map container
type hashContainer struct {
	key       string
	container map[string]*hashEntry
}

// NewHashContainer returns a new hash container
func NewHashContainer(key string) HashContainer {
	return &hashContainer{
		key:       key,
		container: make(map[string]*hashEntry),
	}
}

func (h *hashContainer) isContainer() {
}

func (h *hashContainer) Key() string {
	return h.key
}

func (h *hashContainer) Type() ContainerType {
	return HashType
}

func (h *hashContainer) Set(keys []*StringContainer, values []*StringContainer) error {
	l := len(keys)
	if l != len(values) {
		return ErrHashLengthNotMatch
	}

	for i := 0; i < l; i++ {
		if entry, ok := h.container[keys[i].String()]; !ok {
			h.container[keys[i].String()] = &hashEntry{
				key:   keys[i],
				value: values[i],
			}
		} else {
			entry.value = values[i]
		}
	}

	return nil
}

func (h *hashContainer) Get(keys []*StringContainer) []*StringContainer {
	var ret []*StringContainer

	for _, key := range keys {
		if entry, ok := h.container[key.String()]; !ok {
			ret = append(ret, nil)
		} else {
			ret = append(ret, entry.value)
		}
	}

	return ret
}

func (h *hashContainer) Exists(key *StringContainer) bool {
	_, ret := h.container[key.String()]

	return ret
}

func (h *hashContainer) Del(key *StringContainer) {
	delete(h.container, key.String())
}

func (h *hashContainer) Increase(key *StringContainer, increment int64) (int64, error) {
	entry, ok := h.container[key.String()]
	if !ok {
		return 0, ErrKeyNotExist
	}

	ret, err := entry.value.Increase(increment)

	if err != nil {
		return 0, fmt.Errorf("hash_container: increase with error. err={%w}", err)
	}

	return ret, nil
}

func (h *hashContainer) Keys() []*StringContainer {
	var ret []*StringContainer

	for _, entry := range h.container {
		ret = append(ret, entry.key)
	}

	return ret
}

func (h *hashContainer) Values() []*StringContainer {
	var ret []*StringContainer

	for _, entry := range h.container {
		ret = append(ret, entry.value)
	}

	return ret
}

func (h *hashContainer) Entries() ([]*StringContainer, []*StringContainer) {
	var keys, values []*StringContainer

	for _, entry := range h.container {
		keys = append(keys, entry.key)
		values = append(values, entry.value)
	}

	return keys, values
}

func (h *hashContainer) KeyLen(key *StringContainer) (int, error) {
	entry, ok := h.container[key.String()]

	if !ok {
		return 0, ErrKeyNotExist
	}

	return entry.value.Len(), nil
}

func (h *hashContainer) Len() int {
	return len(h.container)
}
