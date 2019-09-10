package value

import (
	"bytes"
	"strconv"

	"github.com/Salpadding/lua/common"
)

type Instruction uint32

type Value interface {
	value()
	String() string
}

type Nil string

func (n Nil) value() {}

func (n Nil) String() string {
	return string(n)
}

type Boolean bool

func (b Boolean) value() {}

func (b Boolean) String() string {
	if b {
		return "true"
	}
	return "false"
}

type Number float64

func (n Number) value() {}

func (n Number) String() string {
	return strconv.FormatFloat(float64(n), 'f', -1, 64)
}

type Integer int64

func (i Integer) value() {}

func (i Integer) String() string {
	return strconv.FormatInt(int64(i), 64)
}

type String string

func (s String) value() {}

func (s String) String() string {
	return common.Escape(bytes.NewBufferString(string(s)))
}


