package command

import (
	"bytes"
	"codecrafters-redis-go/pkg/resp"
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
		ctx := NewCmdCtx(mock, store.NewMemoryStore(time.Now))
		c, err := GetCmd(ctx, "ECHO", []resp.RESP{{Data: []byte("hi"), Type: resp.RESPBulkString}})
		if err != nil {
			t.Errorf("err: %v", err)
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

func TestCmdPing_Run(t *testing.T) {
	t.Run("It returns the PONG when no arguments are passed.", func(t *testing.T) {
		mock := &mockConn{
			mock: &bytes.Buffer{},
		}
		ctx := NewCmdCtx(mock, store.NewMemoryStore(time.Now))
		c, err := GetCmd(ctx, "PING", []resp.RESP{})
		if err != nil {
			t.Errorf("err: %v", err)
		}

		if err := c.Run(); err != nil {
			t.Errorf("err: %v", err)
		}
		bytes, err := io.ReadAll(mock)
		if err != nil {
			t.Errorf("err: %v", err)
		}
		expected := "+PONG\r\n"
		if string(bytes) != expected {
			t.Errorf("The written value to conn is not as expected. got: %v, expected: %v", string(bytes), expected)
		}
	})

	t.Run("It returns the value when an argument is passed like ECHO.", func(t *testing.T) {
		mock := &mockConn{
			mock: &bytes.Buffer{},
		}
		ctx := NewCmdCtx(mock, store.NewMemoryStore(time.Now))
		c, err := GetCmd(ctx, "PING", []resp.RESP{{Data: []byte("hi there!"), Type: resp.RESPBulkString}})
		if err != nil {
			t.Errorf("err: %v", err)
		}

		if err := c.Run(); err != nil {
			t.Errorf("err: %v", err)
		}
		bytes, err := io.ReadAll(mock)
		if err != nil {
			t.Errorf("err: %v", err)
		}
		expected := "$9\r\nhi there!\r\n"
		if string(bytes) != expected {
			t.Errorf("The written value to conn is not as expected. got: %v, expected: %v", string(bytes), expected)
		}
	})
}

func TestCmdGet_Run(t *testing.T) {
	t.Run("It gets the value from store.", func(t *testing.T) {
		mock := &mockConn{
			mock: &bytes.Buffer{},
		}
		store := store.NewMemoryStore(time.Now)
		store.Set("testKey", "testValue")
		ctx := NewCmdCtx(mock, store)
		c, err := GetCmd(ctx, "GET", []resp.RESP{{Data: []byte("testKey"), Type: resp.RESPBulkString}})
		if err != nil {
			t.Errorf("err: %v", err)
		}

		if err := c.Run(); err != nil {
			t.Errorf("err: %v", err)
		}
		bytes, err := io.ReadAll(mock)
		if err != nil {
			t.Errorf("err: %v", err)
		}
		expected := "$9\r\ntestValue\r\n"
		if string(bytes) != expected {
			t.Errorf("The written value to conn is not as expected. got: %v, expected: %v", string(bytes), expected)
		}
	})

	t.Run("It gets the null bulk string when the key does not exist", func(t *testing.T) {
		mock := &mockConn{
			mock: &bytes.Buffer{},
		}
		store := store.NewMemoryStore(time.Now)
		store.Set("testKey", "testValue")
		ctx := NewCmdCtx(mock, store)
		c, err := GetCmd(ctx, "GET", []resp.RESP{{Data: []byte("testKeyDoesNotExist"), Type: resp.RESPBulkString}})
		if err != nil {
			t.Errorf("err: %v", err)
		}

		if err := c.Run(); err != nil {
			t.Errorf("err: %v", err)
		}
		bytes, err := io.ReadAll(mock)
		if err != nil {
			t.Errorf("err: %v", err)
		}
		expected := "$-1\r\n"
		if string(bytes) != expected {
			t.Errorf("The written value to conn is not as expected. got: %v, exptected: %v", string(bytes), expected)
		}
	})
}

func TestCmdSet_Run(t *testing.T) {
	t.Run("It set the value to the key with store", func(t *testing.T) {
		mock := &mockConn{
			mock: &bytes.Buffer{},
		}
		store := store.NewMemoryStore(time.Now)
		store.Set("testKey", "testValue")
		ctx := NewCmdCtx(mock, store)
		c, err := GetCmd(ctx, "SET", []resp.RESP{
			{Data: []byte("testKey"), Type: resp.RESPBulkString},
			{Data: []byte("hi there!"), Type: resp.RESPBulkString},
		})
		if err != nil {
			t.Errorf("err: %v", err)
		}

		if err := c.Run(); err != nil {
			t.Errorf("err: %v", err)
		}
		bytes, err := io.ReadAll(mock)
		if err != nil {
			t.Errorf("err: %v", err)
		}
		expectedToWrite := "+OK\r\n"
		if string(bytes) != expectedToWrite {
			t.Errorf("The written value to conn is not as expected. got: %v, expected: %v", string(bytes), expectedToWrite)
		}
		got, err := store.Get("testKey")
		if err != nil {
			t.Errorf("err: %v", err)
		}
		expectedToStore := "hi there!"
		if got != expectedToStore {
			t.Errorf("The written value to store is not as expected. got: %v, expected: %v", string(got), expectedToStore)
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
