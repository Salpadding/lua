package parser

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/Salpadding/lua/token"
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
	for p.current.Type() != token.EndOfFile{
		s, err := p.parseStatement()
		if err != nil{
			t.Error(err)
		}
		fmt.Println(s.String())
	}
}

