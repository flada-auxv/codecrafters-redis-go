package main

import (
	"codecrafters-redis-go/pkg/resp"
	"testing"
)

func Test_toReadable(t *testing.T) {
	type args struct {
		r resp.RESP
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "an Array",
			args: args{
				resp.RESP{
					Array: []resp.RESP{
						{
							Data: []byte("one"),
							Type: resp.RESPBulkString,
						},
						{
							Data: []byte("two"),
							Type: resp.RESPBulkString,
						},
						{
							Data: []byte("3"),
							Type: resp.RESPInteger,
						},
					},
					Type: resp.RESPArray,
				},
			},
			want: "1) \"one\"\n2) \"two\"\n3) (integer) 3",
		},
		{
			name: "an empty Array",
			args: args{
				resp.RESP{
					Array: []resp.RESP{},
					Type:  resp.RESPArray,
				},
			},
			want: "(empty array)",
		},
		// TODO: nested array
		{
			name: "a BulkString",
			args: args{
				resp.RESP{
					Data: []byte("So I start a revolution from my bed. 'Cause you said the brains I had went to my head."),
					Type: resp.RESPBulkString,
				},
			},
			want: "\"So I start a revolution from my bed. 'Cause you said the brains I had went to my head.\"",
		},
		{
			name: "an empty BulkString",
			args: args{
				resp.RESP{
					Data: []byte(""),
					Type: resp.RESPBulkString,
				},
			},
			want: "(nil)",
		},

		{
			name: "an Error",
			args: args{
				resp.RESP{
					Data: []byte("ERR something went wrong"),
					Type: resp.RESPError,
				},
			},
			want: "(error) ERR something went wrong",
		},
		{
			name: "a Integer",
			args: args{
				resp.RESP{
					Data: []byte("1"),
					Type: resp.RESPInteger,
				},
			},
			want: "(integer) 1",
		},
		{
			name: "a SimpleString",
			args: args{
				resp.RESP{
					Data: []byte("OK"),
					Type: resp.RESPSimpleString,
				},
			},
			want: "OK", // NOTE: This should not be "\"OK\""
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toReadable(tt.args.r); got != tt.want {
				t.Errorf("toReadable() = %v, want %v", got, tt.want)
			}
		})
	}
}
