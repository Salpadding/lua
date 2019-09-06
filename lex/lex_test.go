package lex

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/Salpadding/lua/token"
)

func Test1(t *testing.T) {
	l := &Lexer{
		RuneReader: bytes.NewBufferString("\r\r12345\r\n123\n123\r123\r\n\r\n123"),
		current:    nil,
		next:       nil,
		line:       1,
		column:     0,
	}
	for l.current == nil || !l.current.isEOF() {
		l.ReadChar()
		if l.current != nil && !l.current.isEOF() && l.column != 0 {
			if l.current.rune() == '\n' || l.current.rune() == '\r' {
				continue
			}
			s := string(l.current.rune())
			fmt.Printf("%s at line %d   column %d\n", s, l.line, l.column)
		}
	}
}

func TestSkipComments(t *testing.T) {
	l := New(bytes.NewBufferString(`-- 这是单行注释
		--[[ 这是多行
		多行注释
		--]]
		-- comments starts with two '-'
`))
	var tokens []token.Token
	for {
		tk, err := l.NextToken()
		if err != nil {
			t.Error(err)
		}
		tokens = append(tokens, tk)
		if tk.Type() == token.EndOfFile {
			break
		}
	}
	if len(tokens) < 0 {
		t.Fail()
	}
	if tokens[0].Type() != token.EndOfFile {
		t.Fail()
	}
	fmt.Println(tokens)
}

func TestOperators(t *testing.T) {
	var buf bytes.Buffer
	for k := range token.Operators {
		buf.WriteString(k)
		buf.WriteRune(' ')
	}
	l := New(&buf)
	var tokens []token.Token
	for {
		tk, err := l.NextToken()
		if err != nil {
			t.Error(err)
		}
		tokens = append(tokens, tk)
		if tk.Type() == token.EndOfFile {
			break
		}
	}
	for _, tk := range tokens {
		if tk.Type() == token.EndOfFile {
			break
		}
		_, ok := tk.(*token.Operator)
		if !ok {
			t.Fail()
		}
		fmt.Println(tk.String())
	}
}

func TestDelimiters(t *testing.T) {
	var buf bytes.Buffer
	for k := range token.Delimiters {
		buf.WriteString(k)
		buf.WriteRune(' ')
	}
	l := New(&buf)
	var tokens []token.Token
	for {
		tk, err := l.NextToken()
		if err != nil {
			t.Error(err)
		}
		tokens = append(tokens, tk)
		if tk.Type() == token.EndOfFile {
			break
		}
	}
	for _, tk := range tokens {
		if tk.Type() == token.EndOfFile {
			break
		}
		_, ok := tk.(*token.Delimiter)
		if !ok {
			t.Fail()
		}
		fmt.Println(tk.String())
	}
}

func TestKeywords(t *testing.T) {
	var buf bytes.Buffer
	for k := range token.Keywords {
		buf.WriteString(k)
		buf.WriteRune(' ')
	}
	l := New(&buf)
	var tokens []token.Token
	for {
		tk, err := l.NextToken()
		if err != nil {
			t.Error(err)
		}
		tokens = append(tokens, tk)
		if tk.Type() == token.EndOfFile {
			break
		}
	}
	for _, tk := range tokens {
		if tk.Type() == token.EndOfFile {
			break
		}
		_, ok := tk.(*token.Keyword)
		if !ok {
			t.Fail()
		}
		fmt.Println(tk.String())
	}
}

func TestStringLiteral(t *testing.T) {
	l := New(bytes.NewBufferString(`
	[[ 这是多行文本
 ++++++ -----
 ]] "这是单行文本"
`))
	var tokens []token.Token
	for {
		tk, err := l.NextToken()
		if err != nil {
			t.Error(err)
		}
		tokens = append(tokens, tk)
		if tk.Type() == token.EndOfFile {
			break
		}
	}
	for _, tk := range tokens {
		if tk.Type() == token.EndOfFile {
			break
		}
		if tk.Type() != token.String {
			t.Fail()
		}
		fmt.Println(tk.String())
	}
}

func TestNumberLiteral(t *testing.T) {
	l := New(bytes.NewBufferString(`
	1 2 3 3.1 333 3e10 0303 0xffff
`))
	var tokens []token.Token
	for {
		tk, err := l.NextToken()
		if err != nil {
			t.Error(err)
		}
		tokens = append(tokens, tk)
		if tk.Type() == token.EndOfFile {
			break
		}
	}
	for _, tk := range tokens {
		if tk.Type() == token.EndOfFile {
			break
		}
		if tk.Type() != token.Number {
			t.Fail()
		}
		fmt.Println(tk.String())
	}
}

func TestID(t *testing.T) {
	l := New(bytes.NewBufferString(`
	abc def ghi
`))
	var tokens []token.Token
	for {
		tk, err := l.NextToken()
		if err != nil {
			t.Error(err)
		}
		tokens = append(tokens, tk)
		if tk.Type() == token.EndOfFile {
			break
		}
	}
	for _, tk := range tokens {
		if tk.Type() == token.EndOfFile {
			break
		}
		if tk.Type() != token.Identifier {
			t.Fail()
		}
		fmt.Println(tk.String())
	}
}
