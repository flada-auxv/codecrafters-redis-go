package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		bufio.NewReader(conn)
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error while reading from connection", err.Error())
			os.Exit(1)
		}
		resp, err := Parse(bytes.Trim(buf, "\x00"))
		if err != nil {
			fmt.Println("Error while parsing request", err.Error())
			os.Exit(1)
		}

		exec(conn, resp)
	}
}

type RESP struct {
	Count int
	Data  []byte
	Raw   []byte
	Type  byte
}

const (
	RESPArray        = '*'
	RESPBulkString   = '$'
	RESPError        = '-'
	RESPInteger      = ':'
	RESPSimpleString = '+'
)

func exec(conn net.Conn, resps []RESP) {
	if resps[0].Type != RESPArray {
		panic("Currently, only array of bulk string is supported")
	}

	// TODO: should be execed according to the type of first RESP
	arr := resps[0]
	cmd := resps[1]
	augs := resps[2 : 2+arr.Count-1]

	// TODO: The redis command group seems to be case insensitive and uses uppercase, but the codecrafters send it in lowercase...?
	switch string(cmd.Data) {
	case "ECHO", "echo":
		message := []byte{}
		for _, v := range augs {
			message = append(message, v.Raw...)
		}
		conn.Write(message)
	case "PING", "ping":
		conn.Write([]byte("+PONG\r\n"))
	default:
		panic("not implemented")
	}
}
