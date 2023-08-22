package coder

import (
	"fmt"
	"istomyang.github.com/like-iam/component-base/base"
	"istomyang.github.com/like-iam/component-base/json"
)

type CodeOpt string

// Those are valid decoders assign to Config.Decode.
const (
	CodeJson    CodeOpt = "json"
	CodeMsgPack CodeOpt = "msgpack"
)

// coders considers different service has different coder.
var coders map[string]*coderImpl

// Register registers a coder for specific name.
// If existed, ignore it.
func Register(name string, opt CodeOpt) error {
	if _, exist := coders[name]; exist {
		return fmt.Errorf("name %s existed", name)
	}
	coders[name] = &coderImpl{coder: opt}
	return nil
}

// Get gets base.Coder associated with name.
// if name key doesn't existed, return nil.
func Get(name string) base.Coder {
	if c, exist := coders[name]; exist {
		return c
	} else {
		return nil
	}
}

type coderImpl struct {
	coder CodeOpt
}

func (c *coderImpl) Decode(data []byte, v any) error {
	switch c.coder {
	case CodeJson:
		return json.Unmarshal(data, v)
	case CodeMsgPack:
		return json.MPUnmarshal(data, v)
	default:
		return fmt.Errorf("you must assign a valid decoder, got: %s", c.coder)
	}
}

func (c *coderImpl) Encode(v any) ([]byte, error) {
	switch c.coder {
	case CodeJson:
		return json.Marshal(v)
	case CodeMsgPack:
		return json.MPMarshal(v)
	default:
		return nil, fmt.Errorf("you must assign a valid encoder, got: %s", c.coder)
	}
}
