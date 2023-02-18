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

	fmt.Fprintf(os.Stdout, "%v>", address)

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

		fmt.Fprintf(os.Stdout, "%v\n", toReadable(response[0]))
		fmt.Fprintf(os.Stdout, "%v>", address)
	}
}

func toReadable(r resp.RESP) string {
	switch r.Type {
	case resp.RESPArray:
		if len(r.Array) == 0 {
			return "(empty array)"
		}

		var str string
		for i, v := range r.Array {
			header := fmt.Sprintf("%v)", i+1)
			if i == 0 {
				str = header + " " + toReadable(v)
			} else {
				str = str + "\n" + header + " " + toReadable(v)
			}
		}
		return str

	case resp.RESPBulkString:
		if string(r.Data) == "" {
			return "(nil)"
		}
		return "\"" + string(r.Data) + "\""

	case resp.RESPError:
		return "(error) " + string(r.Data)

	case resp.RESPInteger:
		return "(integer) " + string(r.Data)

	case resp.RESPSimpleString:
		return string(r.Data)

	default:
		return ""
	}
}
