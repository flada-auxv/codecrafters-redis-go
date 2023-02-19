package main

import (
	"bufio"
	"codecrafters-redis-go/pkg/command"
	resp "codecrafters-redis-go/pkg/resp"
	store "codecrafters-redis-go/pkg/store"
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

	host := flag.String("h", "127.0.0.1", "Server hostname (default: 127.0.0.1)")
	port := flag.String("p", "6379", "Server port (default: 6379)")
	flag.Parse()

	address := *host + ":" + *port

	l, err := net.Listen("tcp", address)
	if err != nil {
		Logger.Fatalf("Failed to listen. address: %v", address)
	}

	defer l.Close()

	store := store.NewMemoryStore(time.Now)

	for {
		conn, err := l.Accept()
		if err != nil {
			Logger.Fatalf("Error accepting connection. error: %v", err)
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
			Logger.Printf("Error while parsing request. error: %v", err)
			continue
		}

		exec(conn, store, resps)
	}
}

func exec(conn net.Conn, store store.Store, resps []resp.RESP) {
	if resps[0].Type != resp.RESPArray || len(resps) != 1 {
		Logger.Printf("Currently, only an Array is supported.")
		return
	}

	cmdCtx := command.NewCmdCtx(conn, store)
	cmdType := resps[0].Array[0].Data
	args := resps[0].Array[1:]
	c, err := command.GetCmd(cmdCtx, string(cmdType), args)
	if err != nil {
		writeError(err, conn)
		return
	}

	if err := c.Run(); err != nil {
		writeError(err, conn)
		return
	}
}

func writeError(e error, conn net.Conn) {
	conn.Write(resp.EncodeError(e))
}
