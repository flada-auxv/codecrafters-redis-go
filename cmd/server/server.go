package main

import (
	"bufio"
	"codecrafters-redis-go/pkg/command"
	resp "codecrafters-redis-go/pkg/resp"
	store "codecrafters-redis-go/pkg/store"
	"errors"
	"flag"
	"io"
	"log"
	"net"
	"os"
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
	cmdCtx := command.NewCmdCtx(conn, store)

	switch string(cmd.Data) {
	case "ECHO", "echo":
		opts, err := command.NewCmdEchoOpts(args)
		if err != nil {
			writeError(err, conn)
			return
		}
		cmd := command.NewCmdEcho(cmdCtx, opts)
		if err := cmd.Run(); err != nil {
			writeError(err, conn)
			return
		}

	case "GET", "get":
		opts, err := command.NewCmdGetOpts(args)
		if err != nil {
			writeError(err, conn)
			return
		}
		cmd := command.NewCmdGet(cmdCtx, opts)
		if err := cmd.Run(); err != nil {
			writeError(err, conn)
			return
		}

	case "SET", "set":
		opts, err := command.NewCmdSetOpts(args)
		if err != nil {
			writeError(err, conn)
			return
		}
		cmd := command.NewCmdSet(cmdCtx, opts)
		if err := cmd.Run(); err != nil {
			writeError(err, conn)
			return
		}

	case "PING", "ping":
		opts, err := command.NewCmdPingOpts(args)
		if err != nil {
			writeError(err, conn)
			return
		}
		cmd := command.NewCmdPing(cmdCtx, opts)
		if err := cmd.Run(); err != nil {
			writeError(err, conn)
			return
		}

	default:
		conn.Write(resp.EncodeError(errors.New("ERR not implemented command")))
	}
}

func writeError(e error, conn net.Conn) {
	conn.Write(resp.EncodeError(e))
}
