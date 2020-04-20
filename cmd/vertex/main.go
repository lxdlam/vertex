package main

import (
	"fmt"
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

func receive(ch chan bool, id int) {
	for {
		select {
		case val := <-ch:
			fmt.Printf("Gorountine %d received %v from channel.\n", id, val)
			return
		}
	}
}

func main() {
	// fmt.Println(banner)

	ch := make(chan bool)

	for i := 0; i < 10; i++ {
		go receive(ch, i)
	}

	ch <- true
	close(ch)

	time.Sleep(1 * time.Second)
}
