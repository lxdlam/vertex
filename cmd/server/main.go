package main

import (
	"fmt"
	"os"

	"github.com/lxdlam/vertex/pkg/common"
	"github.com/lxdlam/vertex/pkg/network"
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
	s := network.NewServer()
	c := common.Config{
		LogPath:  "./log",
		LogLevel: "DEBUG",
		Port:     6789,
	}

	common.InitLog(c, true)

	if !s.Init(c) {
		fmt.Fprintln(os.Stderr, "init server failed!")
		os.Exit(1)
	}

	s.Serve()
}
