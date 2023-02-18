package main

import (
	"bufio"
	"codecrafters-redis-go/pkg/resp"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var Logger *log.Logger

func main() {
	// TODO: Allow writers to specific files specified in the config
	Logger = log.New(os.Stdout, "codecrafters-redis-go", log.LstdFlags)

	host := flag.String("h", "127.0.0.1", "Server hostname (default: 127.0.0.1)")
	port := flag.String("p", "6379", "Server port (default: 6379)")
	flag.Parse()

	dialer := net.Dialer{
		Timeout: time.Second * 3,
	}

	address := *host + ":" + *port
	conn, err := dialer.DialContext(context.TODO(), "tcp4", address)
	if err != nil {
		Logger.Fatalf("Failed to connect with the server. address: %v, error: %v", address, err.Error())
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
				Logger.Printf("Error occurred. error: %v", err.Error())
				continue
			}
			if len(response) > 1 {
				Logger.Printf("Multiple RESPs in a response are not supported")
				continue
			}

			fmt.Fprintf(os.Stdout, "> %v\n", toSpaceSeparated(response[0]))
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
