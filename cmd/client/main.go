package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/lxdlam/vertex/cmd/client/internal"

	"github.com/lxdlam/vertex/pkg/protocol"
)

const prompt = ">>> "

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	c, err := net.Dial("tcp", "127.0.0.1:6789")
	if err != nil {
		panic(err.Error())
	}

	defer c.Close()

	for {
		fmt.Print(prompt)
		if !scanner.Scan() {
			break
		}

		s := scanner.Text()
		if s == "QUIT" || s == "quit" {
			break
		}

		req, err := internal.FormatInput(s)
		if err != nil {
			fmt.Printf("Invalid input! err=%+v\n", err)
			continue
		}

		_, err = c.Write([]byte(req))
		if err != nil {
			fmt.Printf("Write to server failed! addr=%s, err=%+v\n", c.RemoteAddr().String(), err)
			continue
		}

		obj, err := protocol.Parse(bufio.NewReader(c))
		if err != nil {
			fmt.Printf("Read from server meets an error! addr=%s, err=%+v\n", c.RemoteAddr().String(), err)
			os.Exit(1)
		}

		fmt.Printf("%s\n", internal.FormatOutput(obj))
	}
}
