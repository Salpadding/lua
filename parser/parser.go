package parser

import (
	"github.com/Salpadding/lua/lex"
	"github.com/Salpadding/lua/token"
	"io"
)

type Parser struct {
	*lex.Lexer
	current token.Token
	next    token.Token
}

func (p *Parser) nextToken() (token.Token, error) {
	p.current = p.next
	next, err := p.Lexer.NextToken()
	if err != nil {
		return nil, err
	}
	p.next = next
	return p.current, nil
}

func New(reader io.RuneReader) (*Parser, error) {
	p := &Parser{
		Lexer: lex.New(reader),
	}
	if _, err := p.NextToken(); err != nil {
		return nil, err
	}
	if _, err := p.NextToken(); err != nil {
		return nil, err
	}
	return p, nil
}
