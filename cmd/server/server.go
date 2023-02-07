package main

import (
	"bufio"
	resp "codecrafters-redis-go/pkg/resp"
	store "codecrafters-redis-go/pkg/store"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	defer l.Close()

	store := store.NewStore(time.Now)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn, store)
	}
}

func handleConnection(conn net.Conn, store *store.Store) {
	defer conn.Close()

	for {
		resps, err := resp.Parse(bufio.NewReader(conn))
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error while parsing request", err.Error())
			os.Exit(1)
		}

		exec(conn, store, resps)
	}
}

func exec(conn net.Conn, store *store.Store, resps []resp.RESP) {
	if resps[0].Type != resp.RESPArray {
		panic("Currently, only array of bulk string is supported")
	}

	// TODO: should be execed according to the type of first RESP
	// TODO: The redis command group seems to be case insensitive and uses uppercase, but the codecrafters send it in lowercase...?

	respArr := resps[0]
	cmd := respArr.Array[0]
	args := respArr.Array[1:]

	switch string(cmd.Data) {
	case "ECHO", "echo":
		message := []byte{}
		for _, v := range args {
			message = append(message, v.Data...)
		}
		conn.Write([]byte(fmt.Sprintf("$%v\r\n%v\r\n", len(message), string(message))))
	case "GET", "get":
		v, err := store.Get(string(args[0].Data))
		if err != nil {
			conn.Write([]byte(fmt.Sprintf("-ERR something wrong with GET. error: %#v", err)))
			return
		}
		if v == "" {
			conn.Write([]byte(fmt.Sprintf("$%v\r\n", -1)))
		} else {
			conn.Write([]byte(fmt.Sprintf("$%v\r\n%v\r\n", len(v), v)))
		}
	case "SET", "set":
		// TODO: just consider PX being passed, for now
		var err error
		if len(args) <= 2 {
			err = store.Set(string(args[0].Data), string(args[1].Data))
		} else {
			if string(args[2].Data) != "PX" {
				err = errors.New("ERROR: Unknown option")
			}
			milSec, errFromAtoi := strconv.Atoi(string(args[3].Data))
			if errFromAtoi != nil {
				err = errFromAtoi
			}
			err = store.SetWithExpiration(string(args[0].Data), string(args[1].Data), milSec)
		}
		if err != nil {
			conn.Write([]byte(fmt.Sprintf("-ERR something wrong with SET. error: %#v", err)))
			return
		}
		conn.Write([]byte("+OK\r\n"))
	case "PING", "ping":
		conn.Write([]byte("+PONG\r\n"))
	default:
		panic("not implemented")
	}
}