package main

import (
	"bufio"
	"fmt"
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
		reader := bufio.NewReader(conn)
		resp, err := Parse(reader)
		if err != nil {
			fmt.Println("Error while parsing request", err.Error())
			os.Exit(1)
		}

		exec(conn, resp)
	}
}

func exec(conn net.Conn, resps []RESP) {
	if resps[0].Type != RESPArray {
		panic("Currently, only array of bulk string is supported")
	}

	// TODO: should be execed according to the type of first RESP
	// TODO: The redis command group seems to be case insensitive and uses uppercase, but the codecrafters send it in lowercase...?
	switch string(resps[0].Array[0].Data) {
	case "ECHO", "echo":
		message := []byte{}
		for _, v := range resps[0].Array[1:] {
			message = append(message, v.Data...)
		}
		conn.Write(message)
	case "PING", "ping":
		conn.Write([]byte("+PONG\r\n"))
	default:
		panic("not implemented")
	}
}
