package lex

import (
	"bytes"
	"io"
	"strconv"

	"github.com/Salpadding/lua/common"
	"github.com/Salpadding/lua/token"
)

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

func (c character) rune() rune { return []rune(c)[0] }

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
		line:       1,
		column:     0,
	}
	l.nextChar()
	l.nextChar()
	return l
}

func (l *Lexer) readChar() Char {
	next, _, err := l.RuneReader.ReadRune()
	if err != nil {
		return eof(0)
	}
	return character(next)
}

func (l *Lexer) nextChar() Char {
	l.current = l.next
	l.next = l.readChar()
	if l.current == nil {
		return l.current
	}
	if l.current.isEOF() {
		return l.current
	}
	l.column++
	if l.current.rune() == '\r' && l.next.rune() != '\n' {
		l.column = 0
		l.line++
		return l.current
	}
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
	l.nextChar()
	l.nextChar()
	// single-line comment
	if l.current.rune() != '[' || l.next.rune() != '[' {
		for !l.current.isEOF() && l.current.rune() != '\n' {
			l.nextChar()
		}
		l.nextChar()
		return
	}
	l.nextChar()
	l.nextChar()
	// multi-line comment
	for !l.current.isEOF() {
		if l.current.rune() != '-' || l.next.rune() != '-' {
			l.nextChar()
			continue
		}
		l.nextChar()
		l.nextChar()
		if l.current.rune() == ']' && l.next.rune() == ']' {
			l.nextChar()
			l.nextChar()
			break
		}
	}
}

func (l *Lexer) skipWhiteSpaces() {
	// skip white spaces
	for !l.current.isEOF() && isWhiteSpace(l.current.rune()) {
		l.nextChar()
	}
}

func (l *Lexer) NextToken() (token.Token, error) {
	// skip white spaces
	l.skipWhiteSpaces()
	// skip comments
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
			l.nextChar()
			return tk, nil
		}
		tk := token.NewOperator(string(r)+string(n), l.line, l.column)
		l.nextChar()
		l.nextChar()
		return tk, nil
	case '+', '-', '*', '%', '&', '|', '^', '#':
		tk := token.NewOperator(string(l.current.rune()), l.line, l.column)
		l.nextChar()
		return tk, nil
	case ',', ';', '(', ')', ']', '{', '}':
		tk := token.NewDelimiter(string(l.current.rune()), l.line, l.column)
		l.nextChar()
		return tk, nil
	case ':':
		n := l.next.rune()
		if n != ':' {
			tk := token.NewDelimiter(string(r), l.line, l.column)
			l.nextChar()
			return tk, nil
		}
		tk := token.NewDelimiter(string(r)+string(n), l.line, l.column)
		l.nextChar()
		l.nextChar()
		return tk, nil
	case '.':
		n := l.next.rune()
		if n != '.' {
			tk := token.NewOperator(string(r), l.line, l.column)
			l.nextChar()
			return tk, nil
		}
		l.nextChar()
		n = l.next.rune()
		if n != '.' {
			tk := token.NewOperator("..", l.line, l.column)
			l.nextChar()
			return tk, nil
		}
		tk := token.NewDelimiter("...", l.line, l.column)
		l.nextChar()
		l.nextChar()
		return tk, nil
	case '[':
		n := l.next.rune()
		if n != '[' {
			tk := token.NewDelimiter(string(r), l.line, l.column)
			l.nextChar()
			return tk, nil
		}
		// here document
		line, column := l.line, l.column
		l.nextChar()
		l.nextChar()
		var buf bytes.Buffer
		for !l.current.isEOF() && !(l.current.rune() == ']' && l.next.rune() == ']') {
			buf.WriteRune(l.current.rune())
			l.nextChar()
		}
		l.nextChar()
		l.nextChar()
		return token.NewStringLiteral(buf.String(), line, column), nil
	case '"', '\'':
		line, column := l.line, l.column
		l.nextChar()
		var buf bytes.Buffer
		for !l.current.isEOF() && l.current.rune() != r {
			if l.current.rune() == '\\' && l.next.rune() == r {
				buf.WriteRune(r)
				l.nextChar()
				l.nextChar()
			}
			buf.WriteRune(l.current.rune())
			l.nextChar()
		}
		l.nextChar()
		escaped, err := common.FromEscaped(&buf)
		if err != nil {
			return nil, err
		}
		return token.NewStringLiteral(escaped, line, column), nil
	default:
		return l.readLiteralOrKeyword()
	}
}

func (l *Lexer) isNumber(r rune) bool {
	return '0' <= r && r <= '9'
}

func (l *Lexer) readIDOrKeyword() (token.Token, error) {
	line, column := l.line, l.column
	var buf bytes.Buffer
	for !l.current.isEOF() && (l.isID(l.current.rune()) || l.isNumber(l.current.rune())) {
		buf.WriteRune(l.current.rune())
		l.nextChar()
	}
	str := buf.String()
	// and or not is operator
	_, ok := token.Operators[str]
	if ok {
		return token.NewOperator(str, line, column), nil
	}
	// keyword lookup
	_, ok = token.Keywords[str]
	if ok {
		return token.NewKeyword(str, line, column), nil
	}
	return token.NewID(str, line, column), nil
}

func (l *Lexer) isHex(r rune) bool {
	return l.isNumber(r) || ('a' <= r && r <= 'f') || ('A' <= r && r <= 'F')
}

func (l *Lexer) isLetter(r rune) bool {
	return 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z'
}

func (l *Lexer) isID(r rune) bool {
	return r == '_' || l.isLetter(r)
}

func (l *Lexer) readLiteralOrKeyword() (token.Token, error) {
	line, column := l.line, l.column
	fst := l.current.rune()
	if !l.isNumber(fst) {
		return l.readIDOrKeyword()
	}
	// id starts with non-digital

	// peek snd rune
	snd := l.next.rune()
	// try to parse as hex number
	if snd == 'x' || snd == 'X' {
		l.nextChar()
		l.nextChar()
		var buf bytes.Buffer
		for !l.current.isEOF() && l.isHex(l.current.rune()) {
			buf.WriteRune(l.current.rune())
			l.nextChar()
		}
		str := buf.String()
		_, err := strconv.ParseInt(str, 16, 64)
		if err != nil {
			return nil, err
		}
		return token.NewNumberLiteral(str, 16, line, column), nil
	}
	// try to parse as digital number
	var buf bytes.Buffer
	for !l.current.isEOF() && (l.isNumber(l.current.rune()) || l.current.rune() == '.' || l.current.rune() == 'e') {
		buf.WriteRune(l.current.rune())
		l.nextChar()
	}
	_, err := strconv.ParseFloat(buf.String(), 64)
	if err != nil {
		return nil, err
	}
	return token.NewNumberLiteral(buf.String(), 10, line, column), nil
}
