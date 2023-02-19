package command

import (
	"codecrafters-redis-go/pkg/resp"
	"errors"
)

type CmdPing struct {
	Cmd
	CmdPingOpts
}
type CmdPingOpts struct {
	Value string
}

func (c *CmdPing) Run() error {
	if c.CmdPingOpts.Value == "" {
		c.cmdCtx.conn.Write(resp.EncodeSimpleString("PONG"))
		return nil
	}

	c.cmdCtx.conn.Write(resp.EncodeBulkString(c.CmdPingOpts.Value))
	return nil
}

type CmdPingFactory struct{}

func (c *CmdPingFactory) CreateCmd(cmdCtx CmdCtx, args []resp.RESP) (*CmdPing, error) {
	if len(args) == 0 {
		return &CmdPing{
			Cmd: Cmd{
				cmdCtx: cmdCtx,
			},
			CmdPingOpts: CmdPingOpts{
				Value: "",
			},
		}, nil
	}

	if len(args) != 1 {
		return nil, errors.New("ERR invalid argument length for PING")
	}
	if args[0].Type != resp.RESPBulkString {
		return nil, errors.New("ERR invalid argument type for PING")
	}
	return &CmdPing{
		Cmd: Cmd{
			cmdCtx: cmdCtx,
		},
		CmdPingOpts: CmdPingOpts{
			Value: string(args[0].Data),
		},
	}, nil
}
