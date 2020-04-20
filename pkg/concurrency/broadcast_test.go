package concurrency_test

import (
	"testing"

	. "github.com/lxdlam/vertex/pkg/concurrency"
	"github.com/stretchr/testify/assert"
)

func TestBroadCast(t *testing.T) {
	bc := NewBroadCaster()

	msg := "Test Message!"

	for i := 0; i < 100; i++ {
		go func(id int) {
			item := bc.Get()

			if !assert.Equal(t, msg, item) {
				t.Fatalf("Gorountine %d received wrong message, got=%+v", id, item)
			}
		}(i)
	}

	bc.Set(msg)

	for i := 0; i < 100; i++ {
		go func(id int) {
			item := bc.Get()

			if !assert.Equal(t, msg, item) {
				t.Fatalf("Gorountine %d received wrong message, got=%+v", id, item)
			}
		}(i + 101)
	}
}
