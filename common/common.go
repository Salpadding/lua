package common

import (
	"bytes"
	"errors"
	"io"
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
			buf.WriteRune(escapes[n])
		default:
			buf.WriteRune(n)
		}
	}
	return buf.String(), nil
}
