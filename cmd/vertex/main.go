package main

import (
	"github.com/lxdlam/vertex/pkg/common"
	"github.com/lxdlam/vertex/pkg/concurrency"
	"github.com/lxdlam/vertex/pkg/network"
	"github.com/lxdlam/vertex/pkg/types"
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

func startDummyListener() {
	receiver, _ := concurrency.GetEventBus().Subscribe("request", "dummy")

	for {
		e, _ := receiver.Receive()

		d := e.Data().(types.DataMap)
		resp, _ := d.Get("request")
		d.Set("response", resp)

		concurrency.GetEventBus().Publish("response", d, nil)
	}
}

func main() {
	s := network.NewServer()
	c := common.Config{
		LogPath:  "./log",
		LogLevel: "DEBUG",
		Port:     6789,
	}
	s.Init(c)

	go startDummyListener()

	s.Serve()
}
