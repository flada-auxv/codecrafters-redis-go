package command

import (
	"codecrafters-redis-go/pkg/resp"
	"errors"
)

type CmdGet struct {
	Cmd
	CmdGetOpts
}
type CmdGetOpts struct {
	Key string
}

func (c *CmdGet) Run() error {
	v, err := c.cmdCtx.store.Get(c.CmdGetOpts.Key)
	if err != nil {
		return err
	}
	if v == "" {
		c.cmdCtx.conn.Write(resp.EncodeNullBulkString())
		return nil
	}
	c.cmdCtx.conn.Write(resp.EncodeBulkString(v))
	return nil
}

type CmdGetFactory struct{}

func (c *CmdGetFactory) CreateCmd(cmdCtx CmdCtx, args []resp.RESP) (*CmdGet, error) {
	if len(args) != 1 {
		return nil, errors.New("ERR invalid argument length for GET")
	}
	if args[0].Type != resp.RESPBulkString {
		return nil, errors.New("ERR invalid argument type for GET")
	}

	return &CmdGet{
		Cmd: Cmd{
			cmdCtx: cmdCtx,
		},
		CmdGetOpts: CmdGetOpts{
			Key: string(args[0].Data),
		},
	}, nil
}
