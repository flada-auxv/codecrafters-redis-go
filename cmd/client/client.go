package main

import (
	"bufio"
	"codecrafters-redis-go/pkg/resp"
	"context"
	"fmt"
	"net"
	"os"
	"time"
)

// ECHO hey
// > hey

func main() {
	dialer := net.Dialer{
		Timeout: time.Second * 3,
	}

	// TODO: parse from command options
	address := "0.0.0.0:6379"
	conn, err := dialer.DialContext(context.TODO(), "tcp4", address)
	if err != nil {
		fmt.Println(err.Error())
		panic("hi")
	}

	fmt.Println("aaa")

	for {
		inputReader := bufio.NewReader(os.Stdin)
		// TODO: not work
		input, err := resp.Parse(inputReader)
		if err != nil {
			panic("hi")
		}

		// TODO: command
		conn.Write([]byte("*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n"))

		responseReader := bufio.NewReader(conn)
		response, err := resp.Parse(responseReader)
		if err != nil {
			panic("hi")
		}

		fmt.Println(response)
	}
}
