package types

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/Salpadding/lua/common"
	"github.com/Salpadding/lua/types/value"
)

type Value interface {
	value()
	// String is stringer for debug
	String() string
	Type() value.Type
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

func (n *None) Type() value.Type {
	return value.None
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

func (n *Nil) Type() value.Type {
	return value.Nil
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

func (b Boolean) Type() value.Type {
	return value.Boolean
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

func (f Float) Type() value.Type {
	return value.Number
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

func (i Integer) Type() value.Type {
	return value.Number
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

func (s String) Type() value.Type {
	return value.String
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
	array *array
	m     map[Value]Value
}

func NewTable() *Table {
	return &Table{
		array: &array{},
		m:     map[Value]Value{},
	}
}

func (t *Table) ToNumber() (Number, bool) {
	return nil, false
}

func (t *Table) ToInteger() (Integer, bool) {
	return 0, false
}

func (t *Table) Type() value.Type {
	return value.Table
}

func (t *Table) value() {}

func (t *Table) String() string {
	kvs := make([]string, len(t.m))
	i := 0
	for k, v := range t.m {
		kvs[i] = fmt.Sprintf("%s=%s", k, v)
		i++
	}

	return fmt.Sprintf("{ %s }", common.Join(
		append(
			common.ToGeneral(*t.array),
			common.ToGeneral(kvs)...
		), ",",
	))
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
		if int(x) <= t.array.Len()+1 {
			t.array.Set(int(x), v)
			t.expand()
			return nil
		}
		t.m[k] = v
		return nil
	case Float:
		if math.IsNaN(float64(x)) {
			return errors.New("NaN index")
		}
		i, ok := x.ToInteger()
		if ok {
			return t.Set(i, v)
		}
		t.m[k] = v
		return nil
	default:
		if v == nil || v.Type() == value.Nil {
			return nil
		}
		t.m[k] = v
		return nil
	}
}

func (t *Table) expand() {
	idx := t.array.Len() + 1
	for {
		val, ok := t.m[Integer(idx)]
		if !ok || val.Type() == value.Nil {
			break
		}
		delete(t.m, val)
		t.array.Set(idx, val)
		idx++
	}
}

func (t *Table) Len() int {
	return t.array.Len()
}
func (t *Table) Get(k Value) (Value, error) {
	v, ok := t.m[k]
	switch x := k.(type) {
	case *Nil, *None:
		return GetNil(), nil
	case Integer:
		if int(x-1) < t.array.Len() {
			return t.array.Get(int(x))
		}
	case Float:
		if math.IsNaN(float64(x)) {
			return nil, errors.New("NaN index")
		}
		i, ok := x.ToInteger()
		if ok {
			return t.Get(i)
		}
	}
	if ok {
		return v, nil
	}
	return GetNil(), nil
}


type Native func(args ...Value) ([]Value, error)

func (n Native) value() {
}

func (n Native) String() string {
	return "function () \n native code \n end"
}

func (n Native) Type() value.Type {
	return value.Function
}

func (n Native) ToNumber() (Number, bool) {
	return nil, false
}

func (n Native) ToInteger() (Integer, bool) {
	return 0, false
}

func (n Native) ToFloat() (Float, bool) {
	return 0, false
}

func (n Native) ToString() (string, bool) {
	return "", false
}

func (n Native) ToBoolean() Boolean {
	return true
}

type Function struct {
	*Prototype
	UpValues []*ValuePointer
}

func (f *Function) ToFloat() (Float, bool) {
	return 0, false
}

func (f *Function) ToString() (string, bool) {
	return "", false
}

func (f *Function) ToNumber() (Number, bool) {
	return nil, false
}

func (f *Function) ToInteger() (Integer, bool) {
	return 0, false
}

func (f *Function) value() {}

func (f *Function) String() string { return "function" }

func (f *Function) Type() value.Type {
	return value.Function
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

func (t Thread) Type() value.Type {
	return value.Thread
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

func (l *array) Set(idx int, val Value) {
	idx = idx - 1
	for idx == len(*l) {
		*l = append(*l, GetNil())
	}
	(*l)[idx] = val
	if val.Type() == value.Nil {
		l.shrink()
	}
}

func (l *array) Len() int {
	return len(*l)
}

func (l *array) shrink() {
	if len(*l) == 0 {
		return
	}
	for (*l)[len(*l)-1].Type() == value.Nil {
		*l = (*l)[:len(*l)-1]
	}
}

type ValuePointer struct{
	Value
}

func(v *ValuePointer) value() error{
	return nil
}

