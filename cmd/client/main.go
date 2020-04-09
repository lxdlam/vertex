package main

import (
	"sync"
)

func main() {
	var m sync.Map
	m.Store(1, 2)
}
