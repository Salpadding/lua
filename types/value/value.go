package value

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Salpadding/lua/types/value/types"

	"github.com/Salpadding/lua/common"
)

type Value interface {
	value()
	// String is stringer for debug
	String() string
	Type() types.Type
	ToNumber() (Number, bool)
	ToInteger() (Integer, bool)
	ToFloat() (Float, bool)
	ToString() (string, bool)
	ToBoolean() Boolean
}

var (
	luaNil  = &Nil{}
	luaNone = &None{}
)

func GetNil() *Nil {
	return luaNil
}

func GetNone() *None {
	return luaNone
}

type None struct{}

func (n *None) ToString() (string, bool) {
	return "", false
}

func (n *None) value() {}

func (n *None) String() string {
	return "none"
}

func (n *None) ToNumber() (Number, bool) {
	return nil, false
}

func (n *None) ToInteger() (Integer, bool) {
	return 0, false
}

func (n *None) Type() types.Type {
	return types.None
}

func (n *None) ToBoolean() Boolean {
	return false
}

func (n *None) ToFloat() (Float, bool) {
	return 0, false
}

type Nil struct{}

func (n *Nil) ToString() (string, bool) {
	return "", false
}

func (n *Nil) value() {}

func (n *Nil) String() string {
	return "nil"
}

func (n *Nil) Type() types.Type {
	return types.Nil
}

func (n *Nil) ToNumber() (Number, bool) {
	return nil, false
}

func (n *Nil) ToInteger() (Integer, bool) {
	return 0, false
}

func (n *Nil) ToBoolean() Boolean {
	return false
}

func (n *Nil) ToFloat() (Float, bool) {
	return 0, false
}

type Boolean bool

func (b Boolean) ToString() (string, bool) {
	return "", false
}

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

func (b Boolean) ToNumber() (Number, bool) {
	return nil, false
}

func (b Boolean) ToInteger() (Integer, bool) {
	return 0, false
}

func (b Boolean) ToBoolean() Boolean {
	return b
}

func (b Boolean) ToFloat() (Float, bool) {
	return 0, false
}

type Number interface {
	Value
	number()
}

type Float float64

func (f Float) number() {}

func (f Float) ToString() (string, bool) {
	return f.String(), true
}

func (f Float) value() {}

func (f Float) Type() types.Type {
	return types.Number
}

func (f Float) String() string {
	return strconv.FormatFloat(float64(f), 'f', -1, 64)
}

func (f Float) ToNumber() (Number, bool) {
	return f, true
}

func (f Float) ToInteger() (Integer, bool) {
	// todo: correct?
	i := int64(f)
	return Integer(i), Float(i) == f
}

func (f Float) ToBoolean() Boolean {
	return true
}

func (f Float) ToFloat() (Float, bool) {
	return f, true
}

type Integer int64

func (i Integer) number() {}

func (i Integer) ToString() (string, bool) {
	return i.String(), true
}

func (i Integer) value() {}

func (i Integer) String() string {
	return strconv.FormatInt(int64(i), 10)
}

func (i Integer) Type() types.Type {
	return types.Number
}

func (i Integer) ToNumber() (Number, bool) {
	return i, true
}

func (i Integer) ToInteger() (Integer, bool) {
	return i, true
}

func (i Integer) ToBoolean() Boolean {
	return true
}

func (i Integer) ToFloat() (Float, bool) {
	return Float(i), true
}

type String string

func (s String) ToString() (string, bool) {
	return string(s), true
}

func (s String) value() {}

func (s String) String() string {
	return "\"" + common.Escape(bytes.NewBufferString(string(s))) + "\""
}

func (s String) Type() types.Type {
	return types.String
}

func (s String) ToNumber() (Number, bool) {
	return ParseNumber(string(s))
}

func (s String) ToInteger() (Integer, bool) {
	if i, ok := ParseInteger(string(s)); ok {
		return Integer(i), ok
	}
	if f, ok := ParseFloat(string(s)); ok {
		return Integer(f), float64(Integer(f)) == f
	}
	return 0, false
}

func (s String) ToBoolean() Boolean {
	return true
}

func (s String) ToFloat() (Float, bool) {
	f, ok := ParseFloat(string(s))
	if !ok {
		return 0, ok
	}
	return Float(f), true
}

type Table struct {
	array   *array
	m       map[string]Value
}

func (t *Table) ToNumber() (Number, bool) {
	return nil, false
}

func (t *Table) ToInteger() (Integer, bool) {
	return 0, false
}

func (t *Table) Type() types.Type {
	return types.Table
}

func (t *Table) value() {}

func (t *Table) String() string {
	return fmt.Sprintf("[ %s ]", common.Join(common.ToGeneral(t.array), ", "))
}

func (t *Table) ToString() (string, bool) {
	return "", false
}

func (t *Table) ToBoolean() Boolean {
	return true
}

func (t *Table) ToFloat() (Float, bool) {
	return 0, false
}

func (t *Table) Set(k Value, v Value) error {
	switch x := k.(type) {
	case *Nil, *None:
		return nil
	case Integer:
		return t.array.Set(int(x), v)
	case Float:
		if math.IsNaN(float64(x)) {
			return errors.New("NaN index")
		}
		i, ok := x.ToInteger()
		if ok {
			return t.array.Set(int(i), v)
		}

	}
}

func (t *Table) Get(k Value) Value {
	return nil
}

type Function struct {
}

func (f *Function) ToString() string {
	return ""
}

func (f *Function) ToNumber() (Number, bool) {
	return nil, false
}

func (f *Function) ToInteger() (Integer, bool) {
	return 0, false
}

func (f *Function) value() {}

func (f *Function) String() string { return "function" }

func (f *Function) Type() types.Type {
	return types.Function
}

func (f *Function) ToBoolean() Boolean {
	return true
}

type Thread func()

func (t Thread) ToString() (string, bool) {
	return "", false
}

func (t Thread) ToNumber() (Number, bool) {
	return nil, false
}

func (t Thread) ToInteger() (Integer, bool) {
	return 0, false
}

func (t Thread) value() {}

func (t Thread) String() string { return "thread" }

func (t Thread) Type() types.Type {
	return types.Thread
}

func (t Thread) ToBoolean() Boolean {
	return true
}

type array []Value

func (l *array) Get(idx int) (Value, error) {
	idx = idx - 1
	if idx < 0 || idx >= len(*l) {
		return nil, errors.New("index overflow")
	}
	return (*l)[idx], nil
}

func (l *array) Set(idx int, val Value) error {
	idx = idx - 1
	if idx < 0 {
		return errors.New("index overflow")
	}
	if idx > len(*l){
		return errors.New("index overflow")
	}
	for idx == len(*l) {
		*l = append(*l, GetNil())
	}
	(*l)[idx] = val
	return nil
}

func (l *array) Len() int {
	return len(*l)
}

func (l *array) shrink() {
	if len(*l) == 0 {
		return
	}
	for (*l)[len(*l)-1].Type() == types.Nil {
		*l = (*l)[:len(*l)-1]
	}
}
