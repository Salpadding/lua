package parser

import (
	"bytes"
	"fmt"
	"github.com/Salpadding/lua/ast"
	"github.com/Salpadding/lua/token"
	"testing"
)

func TestParse0(t *testing.T){
	p, err := New(bytes.NewBufferString(`
	12 "abc" true false nil
`))
	if err != nil{
		t.Error(err)
	}
	for p.current.Type() != token.EndOfFile{
		exp, err := p.parseExp0()
		if err != nil{
			t.Error(err)
		}
		_, ok := exp.(ast.Expression)
		if !ok{
			t.Fail()
		}
		fmt.Println(exp.String())
	}
}

func TestParse1(t *testing.T){
	p, err := New(bytes.NewBufferString(`
	2 ^ 10
`))
	if err != nil{
		t.Error(err)
	}
	for p.current.Type() != token.EndOfFile{
		exp, err := p.parseExp1()
		if err != nil{
			t.Error(err)
		}
		_, ok := exp.(ast.Expression)
		if !ok{
			t.Fail()
		}
		fmt.Println(exp.String())
	}
}
