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

		handleConnection(conn)
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

// *2\r\n$4\r\nLLEN\r\n$6\r\nmylist\r\n

// *2\r\n
// $4\r\n
// LLEN\r\n
// $6\r\n
// mylist\r\n

// [RESP{*2}, RESP{$4:LLEN}, RESP{$6:mylist}]

func scanCRLF(b []byte) int {
	return bytes.Index(b, []byte("\r\n")) + 2
}

func Tokenize(b []byte) []RESP {
	lineNum := bytes.Count(b, []byte("\r\n"))
	if lineNum == 0 {
		panic("invalid format")
	}

	bSize := len(string(b))
	cursor := 0
	resps := []RESP{}

	fmt.Printf("initial! buffer: %#v, size: %#v", string(b), bSize)

	for len(b) > 0 {
		fmt.Println(fmt.Sprintf("loop start! %#v", string(b)))
		cursor = scanCRLF(b)
		rawLine := b[0:cursor]

		fmt.Println(fmt.Sprintf("cursor: %#v, rawLine: %#v", cursor, string(rawLine)))

		switch b[0] {
		case RESPArray:
			count, err := strconv.Atoi(string(b[1]))
			if err != nil {
				panic("TODO")
			}
			resps = append(resps, RESP{
				Count: count,
				Data:  b[cursor:], // FIXME: parse each of array elements
				Raw:   b[cursor:],
				Type:  RESPArray,
			})
			fmt.Println(fmt.Sprintf("array! buf: %#v, next: %#v", string(b), string(b[cursor:])))
			b = b[cursor:]
		case RESPBulkString:
			count, err := strconv.Atoi(string(b[1]))
			if err != nil {
				panic("TODO")
			}
			resps = append(resps, RESP{
				Count: count,
				Data:  b[cursor : cursor+count],
				Raw:   b[cursor : cursor+count+2],
				Type:  RESPBulkString,
			})
			b = b[cursor+count+2:]
		case RESPSimpleString:
			resps = append(resps, RESP{
				Count: -1,
				Data:  rawLine[1 : len(rawLine)-2],
				Raw:   rawLine,
				Type:  RESPSimpleString,
			})
			b = b[cursor:]
		default:
			fmt.Printf("bSize: %+v, cursor: %+v, Type: %+v, Data: %+v, Raw: %+v\n", bSize, cursor, string(resps[0].Type), string(resps[0].Data), string(resps[0].Raw))
			panic("TODO")
		}
	}

	fmt.Println(fmt.Printf("resps: %#v", resps))

	return resps
}

func Exec(conn net.Conn, resps []RESP) {
	if resps[0].Type != RESPArray {
		panic("Currently, only array of bulk string is supported")
	}

	arr := resps[0]
	cmd := resps[1]
	augs := resps[2:2+arr.Count-1]

	fmt.Println(fmt.Printf("cmd: %#v", string(cmd.Data)))

	switch string(cmd.Data) {
	case "ECHO":
		message := []byte{}
		for _, v := range augs {
			message = append(message, v.Raw...)
		}
		conn.Write(message)
	case "PING":
		conn.Write([]byte("PONG\r\n"))
	default:
		panic("not implemented")
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

		Exec(conn, Tokenize(bytes.Trim(buf, "\x00")))
	}
}
