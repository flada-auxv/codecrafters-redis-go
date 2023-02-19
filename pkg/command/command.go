package command

import (
	"codecrafters-redis-go/pkg/resp"
	"codecrafters-redis-go/pkg/store"
	"errors"
	"net"
)

type ICmd interface {
	Run() error
}
type Cmd struct {
	cmdCtx CmdCtx
}
type CmdCtx struct {
	conn  net.Conn
	store store.Store
}
type CmdFactory interface {
	CreateCmd(cmdCtx CmdCtx, args []resp.RESP) (ICmd, error)
}

func GetCmd(cmdCtx CmdCtx, cmdType string, args []resp.RESP) (ICmd, error) {
	switch cmdType {
	case "ECHO", "echo":
		return new(CmdEchoFactory).CreateCmd(cmdCtx, args)

	case "GET", "get":
		return new(CmdGetFactory).CreateCmd(cmdCtx, args)

	case "PING", "ping":
		return new(CmdPingFactory).CreateCmd(cmdCtx, args)

	case "SET", "set":
		return new(CmdSetFactory).CreateCmd(cmdCtx, args)

	default:
		return nil, errors.New("ERR not implemented or unknown command")
	}
}
func NewCmdCtx(c net.Conn, s store.Store) CmdCtx {
	return CmdCtx{
		conn:  c,
		store: s,
	}
}
