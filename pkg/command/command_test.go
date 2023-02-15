package command

import (
	"bytes"
	"codecrafters-redis-go/pkg/store"
	"io"
	"net"
	"testing"
	"time"
)

func TestCmdEcho_Run(t *testing.T) {
	t.Run("Repeat the value and write to a conneciton.", func(t *testing.T) {
		mock := &mockConn{
			mock: &bytes.Buffer{},
		}
		c := &CmdEcho{
			cmdCtx: cmdCtx{
				conn:  mock,
				store: store.NewMemoryStore(time.Now),
			},
			opts: &CmdEchoOpts{Value: "hi"},
		}

		if err := c.Run(); err != nil {
			t.Errorf("err: %v", err)
		}
		bytes, err := io.ReadAll(mock)
		if err != nil {
			t.Errorf("err: %v", err)
		}
		if string(bytes) != "$2\r\nhi\r\n" {
			t.Errorf("The written value to conn is not as expected. value: %v", string(bytes))
		}
	})
}

type mockConn struct {
	mock io.ReadWriter
}

func (c *mockConn) Read(b []byte) (n int, err error) {
	return c.mock.Read(b)
}
func (c *mockConn) Write(b []byte) (n int, err error) {
	return c.mock.Write(b)
}
func (c *mockConn) Close() error {
	return nil
}
func (c *mockConn) LocalAddr() net.Addr {
	return nil
}
func (c *mockConn) RemoteAddr() net.Addr {
	return nil
}
func (c *mockConn) SetDeadline(t time.Time) error {
	return nil
}
func (c *mockConn) SetReadDeadline(t time.Time) error {
	return nil
}
func (c *mockConn) SetWriteDeadline(t time.Time) error {
	return nil
}
func NewMockConn(mockIO io.ReadWriter) *mockConn {
	return &mockConn{
		mock: mockIO,
	}
}
