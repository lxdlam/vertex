package util

import (
	"math/rand"
	"sync"
	"time"
)

var once sync.Once
var source rand.Source
var globalRand *rand.Rand

// GetNewRandom returns a new random generator while the seed is properly set
func GetNewRandom() *rand.Rand {
	r := rand.New(rand.NewSource(GetGlobalRandom().Int63()))
	return r
}

// GetGlobalRandom returns the global random generator while the seed is properly set
func GetGlobalRandom() *rand.Rand {
	once.Do(func() {
		source = rand.NewSource(time.Now().UnixNano())
		globalRand = rand.New(source)
	})

	return globalRand
}
