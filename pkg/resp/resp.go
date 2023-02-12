package resp

import (
	"bufio"
	"errors"
	"fmt"
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
		return nil, errors.New("ERR Invalid RESP format. It does not end in CRLF")
	}

	switch rawLine[0] {
	case RESPArray:
		count, err := strconv.Atoi(string(rawLine[1]))
		if err != nil {
			return nil, errors.New("ERR Invalid RESP format (Array): Invalid integer")
		}

		respArray := RESP{
			Count: count,
			Type:  RESPArray,
		}

		for i := 0; i < count; i++ {
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
			return nil, errors.New("ERR Invalid RESP format (BulkString): Invalid integer")
		}

		contentAndNewLine := make([]byte, count+len(newLine))
		if _, err := io.ReadFull(s, contentAndNewLine); err != nil {
			return nil, errors.New("ERR Invalid RESP format (BulkString): Wrong content size")
		}
		resps = append(resps, RESP{
			Count: count,
			Data:  contentAndNewLine[:len(contentAndNewLine)-len(newLine)],
			Type:  RESPBulkString,
		})

	case RESPError:
		resps = append(resps, RESP{
			Count: -1,
			Data:  rawLine[1 : len(rawLine)-len(newLine)],
			Type:  RESPError,
		})

	case RESPInteger:
		resps = append(resps, RESP{
			Count: -1,
			Data:  rawLine[1 : len(rawLine)-len(newLine)],
			Type:  RESPInteger,
		})

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

func EncodeArray(array []string) []byte {
	s := []byte(fmt.Sprint(string(RESPArray), len(array), newLine))
	for _, v := range array {
		// TODO: integer...?
		s = append(s, EncodeBulkString(v)...)
	}
	return s
}

func EncodeNullArray() []byte {
	return []byte(fmt.Sprint(string(RESPArray), -1, newLine))
}

func EncodeBulkString(s string) []byte {
	return []byte(fmt.Sprint(string(RESPBulkString), len(s), newLine, s, newLine))
}

func EncodeNullBulkString() []byte {
	return []byte(fmt.Sprint(string(RESPBulkString), -1, newLine))
}

func EncodeError(e error) []byte {
	return []byte(fmt.Sprint(string(RESPError), e.Error(), newLine))
}

func EncodeSimpleString(s string) []byte {
	return []byte(fmt.Sprint(string(RESPSimpleString), s, newLine))
}

func EncodeInteger(i int) []byte {
	return []byte(fmt.Sprint(string(RESPInteger), i, newLine))
}
