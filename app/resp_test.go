package main

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		b *bufio.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []RESP
		wantErr bool
	}{
		{
			name: "Just a Simple String",
			args: args{bufio.NewReader(strings.NewReader("+OK\r\n"))},
			want: []RESP{
				{Type: '+', Count: -1, Data: []byte("OK")},
			},
			wantErr: false,
		},
		{
			name: "PING command",
			args: args{bufio.NewReader(strings.NewReader("*1\r\n$4\r\nping\r\n"))},
			want: []RESP{
				{Type: '*', Count: 1, Array: []RESP{
					{Type: '$', Count: 4, Data: []byte("ping")},
				},
				},
			},
			wantErr: false,
		},
		{
			name: "ECHO command",
			args: args{bufio.NewReader(strings.NewReader("*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n"))},
			want: []RESP{
				{Type: '*', Count: 2, Array: []RESP{
					{Type: '$', Count: 4, Data: []byte("ECHO")},
					{Type: '$', Count: 3, Data: []byte("hey")},
				},
				},
			},
			wantErr: false,
		},
		{
			name: "nested array",
			args: args{bufio.NewReader(strings.NewReader("*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n+Hello\r\n-World\r\n"))},
			want: []RESP{
				{Type: '*', Count: 2, Array: []RESP{
					{Type: '*', Count: 3, Array: []RESP{
						{Type: ':', Count: -1, Data: []byte("1")},
						{Type: ':', Count: -1, Data: []byte("2")},
						{Type: ':', Count: -1, Data: []byte("3")},
					}},
					{Type: '*', Count: 2, Array: []RESP{
						{Type: '+', Count: -1, Data: []byte("Hello")},
						{Type: '-', Count: -1, Data: []byte("World")},
					}},
				}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
