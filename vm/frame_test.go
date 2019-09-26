package vm

import (
	"fmt"
	"os"
	"testing"

	"github.com/Salpadding/lua/types"
	"github.com/stretchr/testify/assert"
)

func TestNewRegister(t *testing.T) {
	r := NewRegister(0)
	assert.NoError(t, r.Push(types.String("hello world")))
	v, err := r.Pop()
	assert.NoError(t, err)
	fmt.Println(v)
	_, err = r.Pop()
	assert.Error(t, err)
}

func TestIndex(t *testing.T) {
	r := NewRegister(64)
	assert.NoError(t, r.Push(types.String("hello world")))
	idx := r.AbsIndex(-1)
	assert.Equal(t, 0, idx)
	v := r.Get(0)
	_, ok := v.(*types.Nil)
	assert.False(t, ok)
	assert.NoError(t, r.Set(0, types.Float(1.2)))
	v = r.Get(0)
	_, ok = v.(types.Float)
	assert.True(t, ok)
}

func TestStack(t *testing.T) {
	s := &Frame{Register: NewRegister(0)}
	if err := s.Push(types.Boolean(true)); err != nil {
		t.Error(err)
	}
	fmt.Println(s)
	if err := s.Push(types.Integer(10)); err != nil {
		t.Error(err)
	}
	fmt.Println(s)
	if err := s.Push(types.GetNil()); err != nil {
		t.Error(err)
	}
	fmt.Println(s)
	if err := s.Push(types.String("hello")); err != nil {
		t.Error(err)
	}
	fmt.Println(s)
	if err := s.Push(types.Integer(-4)); err != nil {
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
	s := &Frame{Register: NewRegister(256)}
	assert.NoError(t, s.Push(types.Integer(1)))
	assert.NoError(t, s.Push(types.String("2.0")))
	assert.NoError(t, s.Push(types.String("3.0")))
	assert.NoError(t, s.Push(types.Float(4.0)))
	fmt.Println(s)
}

func TestBin(t *testing.T) {
	f, err := os.Open("testdata/luac.out")
	assert.NoError(t, err)
	var vm LuaVM
	assert.NoError(t, vm.Load(f))
	assert.NoError(t, vm.Execute())
}

func TestBin2(t *testing.T) {
	f, err := os.Open("testdata/test1.o")
	assert.NoError(t, err)
	var vm LuaVM
	assert.NoError(t, vm.Load(f))
	assert.NoError(t, vm.Execute())
}

