package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strconv"
)

type RESP struct {
	Array []RESP
	Count int
	Data  []byte
	Type  byte
}

const (
	RESPArray        = '*'
	RESPBulkString   = '$'
	RESPError        = '-'
	RESPInteger      = ':'
	RESPSimpleString = '+'
	newLine          = "\r\n"
)

func Parse(s *bufio.Reader) ([]RESP, error) {
	resps := []RESP{}

	rawLine, err := readLine(s)
	if err != nil {
		if err == io.EOF {
			return nil, err
		}
		return nil, errors.New("Invalid RESP format. It does not end in CRLF.")
	}

	switch rawLine[0] {
	case RESPArray:
		count, err := strconv.Atoi(string(rawLine[1]))
		if err != nil {
			return nil, errors.New("Invalid RESP format (Array): Invalid integer.")
		}

		respArray := RESP{
			Count: count,
			Type:  RESPArray,
		}

		for i := 1; i <= count; i++ {
			r, err := Parse(s)
			if err != nil {
				return nil, err
			}
			respArray.Array = append(respArray.Array, r...)
		}
		resps = append(resps, respArray)
	case RESPBulkString:
		count, err := strconv.Atoi(string(rawLine[1]))
		if err != nil {
			return nil, errors.New("Invalid RESP format (BulkString): Invalid integer.")
		}

		contentAndNewLine := make([]byte, count+len(newLine))
		if _, err := io.ReadFull(s, contentAndNewLine); err != nil {
			return nil, errors.New("Invalid RESP format (BulkString): Wrong content size.")
		}
		resps = append(resps, RESP{
			Count: count,
			Data:  contentAndNewLine[:len(contentAndNewLine)-len(newLine)],
			Type:  RESPBulkString,
		})
	case RESPError:
		panic("TODO")
	case RESPInteger:
		panic("TODO")
	case RESPSimpleString:
		resps = append(resps, RESP{
			Count: -1,
			Data:  rawLine[1 : len(rawLine)-len(newLine)],
			Type:  RESPSimpleString,
		})
	default:
		panic("TODO")
	}

	return resps, nil
}

func readLine(s *bufio.Reader) ([]byte, error) {
	bytes := []byte{}

	for {
		b, err := s.ReadBytes('\n')
		if err != nil {
			// does not end in '\n'
			return nil, err
		}

		bytes = append(bytes, b...)

		if bytes[len(bytes)-2] == '\r' {
			break
		}
	}

	return bytes, nil
}

func scanCRLF(b []byte) int {
	return bytes.Index(b, []byte(newLine)) + len(newLine)
}
