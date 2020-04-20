package types_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	. "github.com/lxdlam/vertex/pkg/types"
)

type testCase struct {
	key   string
	value interface{}
}

func TestSimpleDataMap(t *testing.T) {
	m := NewSimpleDataMap()

	if m.Len() != 0 {
		t.Errorf("invalid size of map. m.Len()=%v", m.Len())
	}

	if !testBasicOperation(m, t) {
		t.Fatal("test basic operation for SimpleDataMap failed!")
	}
}

func TestLockedDataMap(t *testing.T) {
	m := NewLockedDataMap()

	if m.Len() != 0 {
		t.Errorf("invalid size of map. m.Len()=%v", m.Len())
	}

	if !testBasicOperation(m, t) {
		t.Fatal("test basic operation for LockedDataMap failed!")
	}

	if !testSyncOperation(m, t) {
		t.Fatal("test sync operation for LockedDataMap failed!")
	}
}

func TestSyncDataMap(t *testing.T) {
	m := NewSyncDataMap()

	if m.Len() != 0 {
		t.Errorf("invalid size of map. m.Len()=%v", m.Len())
	}

	if !testBasicOperation(m, t) {
		t.Fatal("test basic operation for SyncDataMap failed!")
	}

	if !testSyncOperation(m, t) {
		t.Fatal("test sync operation for SyncDataMap failed!")
	}
}

func testBasicOperation(m DataMap, t *testing.T) bool {
	// Insert
	insertTests := []testCase{
		{"a", "b"},
		{"c", 15},
	}

	if !testInsertAndGet(m, insertTests, t) {
		t.Errorf("test insert failed. testCase=%+v", insertTests)
		return false
	}

	// Replacement
	replaceTests := []testCase{
		{"a", 16},
		{"c", "d"},
	}

	if !testInsertAndGet(m, replaceTests, t) {
		t.Errorf("test replace failed. testCase=%+v", replaceTests)
		return false
	}

	for _, tt := range replaceTests {
		m.Remove(tt.key)
	}

	if m.Len() != 0 {
		t.Errorf("map still has item. m=%v", m)
		return false
	}

	return true
}

func testInsertAndGet(m DataMap, tests []testCase, t *testing.T) bool {
	for _, tc := range tests {
		m.Set(tc.key, tc.value)
	}

	for _, tc := range tests {
		if !m.Has(tc.key) {
			t.Errorf("key not exist. key=%s", tc.key)
			return false
		}

		if val, ok := m.Get(tc.key); !ok || !assert.Equal(t, val, tc.value) {
			t.Errorf("Get failed. key=%s, value=%+v", tc.key, tc.value)
			return false
		}
	}

	return true
}

func testSyncOperation(ctx DataMap, t *testing.T) bool {
	testCases := [][]testCase{
		{
			{"a", "b"},
			{"c", 1},
		},
		{
			{"e", 1.2},
			{"f", -7},
			{"g", "opopop"},
		},
		{
			{"hello", "world"},
		},
	}

	ch := make(chan bool)
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	for _, tcs := range testCases {
		go func(tc []testCase) {
			time.Sleep(time.Duration(r.Intn(10)) * time.Millisecond)
			ch <- testInsertAndGet(ctx, tc, t)
		}(tcs)
	}

	for i := 0; i < len(testCases); i++ {
		select {
		case <-ch:
			continue
		}
	}

	close(ch)

	for ret := range ch {
		if !ret {
			t.Error("test insert and get failed!")
			return false
		}
	}

	return true
}
