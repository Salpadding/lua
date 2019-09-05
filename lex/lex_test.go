package lex

import (
	"bytes"
	"fmt"
	"github.com/Salpadding/lua/token"
	"testing"
)

func Test1(t *testing.T) {
	l := &Lexer{
		RuneReader: bytes.NewBufferString("12345\r\n123\n123\r123\r\n\r\n123"),
		current:    nil,
		next:       nil,
		line:       0,
		column:     0,
	}
	for l.current == nil || !l.current.isEOF() {
		l.ReadChar()
		if l.current != nil && !l.current.isEOF() && l.column != 0{
			s := string(l.current.rune())
			fmt.Printf("%s at line %d   column %d\n", s, l.line, l.column)
		}
	}
}

func TestSkipComments(t *testing.T){
	l := New(bytes.NewBufferString(`-- 这是一行注释
		--[[ 这是多行
		多行注释
		--]]
`))
	var tokens []token.Token
	for {
		tk, err := l.NextToken()
		if err != nil{
			t.Error(err)
		}
		tokens = append(tokens, tk)
		if tk.Type() == token.EndOfFile{
			break
		}
	}
	fmt.Println(tokens)
}

func TestOperators(t *testing.T){
	l := New(bytes.NewBufferString(`
		+ - * % & | ^ = == # <= < >= > << >> ~ ~= / // ~ ~=
`))
	var tokens []token.Token
	for {
		tk, err := l.NextToken()
		if err != nil{
			t.Error(err)
		}
		tokens = append(tokens, tk)
		if tk.Type() == token.EndOfFile{
			break
		}
	}
	for _, tk := range tokens{
		fmt.Println(tk.String())
	}
}

func TestDelimiters(t *testing.T){
	l := New(bytes.NewBufferString(`
	 , ;  (  )  [ ] : :: . .. ... . .. ...
`))
	var tokens []token.Token
	for {
		tk, err := l.NextToken()
		if err != nil{
			t.Error(err)
		}
		tokens = append(tokens, tk)
		if tk.Type() == token.EndOfFile{
			break
		}
	}
	for _, tk := range tokens{
		fmt.Println(tk.String())
	}
}

func TestLiteralKeywords(t *testing.T){
	l := New(bytes.NewBufferString(`
	 for break goto function 12 333 "aaaa" [[
aaa ff
]]
`))
	var tokens []token.Token
	for {
		tk, err := l.NextToken()
		if err != nil{
			t.Error(err)
		}
		tokens = append(tokens, tk)
		if tk.Type() == token.EndOfFile{
			break
		}
	}
	for _, tk := range tokens{
		fmt.Println(tk.String())
	}
}