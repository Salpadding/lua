package vm

import (
	"fmt"
	"github.com/Salpadding/lua/types/value"
	"testing"
)

func TestNewStack(t *testing.T) {
	s := NewStack(1)
	if err := s.push(value.String("hello world")); err != nil{
		t.Error(err)
	}
	if err := s.push(value.String("hello world")); err == nil{
		t.Fail()
	}
	v, err := s.pop()
	if err != nil{
		t.Error(err)
	}
	fmt.Println(v)
	_, err = s.pop()
	if err == nil{
		t.Fail()
	}
}
