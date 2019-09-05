package lex

import (
	"bytes"
	"fmt"
	"testing"
)

func Test(t *testing.T) {
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
