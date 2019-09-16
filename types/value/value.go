package value

import (
	"bytes"
	"fmt"
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
	Compare(Number) int
}

type Float float64

func (n Float) ToString() (string, bool) {
	return n.String(), true
}

func (n Float) value() {}

func (n Float) Type() types.Type {
	return types.Number
}

func (n Float) Compare(Number) int {
	return 0
}

func (n Float) String() string {
	return strconv.FormatFloat(float64(n), 'f', -1, 64)
}

func (n Float) ToNumber() (Number, bool) {
	return n, true
}

func (n Float) ToInteger() (Integer, bool) {
	// todo: correct?
	i := int64(n)
	return Integer(i), Float(i) == n
}

func (n Float) ToBoolean() Boolean {
	return true
}

func (n Float) ToFloat() (Float, bool) {
	return n, true
}

type Integer int64

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

func (i Integer) Compare(Number) int {
	return 0
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
	array   []Value
	m       map[string]Value
	isArray bool
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

func (t *Table) ToString() (string, bool) {
	return "", false
}

func (t *Table) ToBoolean() Boolean {
	return true
}

func(t *Table) ToFloat() (Float, bool){
	return 0, false
}

func (t *Table) Set(k Value, v Value) error {
	return nil
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
