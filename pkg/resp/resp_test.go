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
			name: "Just a simple Integer",
			args: args{bufio.NewReader(strings.NewReader(":1000\r\n"))},
			want: []RESP{
				{Type: ':', Count: -1, Data: []byte("1000")},
			},
			wantErr: false,
		}, {
			name:    "Type is Integer but data is not",
			args:    args{bufio.NewReader(strings.NewReader(":n\r\n"))},
			want:    nil,
			wantErr: true,
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
		{
			name: "Null Bulk String",
			args: args{bufio.NewReader(strings.NewReader("$-1\r\n"))},
			want: []RESP{{
				Type: '$', Count: -1, Data: []byte(""),
			}},
			wantErr: false,
		},
		{
			name: "Not to short Bulk String",
			args: args{bufio.NewReader(strings.NewReader("$52\r\nSo, what do you wanna do, what's your point-of-view?\r\n"))},
			want: []RESP{{
				Type: '$', Count: 52, Data: []byte("So, what do you wanna do, what's your point-of-view?"),
			}},
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

func TestEncodeArray(t *testing.T) {
	type args struct {
		array []string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "hi",
			args: args{[]string{"hi", "hey"}},
			want: []byte("*2\r\n$2\r\nhi\r\n$3\r\nhey\r\n"),
		},
		{
			name: "empty",
			args: args{[]string{}},
			want: []byte("*0\r\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodeArray(tt.args.array); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncodeArray() = %#v, want %#v", string(got), string(tt.want))
			}
		})
	}
}

func TestEncodeBulkString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "hi",
			args: args{"hi"},
			want: []byte("$2\r\nhi\r\n"),
		},
		{
			name: "empty",
			args: args{""},
			want: []byte("$0\r\n\r\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodeBulkString(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncodeBulkString() = %#v, want %#v", string(got), string(tt.want))
			}
		})
	}
}

func TestEncodeNullArray(t *testing.T) {
	tests := []struct {
		name string
		want []byte
	}{
		{
			name: "hi",
			want: []byte("*-1\r\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodeNullArray(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncodeNullArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodeNullBulkString(t *testing.T) {
	tests := []struct {
		name string
		want []byte
	}{
		{
			name: "hi",
			want: []byte("$-1\r\n"),
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodeNullBulkString(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncodeNullBulkString() = %v, want %v", got, tt.want)
			}
		})
	}
}
