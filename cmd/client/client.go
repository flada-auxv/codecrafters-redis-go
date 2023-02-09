package main

import (
	"bufio"
	"codecrafters-redis-go/pkg/resp"
	"context"
	"fmt"
	"net"
	"os"
	"strings"
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

	var line

	for {
		stdinScanner := bufio.NewScanner(os.Stdin)
		for stdinScanner.Scan() {
			line := stdinScanner.Text()
			fields := strings.Fields(line)
			os.Stdin.

			// TODO:
			conn.Write([]byte{fmt.Sprintf("*1\r\n$%v\r\n%v\r\n", len(fields[0]), fields[0])},)

			responseReader := bufio.NewReader(conn)
			response, err := resp.Parse(responseReader)
			if err != nil {
				panic("hi")
			}

			for i, resp := range response {

			}

		}
	}

	for _, v := range response {
		fmt.Printf("res: %#v\n", string(v.Data))
	}
}
