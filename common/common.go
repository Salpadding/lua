package common

import (
	"bytes"
	"errors"
	"io"
)

var escapeChars = map[rune]rune{
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

var escapes = map[rune]string{
	'\a': `\a`,
	'\b': `\b`,
	'\f': `\f`,
	'\n': `\n`,
	'\r': `\r`,
	'\t': `\t`,
	'\v': `\v`,
	'"':  `\"`,
	'\\': `\\`,
	'\'': `\'`,
}

func FromEscaped(rd io.RuneReader) (string, error) {
	var buf bytes.Buffer
	for {
		r, _, err := rd.ReadRune()
		if err != nil {
			break
		}
		if r != '\\' {
			buf.WriteRune(r)
			continue
		}
		n, _, err := rd.ReadRune()
		if err != nil {
			return "", errors.New("unexpected string end after \\")
		}
		switch n {
		case 'a', 'b', 'f', 'n', 'r', 't', 'v', '"', '\'', '\\':
			buf.WriteRune(escapeChars[n])
		default:
			buf.WriteRune(n)
		}
	}
	return buf.String(), nil
}

func Escape(rd io.RuneReader) string {
	var buf bytes.Buffer
	for {
		r, _, err := rd.ReadRune()
		if err != nil {
			break
		}
		s, ok := escapes[r]
		if ok {
			buf.WriteString(s)
			continue
		}
		buf.WriteRune(r)
	}
	return buf.String()
}
