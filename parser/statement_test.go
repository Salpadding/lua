package parser

import (
	"bytes"
	"fmt"
	"testing"
)

func TestParseAssign(t *testing.T){
	p, err := New(bytes.NewBufferString(`
	a, a["bb"], a.name = "aa", 1 + 333, call("a", "b")
`))
	if err != nil {
		t.Error(err)
	}
	s, err := p.parseAssign()
	if err != nil{
		t.Error(err)
	}
	fmt.Println(s.String())
}

func TestParseLocalAssign(t *testing.T){
	p, err := New(bytes.NewBufferString(`
	local a, ab, ce = "aa", 1 + 333, call("a", "b")
`))
	if err != nil {
		t.Error(err)
	}
	s, err := p.parseLocalAssign()
	if err != nil{
		t.Error(err)
	}
	fmt.Println(s.String())
}

func TestParseSimples(t *testing.T){
	p, err := New(bytes.NewBufferString(`
	;
	break
	:: label ::
	goto label
`))
	if err != nil {
		t.Error(err)
	}
	stmts, err := p.parseStatements()
	if err != nil {
		t.Error(err)
	}
	for _, st := range stmts{
		fmt.Println(st.String())
	}
}

func TestParseBlock(t *testing.T){
	p, err := New(bytes.NewBufferString(`
	a = 1 + 2
	b = 1 + a
	return a + b, 1 ,2 ;
`))
	if err != nil {
		t.Error(err)
	}
	blk, err := p.parseBlock()
	if err != nil {
		t.Error(err)
	}
	for _, st := range blk.Statements{
		fmt.Println(st.String())
	}
	if blk.Return != nil{
		fmt.Println(blk.Return.String())
	}
}

