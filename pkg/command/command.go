package command

import (
	"codecrafters-redis-go/pkg/resp"
	"codecrafters-redis-go/pkg/store"
	"net"
)

type command interface {
	Run() error
}

type cmdCtx struct {
	conn  net.Conn
	store store.Store
}

func NewCmdCtx(c net.Conn, s store.Store) cmdCtx {
	return cmdCtx{
		conn:  c,
		store: s,
	}
}

type CmdPing struct {
	cmdCtx
	opts CmdPingOpts
}
type CmdPingOpts struct {
	value string
}

func NewCmdPing(cmdCtx cmdCtx, opts CmdPingOpts) *CmdPing {
	return &CmdPing{
		cmdCtx: cmdCtx,
		opts:   opts,
	}
}
func (c CmdPing) Run() error {
	if c.opts.value == "" {
		c.cmdCtx.conn.Write(resp.EncodeSimpleString("PONG"))
	}

	c.cmdCtx.conn.Write(resp.EncodeBulkString(c.opts.value))

	return nil
}

type CmdEcho struct {
	cmdCtx
	opts CmdEchoOpts
}
type CmdEchoOpts struct {
	Value string
}

func NewCmdEcho(cmdCtx cmdCtx, opts CmdEchoOpts) *CmdEcho {
	return &CmdEcho{
		cmdCtx: cmdCtx,
		opts:   opts,
	}
}
func (c CmdEcho) Run() error {
	c.cmdCtx.conn.Write(resp.EncodeBulkString(c.opts.Value))
	return nil
}

type CmdGet struct {
	cmdCtx
	opts CmdGetOpts
}
type CmdGetOpts struct {
	Key string
}

func NewCmdGet(cmdCtx cmdCtx, opts CmdGetOpts) *CmdGet {
	return &CmdGet{
		cmdCtx: cmdCtx,
		opts:   opts,
	}
}
func (c CmdGet) Run() (string, error) {
	v, err := c.cmdCtx.store.Get(c.opts.Key)
	if err != nil {
		return "", err
	}
	return v, nil
}

type CmdSet struct {
	cmdCtx
	opts CmdSetOpts
}
type CmdSetOpts struct {
	Key        string
	Value      string
	Expiration int
}

func NewCmdSet(cmdCtx cmdCtx, opts CmdSetOpts) *CmdSet {
	return &CmdSet{
		cmdCtx: cmdCtx,
		opts:   opts,
	}
}
func (c CmdSet) Run() error {
	if c.opts.Expiration == 0 {
		err := c.cmdCtx.store.Set(c.opts.Key, c.opts.Value)
		if err != nil {
			return err
		}
	}

	err := c.cmdCtx.store.SetWithExpiration(c.opts.Key, c.opts.Value, c.opts.Expiration)
	if err != nil {
		return err
	}

	return nil
}
