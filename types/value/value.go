package value

import (
	"bytes"
	"fmt"
	"github.com/Salpadding/lua/types/value/types"
	"strconv"
	"strings"

	"github.com/Salpadding/lua/common"
)

type Instruction uint32

type Value interface {
	value()
	String() string
	Type() types.Type
}

type None string

func (n None) value() {}

func (n None) String() string {
	return string(n)
}

func (n None) Type() types.Type {
	return types.None
}

type Nil string

func (n Nil) value() {}

func (n Nil) String() string {
	return string(n)
}

func (n Nil) Type() types.Type {
	return types.Nil
}

type Boolean bool

func (b Boolean) value() {}

func (b Boolean) Type() types.Type {
	return types.Boolean
}

func (b Boolean) String() string {
	if b {
		return "true"
	}
	return "false"
}

type Number float64

func (n Number) value() {}

func (n Number) Type() types.Type {
	return types.Number
}

func (n Number) String() string {
	return strconv.FormatFloat(float64(n), 'f', -1, 64)
}

type Integer int64

func (i Integer) value() {}

func (i Integer) String() string {
	return strconv.FormatInt(int64(i), 64)
}

func (i Integer) Type() types.Type {
	return types.Number
}

type String string

func (s String) value() {}

func (s String) String() string {
	return common.Escape(bytes.NewBufferString(string(s)))
}

func (s String) Type() types.Type {
	return types.String
}

type Table struct {
	array   []Value
	m       map[string]Value
	isArray bool
}

func (t *Table) Type() types.Type {
	return types.Table
}

func (t *Table) value() {}

func (t *Table) String() string {
	if t.isArray {
		return fmt.Sprintf("[ %s ]", common.Join(common.ToGeneral(t.array), ", "))
	}
	res := make([]string, len(t.m))
	i := 0
	for k, v := range t.m {
		res[i] = fmt.Sprintf("%s = %s", k, v)
		i++
	}
	return fmt.Sprintf("{ %s }", strings.Join(res, ", "))
}

type Function struct {
}

func (f *Function) value() {}

func (f *Function) String() string { return "function" }

func (f *Function) Type() types.Type {
	return types.Function
}

type Thread func()

func (t Thread) value() {}

func (t Thread) String() string { return "thread" }

func (t Thread) Type() types.Type {
	return types.Thread
}
