package lex

import (
	"github.com/Salpadding/lua/token"
	"io"
)

type Char interface {
	rune() rune
	isEOF() bool
}

type eof rune

func (e eof) rune() rune { return 0 }

func (e eof) isEOF() bool { return true }

type character string

func (c character) rune() rune { return rune(c[0]) }

func (c character) isEOF() bool { return false }

type Lexer struct {
	io.RuneReader
	current Char
	next    Char
	line    int
	column  int
}

func New(reader io.RuneReader) *Lexer {
	l := &Lexer{
		RuneReader: reader,
		line:       0,
		column:     0,
	}
	l.ReadChar()
	l.ReadChar()
	return l
}

func (l *Lexer) readChar() Char {
	next, _, err := l.RuneReader.ReadRune()
	if err != nil {
		return eof(0)
	}
	return character(next)
}

func (l *Lexer) ReadChar() Char {
	l.current = l.next
	l.next = l.readChar()
	if l.current == nil || l.next == nil {
		return l.current
	}
	if l.current.rune() == '\r' && l.next.rune() == '\n' {
		l.current = l.next
		l.next = l.readChar()
	}
	if l.current.rune() == '\r' && l.next.rune() != '\n' {
		l.current = character('\n')
	}
	l.column++
	if l.current.rune() == '\n' {
		l.column = 0
		l.line++
	}
	return l.current
}

func isWhiteSpace(r rune) bool {
	switch r {
	case '\t', '\n', '\v', '\f', '\r', ' ':
		return true
	}
	return false
}

func (l *Lexer) skipComment() {
	// skip comments
	if l.current.rune() != '-' || l.next.rune() != '-' {
		return
	}
	l.ReadChar()
	l.ReadChar()
	if l.current.rune() != '[' || l.next.rune() != '[' {
		for !l.current.isEOF() && l.current.rune() != '\n' {
			l.ReadChar()
		}
		l.ReadChar()
		return
	}
	l.ReadChar()
	l.ReadChar()
	for !l.current.isEOF() {
		if l.current.rune() != '-' || l.next.rune() != '-' {
			l.ReadChar()
			continue
		}
		l.ReadChar()
		l.ReadChar()
		if l.current.rune() == ']' && l.next.rune() == ']' {
			l.ReadChar()
			l.ReadChar()
			break
		}
	}
}

func(l *Lexer) skipWhiteSpaces(){
	// skip white spaces
	for !l.current.isEOF() && isWhiteSpace(l.current.rune()) {
		l.ReadChar()
	}
}

func (l *Lexer) NextToken() (token.Token, error) {
	// 跳过空白
	l.skipWhiteSpaces()
	// 发现注释则跳过
	if l.current.rune() == '-' && l.next.rune() == '-' {
		l.skipComment()
		return l.NextToken()
	}
	if l.current.isEOF() {
		return token.EOF("EOF"), nil
	}
	return nil, nil
}
