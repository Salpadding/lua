package lex

import (
	"bytes"
	"errors"
	"github.com/Salpadding/lua/token"
	"io"
	"strconv"
)

var escapes = map[rune]rune{
	'a':  '\a',
	'b':  '\b',
	'f':  '\f',
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
	'v':  '\v',
	'"':  '"',
	'\\': '\\',
	'\'': '\'',
}

var ops = map[rune]map[rune]bool{
	'=': {'=': true},
	'<': {'=': true, '<': true},
	'>': {'=': true, '>': true},
	'~': {'=': true},
	'/': {'/': true},
}

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

func (l *Lexer) skipWhiteSpaces() {
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
	r := l.current.rune()
	switch r {
	case '=', '<', '>', '~', '/':
		n := l.next.rune()
		ok := ops[r][n]
		if !ok {
			tk := token.NewOperator(string(r), l.line, l.column)
			l.ReadChar()
			return tk, nil
		}
		tk := token.NewOperator(string(r)+string(n), l.line, l.column)
		l.ReadChar()
		l.ReadChar()
		return tk, nil
	case '+', '-', '*', '%', '&', '|', '^', '#':
		tk := token.NewOperator(string(l.current.rune()), l.line, l.column)
		l.ReadChar()
		return tk, nil
	case ',', ';', '(', ')', ']':
		tk := token.NewDelimiter(string(l.current.rune()), l.line, l.column)
		l.ReadChar()
		return tk, nil
	case ':':
		n := l.next.rune()
		if n != ':' {
			tk := token.NewDelimiter(string(r), l.line, l.column)
			l.ReadChar()
			return tk, nil
		}
		tk := token.NewDelimiter(string(r)+string(n), l.line, l.column)
		l.ReadChar()
		l.ReadChar()
		return tk, nil
	case '.':
		n := l.next.rune()
		if n != '.' {
			tk := token.NewOperator(string(r), l.line, l.column)
			l.ReadChar()
			return tk, nil
		}
		l.ReadChar()
		n = l.next.rune()
		if n != '.' {
			tk := token.NewOperator("..", l.line, l.column)
			l.ReadChar()
			return tk, nil
		}
		tk := token.NewDelimiter("...", l.line, l.column)
		l.ReadChar()
		l.ReadChar()
		return tk, nil
	case '[':
		n := l.next.rune()
		if n != '[' {
			tk := token.NewDelimiter(string(r), l.line, l.column)
			l.ReadChar()
			return tk, nil
		}
		// here document
		line, column := l.line, l.column
		l.ReadChar()
		l.ReadChar()
		var buf bytes.Buffer
		for !l.current.isEOF() && !(l.current.rune() == ']' && l.next.rune() == ']') {
			buf.WriteRune(l.current.rune())
			l.ReadChar()
		}
		l.ReadChar()
		l.ReadChar()
		return token.NewLiteral(token.String, buf.String(), line, column), nil
	case '"':
		line, column := l.line, l.column
		l.ReadChar()
		var buf bytes.Buffer
		for !l.current.isEOF() && l.current.rune() != '"' {
			buf.WriteRune(l.current.rune())
			l.ReadChar()
		}
		l.ReadChar()
		escaped, err := l.escape(&buf)
		if err != nil{
			return nil, err
		}
		return token.NewLiteral(token.String, escaped, line, column), nil
	default:
		return l.readLiteralOrKeyword()
	}
}

func (l *Lexer) escape(rd io.RuneReader) (string, error) {
	var buf bytes.Buffer
	for {
		r, _, err := rd.ReadRune()
		if err != nil {
			break
		}
		if r != '\\' {
			buf.WriteRune(r)
		}
		n, _, err := rd.ReadRune()
		if err != nil {
			return "", errors.New("unexpected string end after \\")
		}
		switch n {
		case 'a', 'b', 'f', 'n', 'r', 't', 'v', '"', '\'', '\\':
			buf.WriteRune(escapes[n])
		default:
			buf.WriteRune(n)
		}
	}
	return buf.String(), nil
}

func (l *Lexer) readLiteralOrKeyword() (token.Token, error) {
	var buf bytes.Buffer
	line, column := l.line, l.column
	for !l.current.isEOF() && !isWhiteSpace(l.current.rune()) {
		buf.WriteRune(l.current.rune())
		l.ReadChar()
	}
	str := buf.String()
	_, ok := token.Operators[str]
	if ok {
		return token.NewOperator(str, line, column), nil
	}
	_, ok = token.Keywords[str]
	if ok {
		return token.NewKeyword(str, line, column), nil
	}
	_, err := strconv.ParseFloat(str, 64)
	if err == nil {
		return token.NewLiteral(token.Number, str, line, column), nil
	}
	return token.NewID(str, line, column), nil
}
