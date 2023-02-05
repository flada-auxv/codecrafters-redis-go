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
		name string
		args args
		want []RESP
	}{
		{
			name: "hi",
			args: args{[]byte("*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n")},
			want: []RESP{
				{Type: '*', Count: 2, Data: []byte("$4\r\nECHO\r\n$3\r\nhey\r\n"), Raw: []byte("*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n")},
				{Type: '$', Count: 4, Data: []byte("ECHO"), Raw: []byte("$4\r\nECHO\r\n")},
				{Type: '$', Count: 3, Data: []byte("hey"), Raw: []byte("$3\r\nhey\r\n")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Parse(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				// for _, v := range got {
				// 	t.Errorf("Type: %#v, Count: %#v, Data: %#v, Raw: %#v", string(v.Type), v.Count, string(v.Data), string(v.Raw))
				// }
				t.Errorf("Parse() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
