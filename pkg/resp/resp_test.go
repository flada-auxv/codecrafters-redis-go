package resp

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

func TestRESP_ToString(t *testing.T) {
	type fields struct {
		Array []RESP
		Count int
		Data  []byte
		Type  byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "BulkString",
			fields: fields{
				Count: -1,
				Data:  []byte("hi hoi hey"),
				Type:  RESPBulkString,
			},
			want: "hi hoi hey",
		},
		{
			name: "SimpleString",
			fields: fields{
				Count: -1,
				Data:  []byte("OK"),
				Type:  RESPSimpleString,
			},
			want: "OK",
		}, {
			name: "Error",
			fields: fields{
				Count: -1,
				Data:  []byte("ERR: something went wrong"),
				Type:  RESPError,
			},
			want: "ERR: something went wrong",
		},
		{
			name: "Array",
			fields: fields{
				Array: []RESP{
					{
						Count: -1,
						Data:  []byte("SET"),
						Type:  RESPBulkString,
					},
					{
						Count: -1,
						Data:  []byte("testKey"),
						Type:  RESPBulkString,
					},
					{
						Count: -1,
						Data:  []byte("testValue"),
						Type:  RESPBulkString,
					},
					{
						Count: -1,
						Data:  []byte("someOpts"),
						Type:  RESPBulkString,
					},
				},
				Count: 2,
				Type:  RESPArray,
			},
			want: "SET testKey testValue someOpts",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RESP{
				Array: tt.fields.Array,
				Count: tt.fields.Count,
				Data:  tt.fields.Data,
				Type:  tt.fields.Type,
			}
			if got := r.ToString(); got != tt.want {
				t.Errorf("RESP.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
