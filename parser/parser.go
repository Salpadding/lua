package parser

import (
	"github.com/Salpadding/lua/ast"
	"github.com/Salpadding/lua/lex"
	"github.com/Salpadding/lua/token"
	"io"
	"strconv"
)

/*
exp   ::= exp12
exp12 ::= exp11 {or exp11}
exp11 ::= exp10 {and exp10}
exp10 ::= exp9 {('<' | '>' | '<=' | '>=' | '~=' | '==') exp9}
exp9  ::= exp8 {'|' exp8}
exp8  ::= exp7 {'~' exp7}
exp7  ::= exp6 {'&' exp6}
exp6  ::= exp5 {('<<' | '>>') exp5}
exp5  ::= exp4 {'..' exp5}
exp4  ::= exp3 {('+' | '-') exp3}
exp3  ::= exp2 {('*' | '/' | '//' | '%') exp2}
exp2  ::= {('not' | '#' | '-' | '~')} exp2
exp1  ::= exp0 {'^' exp1}
exp0  ::= nil | false | true | Numeral | LiteralString | ID
		| '...' | (exp) | (ID | ( ) (' (exp,)*exp ')' | functionLiteral
*/

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
	if _, err := p.nextToken(); err != nil {
		return nil, err
	}
	if _, err := p.nextToken(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Parser) parseExpression() (ast.Expression, error) {
	return p.parseExp12()
}

func (p *Parser) parseExp12() (ast.Expression, error) {
	left, err := p.parseExp11()
	if err != nil {
		return nil, err
	}
	for {
		op := p.current
		switch op.Type() {
		case token.LogicalOr:
			if _, err = p.nextToken(); err != nil {
				return nil, err
			}
			right, err := p.parseExp11()
			if err != nil {
				return nil, err
			}
			left = &ast.InfixExpression{
				Operator: op.(*token.Operator),
				Left:     left,
				Right:    right,
			}
		default:
			return left, nil
		}
	}
}

func (p *Parser) parseExp11() (ast.Expression, error) {
	left, err := p.parseExp10()
	if err != nil {
		return nil, err
	}
	for {
		op := p.current
		switch op.Type() {
		case token.LogicalAnd:
			if _, err = p.nextToken(); err != nil {
				return nil, err
			}
			right, err := p.parseExp10()
			if err != nil {
				return nil, err
			}
			left = &ast.InfixExpression{
				Operator: op.(*token.Operator),
				Left:     left,
				Right:    right,
			}
		default:
			return left, nil
		}
	}
}

func (p *Parser) parseExp10() (ast.Expression, error) {
	left, err := p.parseExp9()
	if err != nil {
		return nil, err
	}
	for {
		op := p.current
		switch op.Type() {
		case token.LessThan, token.LessThanOrEqual, token.Equal,
			token.GreaterThan, token.GreaterThanOrEqual, token.NotEqual:
			if _, err = p.nextToken(); err != nil {
				return nil, err
			}
			right, err := p.parseExp9()
			if err != nil {
				return nil, err
			}
			left = &ast.InfixExpression{
				Operator: op.(*token.Operator),
				Left:     left,
				Right:    right,
			}
		default:
			return left, nil
		}
	}
}

func (p *Parser) parseExp9() (ast.Expression, error) {
	left, err := p.parseExp8()
	if err != nil {
		return nil, err
	}
	for {
		op := p.current
		switch op.Type() {
		case token.BitwiseOr:
			if _, err = p.nextToken(); err != nil {
				return nil, err
			}
			right, err := p.parseExp8()
			if err != nil {
				return nil, err
			}
			left = &ast.InfixExpression{
				Operator: op.(*token.Operator),
				Left:     left,
				Right:    right,
			}
		default:
			return left, nil
		}
	}
}

func (p *Parser) parseExp8() (ast.Expression, error) {
	left, err := p.parseExp7()
	if err != nil {
		return nil, err
	}
	for {
		op := p.current
		switch op.Type() {
		case token.Wave:
			if _, err = p.nextToken(); err != nil {
				return nil, err
			}
			right, err := p.parseExp7()
			if err != nil {
				return nil, err
			}
			left = &ast.InfixExpression{
				Operator: op.(*token.Operator),
				Left:     left,
				Right:    right,
			}
		default:
			return left, nil
		}
	}
}

func (p *Parser) parseExp7() (ast.Expression, error) {
	left, err := p.parseExp6()
	if err != nil {
		return nil, err
	}
	for {
		op := p.current
		switch op.Type() {
		case token.BitwiseAnd:
			if _, err = p.nextToken(); err != nil {
				return nil, err
			}
			right, err := p.parseExp6()
			if err != nil {
				return nil, err
			}
			left = &ast.InfixExpression{
				Operator: op.(*token.Operator),
				Left:     left,
				Right:    right,
			}
		default:
			return left, nil
		}
	}
}

func (p *Parser) parseExp6() (ast.Expression, error) {
	left, err := p.parseExp5()
	if err != nil {
		return nil, err
	}
	for {
		op := p.current
		switch op.Type() {
		case token.LeftShift, token.RightShift:
			if _, err = p.nextToken(); err != nil {
				return nil, err
			}
			right, err := p.parseExp5()
			if err != nil {
				return nil, err
			}
			left = &ast.InfixExpression{
				Operator: op.(*token.Operator),
				Left:     left,
				Right:    right,
			}
		default:
			return left, nil
		}
	}
}

func (p *Parser) parseExp5() (ast.Expression, error) {
	left, err := p.parseExp4()
	if err != nil {
		return nil, err
	}
	op := p.current
	if op.Type() != token.Concat {
		return left, nil
	}
	if _, err = p.nextToken(); err != nil {
		return nil, err
	}
	right, err := p.parseExp5()
	if err != nil {
		return nil, err
	}
	return &ast.InfixExpression{
		Operator: op.(*token.Operator),
		Left:     left,
		Right:    right,
	}, nil
}

func (p *Parser) parseExp4() (ast.Expression, error) {
	left, err := p.parseExp3()
	if err != nil {
		return nil, err
	}
	for {
		op := p.current
		switch op.Type() {
		case token.Minus, token.Plus:
			if _, err = p.nextToken(); err != nil {
				return nil, err
			}
			right, err := p.parseExp3()
			if err != nil {
				return nil, err
			}
			left = &ast.InfixExpression{
				Operator: op.(*token.Operator),
				Left:     left,
				Right:    right,
			}
		default:
			return left, nil
		}
	}
}

func (p *Parser) parseExp3() (ast.Expression, error) {
	left, err := p.parseExp2()
	if err != nil {
		return nil, err
	}
	for {
		op := p.current
		switch op.Type() {
		case token.Asterisk, token.Divide, token.IntegerDivide, token.Modular:
			if _, err = p.nextToken(); err != nil {
				return nil, err
			}
			right, err := p.parseExp2()
			if err != nil {
				return nil, err
			}
			left = &ast.InfixExpression{
				Operator: op.(*token.Operator),
				Left:     left,
				Right:    right,
			}
		default:
			return left, nil
		}
	}
}

func (p *Parser) parseExp2() (ast.Expression, error) {
	op := p.current
	switch op.Type() {
	case token.LogicalNot, token.Len, token.Minus, token.Wave:
		if _, err := p.nextToken(); err != nil {
			return nil, err
		}
		exp, err := p.parseExp2()
		if err != nil {
			return nil, err
		}
		return &ast.PrefixExpression{
			Operator: op.(*token.Operator),
			Right:    exp,
		}, nil
	default:
		return p.parseExp1()
	}
}

func (p *Parser) parseExp1() (ast.Expression, error) {
	left, err := p.parseExp0()
	if err != nil {
		return nil, err
	}

	for {
		op := p.current
		if op.Type() != token.Power {
			return left, nil
		}
		if _, err = p.nextToken(); err != nil {
			return nil, err
		}
		right, err := p.parseExp1()
		if err != nil {
			return nil, err
		}
		left = &ast.InfixExpression{
			Operator: op.(*token.Operator),
			Left:     left,
			Right:    right,
		}
	}
}

func (p *Parser) parseExp0() (ast.Expression, error) {
	current := p.current
	switch c := current.(type) {
	case *token.NumberLiteral:
		if c.Base() == 10 {
			f, err := strconv.ParseFloat(c.Literal(), 64)
			if err != nil {
				return nil, err
			}
			if _, err = p.nextToken(); err != nil {
				return nil, err
			}
			return ast.Number(f), nil
		}
		n, err := strconv.ParseInt(c.Literal(), 16, 64)
		if err != nil {
			return nil, err
		}
		if _, err = p.nextToken(); err != nil {
			return nil, err
		}
		return ast.Number(n), nil
	case *token.StringLiteral:
		if _, err := p.nextToken(); err != nil {
			return nil, err
		}
		return ast.String(c.Literal()), nil
	case *token.Keyword:
		switch c.Type() {
		case token.True:
			if _, err := p.nextToken(); err != nil {
				return nil, err
			}
			return ast.Boolean(true), nil
		case token.False:
			if _, err := p.nextToken(); err != nil {
				return nil, err
			}
			return ast.Boolean(false), nil
		case token.Nil:
			if _, err := p.nextToken(); err != nil {
				return nil, err
			}
			return &ast.Nil{}, nil
		default:
			return nil, errUnexpectedError(current)
		}
	default:
		return p.parsePrefix0()
	}
}

/*
	prefix2: prefix1 | prefix1 '(' ')' | prefix1 '(' (exp ',')* exp ')'
	prefix1: prefix0 ('.' id)* | prefix0 ('[' prefix1 ']')*;
	prefix0: id | '(' exp ')';
 */
func(p *Parser) parsePrefix1() (ast.Expression, error){
	return nil, nil
}


func (p *Parser) parsePrefix0() (ast.Expression, error) {
	current := p.current
	switch p.current.Type() {
	case token.LeftParenthesis:
		if _, err := p.nextToken(); err != nil {
			return nil, err
		}
		exp, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if p.current.Type() != token.RightParenthesis {
			return nil, errUnexpectedError(p.current)
		}
		if _, err := p.nextToken(); err != nil {
			return nil, err
		}
		return exp, nil
	case token.Identifier:
		if _, err := p.nextToken(); err != nil {
			return nil, err
		}
		return ast.Identifier(current.String()), nil
	default:
		return nil, errUnexpectedError(p.current)
	}
}
