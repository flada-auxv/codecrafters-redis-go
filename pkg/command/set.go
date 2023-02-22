package command

import (
	"codecrafters-redis-go/pkg/resp"
	"errors"
	"strconv"
)

type CmdSet struct {
	Cmd
	CmdSetOpts
}
type CmdSetOpts struct {
	Key        string
	Value      string
	Expiration int
}

func (c *CmdSet) Run() error {
	if c.CmdSetOpts.Expiration == 0 {
		err := c.cmdCtx.store.Set(c.CmdSetOpts.Key, c.CmdSetOpts.Value)
		if err != nil {
			return err
		}
		c.cmdCtx.conn.Write(resp.EncodeSimpleString("OK"))
		return nil
	}

	err := c.cmdCtx.store.SetWithExpiration(c.CmdSetOpts.Key, c.CmdSetOpts.Value, c.CmdSetOpts.Expiration)
	if err != nil {
		return err
	}
	c.cmdCtx.conn.Write(resp.EncodeSimpleString("OK"))
	return nil
}

type CmdSetFactory struct{}

func (c *CmdSetFactory) CreateCmd(cmdCtx CmdCtx, args []resp.RESP) (*CmdSet, error) {
	if len(args) < 2 {
		return nil, errors.New("ERR invalid argument length for SET")
	}
	key := args[0]
	value := args[1]
	respOpts := args[2:]

	if key.Type != resp.RESPBulkString || value.Type != resp.RESPBulkString {
		return nil, errors.New("ERR invalid argument type for SET")
	}
	opts := CmdSetOpts{
		Key:   string(key.Data),
		Value: string(value.Data),
	}

	for i := 0; i < len(respOpts); {
		switch string(respOpts[i].Data) {
		case "PX":
			v := respOpts[i+1]
			if v.Type != resp.RESPInteger {
				return nil, errors.New("ERR invalid argument type for SET")
			}
			px, err := strconv.Atoi(string(v.Data))
			if err != nil {
				return nil, errors.New("ERR invalid argument type for SET")
			}
			opts.Expiration = px
			i = i + 2
		case "EX":
			v := respOpts[i+1]
			if v.Type != resp.RESPInteger {
				return nil, errors.New("ERR invalid argument type for SET")
			}
			ex, err := strconv.Atoi(string(v.Data))
			if err != nil {
				return nil, errors.New("ERR invalid argument type for SET")
			}
			opts.Expiration = ex * 1000
			i = i + 2
		default:
			return nil, errors.New("ERR invalid argument for SET")
		}
	}

	return &CmdSet{
		Cmd: Cmd{
			cmdCtx: cmdCtx,
		},
		CmdSetOpts: opts,
	}, nil
}
