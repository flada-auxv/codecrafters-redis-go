package main

import (
	"bufio"
	"fmt"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
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

		s := bufio.NewScanner(conn)
		s.Scan()
		if err = s.Err(); err != nil {
			fmt.Println("Error while reading from connection", err.Error())
			os.Exit(1)
		}

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
