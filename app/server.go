package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
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
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error while reading from connection", err.Error())
			os.Exit(1)
		}

		Exec(conn, Parse(bytes.Trim(buf, "\x00")))
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

func scanCRLF(b []byte) int {
	return bytes.Index(b, []byte("\r\n")) + 2
}

func Parse(b []byte) []RESP {

	lineNum := bytes.Count(b, []byte("\r\n"))
	if lineNum == 0 {
		panic("invalid format")
	}

	resps := []RESP{}

	for len(b) > 0 {
		eol := scanCRLF(b)
		rawLine := b[0:eol]

		switch b[0] {
		case RESPArray:
			count, err := strconv.Atoi(string(b[1]))
			if err != nil {
				panic("TODO")
			}
			resps = append(resps, RESP{
				Count: count,
				Data:  b[eol:], // FIXME: parse each of array elements
				Raw:   b,
				Type:  RESPArray,
			})
			b = b[eol:]
		case RESPBulkString:
			count, err := strconv.Atoi(string(b[1]))
			if err != nil {
				panic("TODO")
			}
			resps = append(resps, RESP{
				Count: count,
				Data:  b[eol : eol+count],
				Raw:   b[0 : eol+count+2],
				Type:  RESPBulkString,
			})
			b = b[eol+count+2:]
		case RESPError:
			panic("TODO")
		case RESPInteger:
			panic("TODO")
		case RESPSimpleString:
			resps = append(resps, RESP{
				Count: -1,
				Data:  rawLine[1 : len(rawLine)-2],
				Raw:   rawLine,
				Type:  RESPSimpleString,
			})
			b = b[eol:]
		default:
			panic("TODO")
		}
	}

	return resps
}

func Exec(conn net.Conn, resps []RESP) {
	if resps[0].Type != RESPArray {
		panic("Currently, only array of bulk string is supported")
	}

	// TODO: should be execed according to the type of first RESP
	arr := resps[0]
	cmd := resps[1]
	augs := resps[2:2+arr.Count-1]

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
