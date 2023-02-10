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

	for {
		stdinScanner := bufio.NewScanner(os.Stdin)
		for stdinScanner.Scan() {
			line := stdinScanner.Text()
			fields := strings.Fields(line)

			conn.Write(resp.EncodeArray(fields))

			responseReader := bufio.NewReader(conn)
			response, err := resp.Parse(responseReader)
			if err != nil {
				panic("hi")
			}

			// TODO: response[1:]
			toSpaceSeparated(response[0])

			for _, v := range response {
				fmt.Printf("> %#v\n", string(v.Data))
			}
		}
	}
}

func toSpaceSeparated(r resp.RESP) string {
	switch r.Type {
	case resp.RESPArray:
		var str string
		for i, v := range r.Array {
			if i == 0 {
				str = str + toSpaceSeparated(v)
			} else {
				str = str + " " + toSpaceSeparated(v)
			}
		}
		return str
	case resp.RESPBulkString, resp.RESPError, resp.RESPInteger, resp.RESPSimpleString:
		return string(r.Data)
	default:
		return ""
	}
}
