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

func main() {
	fmt.Println(banner)

	ch := make(chan int)
	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func(idx int) {
			ch <- idx
			wg.Done()
		}(i)
	}

	var nums []int

	go func() {
		for item := range ch {
			nums = append(nums, item)
		}
	}()

	wg.Wait()

	go func() {
		for _, item := range nums {
			fmt.Println(item)
		}
		fmt.Println(len(nums))
	}()

	time.Sleep(time.Millisecond)
}
