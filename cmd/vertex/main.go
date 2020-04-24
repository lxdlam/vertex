package main

import (
	"fmt"
	"sync"
	"time"
)

const banner string = `
===================================================================
|                                                                 |
| '##::::'##:'########:'########::'########:'########:'##::::'##: |
|  ##:::: ##: ##.....:: ##.... ##:... ##..:: ##.....::. ##::'##:: |
|  ##:::: ##: ##::::::: ##:::: ##:::: ##:::: ##::::::::. ##'##::: |
|  ##:::: ##: ######::: ########::::: ##:::: ######:::::. ###:::: |
| . ##:: ##:: ##...:::: ##.. ##:::::: ##:::: ##...:::::: ## ##::: |
| :. ## ##::: ##::::::: ##::. ##::::: ##:::: ##:::::::: ##:. ##:: |
| ::. ###:::: ########: ##:::. ##:::: ##:::: ########: ##:::. ##: |
| :::...:::::........::..:::::..:::::..:::::........::..:::::..:: |
|                                                                 |
|   A Go implementation of Redis.     lxdlam(lxdlam@gmail.com)    |
|   Get the latest version at https://github.com/lxdlam/vertex    |
|                                                                 |
===================================================================
`

var wg sync.WaitGroup

func blockCall() int {
	for {
		fmt.Println("Blocking...")
	}

	wg.Done()
	return 15
}

func test(c chan struct{}) int {
	select {
	case <-c:
		wg.Done()
		return 10
	default:
		ret := blockCall()
		return ret
	}
}

func main() {
	fmt.Println(banner)
	wg.Add(2)

	ch := make(chan struct{})
	go test(ch)

	time.Sleep(100 * time.Millisecond)

	close(ch)

	wg.Wait()
}
