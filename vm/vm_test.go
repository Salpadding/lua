package vm

import (
	"fmt"
	"testing"

	"github.com/Salpadding/lua/types/value"
)

func TestNewStack(t *testing.T) {
	s := NewStack(1)
	if err := s.Push(value.String("hello world")); err != nil {
		t.Error(err)
	}
	if err := s.Push(value.String("hello world")); err == nil {
		t.Fail()
	}
	v, err := s.Pop()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(v)
	_, err = s.Pop()
	if err == nil {
		t.Fail()
	}
}

func TestIndex(t *testing.T) {
	s := NewStack(100)
	if err := s.Push(value.String("hello world")); err != nil {
		t.Error(err)
	}
	idx := s.AbsIndex(-1)
	if idx != 1 {
		t.Fail()
	}
	v := s.Get(1)
	_, ok := v.(*value.Nil)
	if ok {
		t.Fail()
	}
	if err := s.Set(1, value.Float(1.2)); err != nil {
		t.Error(err)
	}
	v = s.Get(1)
	_, ok = v.(value.Number)
	if !ok {
		t.Fail()
	}
}

func TestCheck(t *testing.T) {
	s := NewStack(1)
	if err := s.Push(value.String("hello world")); err != nil {
		t.Error(err)
	}
	s.Check(4)
	if len(s.slots) != 5 {
		t.Fail()
	}
	v := s.Get(1)
	if v.(value.String) != "hello world" {
		t.Fail()
	}
	if err := s.Push(value.String("hello world")); err != nil {
		t.Error(err)
	}
}

func TestStack(t *testing.T) {
	s := &LuaVM{Stack: NewStack(256)}
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
	if err := s.SetTop(6); err != nil{
		t.Error(err)
	}
	fmt.Println(s)
	if err := s.Remove(-3); err != nil{
		t.Error(err)
	}
	fmt.Println(s)
	if err := s.SetTop(-5); err != nil{
		t.Error(err)
	}
	fmt.Println(s)
}
