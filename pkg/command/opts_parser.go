package command

import (
	"errors"
	"strconv"
)

type OptParser struct {
	opts   map[string]Opt
	parsed bool
}

func NewOptParser() *OptParser {
	return &OptParser{
		opts: map[string]Opt{},
	}
}

type Opt struct {
	Name  string
	Type  OptType
	Value OptValue
}
type OptType int

const (
	OptTypeInt OptType = iota
	OptTypeString
	OptTypeBool
)

type OptValue interface {
	Set(string) error
}
type BoolValue struct {
	Ptr *bool
}

func (o *BoolValue) Set(v string) error {
	b, err := strconv.ParseBool(v)
	if err != nil {
		return errors.New("TODO")
	}
	*o.Ptr = b
	return nil
}

type IntValue struct {
	Ptr *int
}

func (o *IntValue) Set(v string) error {
	i, err := strconv.Atoi(v)
	if err != nil {
		return errors.New("TODO")
	}
	*o.Ptr = i
	return nil
}

func (c *OptParser) Parse(args []string) error {
	defer func() {
		c.parsed = true
	}()

	if c.parsed {
		return errors.New("already parsed")
	}

	for i := 0; i < len(args); {
		v, ok := c.opts[args[i]]
		if !ok {
			return errors.New("unknown option")
		}
		switch v.Type {
		case OptTypeInt:
			err := v.Value.Set(args[i+1])
			if err != nil {
				return errors.New("TODO")
			}
			i = i + 2

		case OptTypeString:
			return errors.New("unsupported")

		case OptTypeBool:
			err := v.Value.Set("true")
			if err != nil {
				return errors.New("TODO")
			}
			i = i + 1

		default:
			return errors.New("unsupported")
		}
	}

	return nil
}

func (c *OptParser) SetInt(name string) *int {
	p := new(int)
	c.opts[name] = Opt{
		Name: name,
		Type: OptTypeInt,
		Value: &IntValue{
			Ptr: p,
		},
	}
	return p
}

func (c *OptParser) SetBool(name string) *bool {
	p := new(bool)
	c.opts[name] = Opt{
		Name: name,
		Type: OptTypeBool,
		Value: &BoolValue{
			Ptr: p,
		},
	}
	return p
}
