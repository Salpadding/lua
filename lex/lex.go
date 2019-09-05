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
	// single-line comment
	if l.current.rune() != '[' || l.next.rune() != '[' {
		for !l.current.isEOF() && l.current.rune() != '\n' {
			l.ReadChar()
		}
		l.ReadChar()
		return
	}
	l.ReadChar()
	l.ReadChar()
	// multi-line comment
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
	case ',', ';', '(', ')', ']', '{', '}':
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
		escaped, err := common.FromEscaped(&buf)
		if err != nil {
			return nil, err
		}
		return token.NewLiteral(token.String, escaped, line, column), nil
	default:
		return l.readLiteralOrKeyword()
	}
}

func (l *Lexer) isNumber(r rune) bool {
	return '0' <= r && r <= '9'
}

func (l *Lexer) readLiteralOrKeyword() (token.Token, error) {
	var buf bytes.Buffer
	line, column := l.line, l.column
	for !l.current.isEOF() && !isWhiteSpace(l.current.rune()) {
		buf.WriteRune(l.current.rune())
		l.ReadChar()
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
	fst := []rune(str)[0]
	// id starts with non-digital
	if !l.isNumber(fst) {
		return token.NewID(str, line, column), nil
	}
	// peek snd rune
	if len([]rune(str)) == 1 {
		return token.NewLiteral(token.Number, str, line, column), nil
	}
	snd := []rune(str)[1]
	if snd == 'x' {
		_, err := strconv.ParseInt(str[2:], 16, 64)
		if err != nil {
			return nil, err
		}
		return token.NewLiteral(token.Number, str, line, column), nil
	}
	_, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return nil, err
	}
	return token.NewLiteral(token.Number, str, line, column), nil
}
