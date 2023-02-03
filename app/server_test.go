package main

import (
	"reflect"
	"testing"
)

func TestTokenize(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want []RESP
	}{
		{ name: "hi", args: args{[]byte("*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n")}, want: []RESP{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Tokenize(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tokenize() = %v, want %v", got, tt.want)
			}
		})
	}
}
