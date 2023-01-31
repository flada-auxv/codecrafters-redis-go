package main

import (
	"bufio"
	"fmt"

	// Uncomment this block to pass the first stage
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
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	for {
		s := bufio.NewScanner(conn)
		s.Scan()
		if err = s.Err(); err != nil {
			fmt.Println("Error while reading from connection", err.Error())
			os.Exit(1)
		}
		fmt.Println("Received message:", s.Text())

		// FIXME: Just sending "PONG" back in RESP format
		w := bufio.NewWriter(conn)
		fmt.Fprint(w, "+PONG\r\n")
		err = w.Flush()
		if err != nil {
			fmt.Println("Error while writing to connection")
			os.Exit(1)
		}
	}
}
