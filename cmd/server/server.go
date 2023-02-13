package main

import (
	"bufio"
	resp "codecrafters-redis-go/pkg/resp"
	store "codecrafters-redis-go/pkg/store"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

var Logger *log.Logger

func main() {
	// TODO: Allow writers to specific files specified in the config
	Logger = log.New(os.Stdout, "codecrafters-redis-go", log.LstdFlags)

	host := flag.String("h", "0.0.0.0", "Server hostname (default: 0.0.0.0)")
	port := flag.String("p", "6379", "Server port (default: 6379)")
	flag.Parse()

	address := *host + ":" + *port

	l, err := net.Listen("tcp", address)
	if err != nil {
		Logger.Fatalf("Failed to listen. address: %v\n", address)
	}

	defer l.Close()

	store := store.NewMemoryStore(time.Now)

	for {
		conn, err := l.Accept()
		if err != nil {
			Logger.Fatalf("Error accepting connection. error: %v\n", err.Error())
		}

		go handleConnection(conn, store)
	}
}

func handleConnection(conn net.Conn, store store.Store) {
	defer conn.Close()

	for {
		resps, err := resp.Parse(bufio.NewReader(conn))
		if err == io.EOF {
			break
		}
		if err != nil {
			Logger.Fatalf("Error while parsing request. error: %v\n", err.Error())
		}

		exec(conn, store, resps)
	}
}

func exec(conn net.Conn, store store.Store, resps []resp.RESP) {
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
		message := ""
		for i, v := range args {
			if i == 0 {
				message = message + string(v.Data)
			} else {
				message = message + " " + string(v.Data)
			}
		}
		conn.Write(resp.EncodeBulkString(message))

	case "GET", "get":
		v, err := store.Get(string(args[0].Data))
		if err != nil {
			conn.Write(resp.EncodeError(fmt.Errorf("ERR something wrong with GET. error: %v", err.Error())))
			return
		}
		if v == "" {
			conn.Write([]byte(fmt.Sprintf("$%v\r\n", -1)))
		} else {
			conn.Write(resp.EncodeBulkString(v))
		}

	case "SET", "set":
		// TODO: just consider PX being passed, for now
		if len(args) <= 2 {
			err := store.Set(string(args[0].Data), string(args[1].Data))
			conn.Write(resp.EncodeError(fmt.Errorf("ERR something wrong with SET. error: %v", err.Error())))
			return
		}

		if string(args[2].Data) != "PX" {
			conn.Write(resp.EncodeError(errors.New("ERR unknown option for SET")))
			return
		}

		milSec, errFromAtoi := strconv.Atoi(string(args[3].Data))
		if errFromAtoi != nil {
			conn.Write(resp.EncodeError(fmt.Errorf("ERR invalid opition value. error: %v", errFromAtoi)))
			return
		}

		err := store.SetWithExpiration(string(args[0].Data), string(args[1].Data), milSec)
		if err != nil {
			conn.Write(resp.EncodeError(fmt.Errorf("ERR something wrong with SET. error: %v", err.Error())))
			return
		}

		conn.Write(resp.EncodeSimpleString("OK"))

	case "PING", "ping":
		conn.Write(resp.EncodeSimpleString("PONG"))

	default:
		conn.Write(resp.EncodeError(errors.New("ERR not implemented command")))
	}
}
