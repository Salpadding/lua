package vm

import (
	"fmt"
	"os"
	"testing"

	"github.com/Salpadding/lua/types/chunk"
	"github.com/Salpadding/lua/types/value"
	"github.com/Salpadding/lua/types/value/types"
	"github.com/stretchr/testify/assert"
)

func TestNewRegister(t *testing.T) {
	r := NewRegister(0)
	assert.NoError(t, r.Push(value.String("hello world")))
	v, err := r.Pop()
	assert.NoError(t, err)
	fmt.Println(v)
	_, err = r.Pop()
	assert.Error(t, err)
}

func TestIndex(t *testing.T) {
	r := NewRegister(64)
	assert.NoError(t, r.Push(value.String("hello world")))
	idx := r.AbsIndex(-1)
	assert.Equal(t, 0, idx)
	v := r.Get(0)
	_, ok := v.(*value.Nil)
	assert.False(t, ok)
	assert.NoError(t, r.Set(0, value.Float(1.2)))
	v = r.Get(0)
	_, ok = v.(value.Float)
	assert.True(t, ok)
}

func TestStack(t *testing.T) {
	s := &LuaVM{Register: NewRegister(0)}
	if err := s.Push(value.Boolean(true)); err != nil {
		t.Error(err)
	}
	fmt.Println(s)
	if err := s.Push(value.Integer(10)); err != nil {
		t.Error(err)
	}
	fmt.Println(s)
	if err := s.Push(value.GetNil()); err != nil {
		t.Error(err)
	}
	fmt.Println(s)
	if err := s.Push(value.String("hello")); err != nil {
		t.Error(err)
	}
	fmt.Println(s)
	if err := s.Push(value.Integer(-4)); err != nil {
		t.Error(err)
	}
	fmt.Println(s)
	if err := s.Replace(3); err != nil {
		t.Error(err)
	}
	fmt.Println(s)
	if err := s.SetTop(6); err != nil {
		t.Error(err)
	}
	fmt.Println(s)
	if err := s.Remove(-3); err != nil {
		t.Error(err)
	}
	fmt.Println(s)
	if err := s.SetTop(-5); err != nil {
		t.Error(err)
	}
	fmt.Println(s)
}

func TestArithmetic(t *testing.T) {
	s := &LuaVM{Register: NewRegister(256)}
	assert.NoError(t, s.Push(value.Integer(1)))
	assert.NoError(t, s.Push(value.String("2.0")))
	assert.NoError(t, s.Push(value.String("3.0")))
	assert.NoError(t, s.Push(value.Float(4.0)))
	fmt.Println(s)
	assert.NoError(t, s.Arithmetic(types.Add))
	fmt.Println(s)
	assert.NoError(t, s.Arithmetic(types.BitwiseNot))
	fmt.Println(s)
	assert.NoError(t, s.Concat(3))
	fmt.Println(s)
}

func TestBin(t *testing.T) {
	f, err := os.Open("testdata/luac.out")
	assert.NoError(t, err)
	proto, err := chunk.ReadPrototype(f)
	assert.NoError(t, err)
	vm := &LuaVM{
		proto: proto,
		pc:    0,
	}
	assert.NoError(t, vm.execute())
}
