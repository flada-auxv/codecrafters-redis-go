package main

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []RESP
		wantErr bool
	}{
		{
			name: "Just a Simple String",
			args: args{[]byte("+OK\r\n")},
			want: []RESP{
				{Type: '+', Count: -1, Data: []byte("OK"), Raw: []byte("+OK\r\n")},
			},
			wantErr: false,
		},
		{
			name: "ECHO command",
			args: args{[]byte("*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n")},
			want: []RESP{
				{Type: '*', Count: 2, Data: []byte("$4\r\nECHO\r\n$3\r\nhey\r\n"), Raw: []byte("*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n")},
				{Type: '$', Count: 4, Data: []byte("ECHO"), Raw: []byte("$4\r\nECHO\r\n")},
				{Type: '$', Count: 3, Data: []byte("hey"), Raw: []byte("$3\r\nhey\r\n")},
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
