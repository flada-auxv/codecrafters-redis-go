package command

import (
	"codecrafters-redis-go/pkg/resp"
	"errors"
)

type CmdEcho struct {
	Cmd
	CmdEchoOpts
}
type CmdEchoOpts struct {
	Value string
}

func (c *CmdEcho) Run() error {
	c.cmdCtx.conn.Write(resp.EncodeBulkString(c.CmdEchoOpts.Value))
	return nil
}

type CmdEchoFactory struct{}

func (c *CmdEchoFactory) CreateCmd(cmdCtx CmdCtx, args []resp.RESP) (*CmdEcho, error) {
	if len(args) != 1 {
		return nil, errors.New("ERR invalid argument length for ECHO")
	}
	if args[0].Type != resp.RESPBulkString {
		return nil, errors.New("ERR invalid argument type for ECHO")
	}

	return &CmdEcho{
		Cmd: Cmd{
			cmdCtx: cmdCtx,
		},
		CmdEchoOpts: CmdEchoOpts{
			Value: string(args[0].Data),
		},
	}, nil
}
