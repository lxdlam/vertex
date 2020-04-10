package main

import (
	"fmt"

	"github.com/lxdlam/vertex/cmd/vertex/internal"
)

func main() {
	internal.ParseFlags()
	fmt.Println(internal.Verbose)
	fmt.Println(internal.Debug)
	fmt.Println(internal.ConfPath)
	fmt.Println("Welcome to Vertex Server!")
}
