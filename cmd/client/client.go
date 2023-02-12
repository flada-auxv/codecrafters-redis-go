package main

import (
	"bufio"
	"codecrafters-redis-go/pkg/resp"
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	host := flag.String("h", "0.0.0.0", "Server hostname (default: 0.0.0.0)")
	port := flag.String("p", "6379", "Server port (default: 6379)")
	flag.Parse()

	dialer := net.Dialer{
		Timeout: time.Second * 3,
	}

	address := *host + ":" + *port
	conn, err := dialer.DialContext(context.TODO(), "tcp4", address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect with the server. address: %v, error: %v", address, err.Error())
		os.Exit(1)
	}

	for {
		stdinScanner := bufio.NewScanner(os.Stdin)

		for stdinScanner.Scan() {
			line := stdinScanner.Text()
			fields := strings.Fields(line)

			conn.Write(resp.EncodeArray(fields))

			responseReader := bufio.NewReader(conn)
			response, err := resp.Parse(responseReader)
			fmt.Println(response)
			if err != nil {
				fmt.Printf("Error occurred. error: %v\r\n", err.Error())
				continue
			}

			// TODO: response[1:]
			toSpaceSeparated(response[0])

			for _, v := range response {
				fmt.Printf("> %v\n", string(v.Data))
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
