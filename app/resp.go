package main

import (
	"bytes"
	"errors"
	"strconv"
)

const (
	newLine = "\r\n"
)

func Parse(b []byte) ([]RESP, error) {
	if bytes.Count(b, []byte(newLine)) == 0 {
		return nil, errors.New("no new line")
	}

	resps := []RESP{}

	for len(b) > 0 {
		eol := scanCRLF(b)
		rawLine := b[0:eol]

		switch b[0] {
		case RESPArray:
			count, err := strconv.Atoi(string(b[1]))
			if err != nil {
				return nil, errors.New("Invalid RESP format.")
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
				return nil, errors.New("Invalid RESP format.")
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

	return resps, nil
}

func scanCRLF(b []byte) int {
	return bytes.Index(b, []byte(newLine)) + len(newLine)
}
