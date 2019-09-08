package parser

import (
	"strconv"

	"github.com/Salpadding/lua/ast"
	"github.com/Salpadding/lua/token"
)

func (p *Parser) parseExp12() (ast.Expression, error) {
	left, err := p.parseExp11()
	if err != nil {
		return nil, err
	}
	for {
		op := p.current
		switch op.Type() {
		case token.LogicalOr:
			if _, err = p.nextToken(1); err != nil {
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
			if _, err = p.nextToken(1); err != nil {
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
			if _, err = p.nextToken(1); err != nil {
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
			if _, err = p.nextToken(1); err != nil {
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
			if _, err = p.nextToken(1); err != nil {
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
			if _, err = p.nextToken(1); err != nil {
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
			if _, err = p.nextToken(1); err != nil {
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
	if _, err = p.nextToken(1); err != nil {
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
			if _, err = p.nextToken(1); err != nil {
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
			if _, err = p.nextToken(1); err != nil {
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
		if _, err := p.nextToken(1); err != nil {
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
		if _, err = p.nextToken(1); err != nil {
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
			if _, err = p.nextToken(1); err != nil {
				return nil, err
			}
			return ast.Number(f), nil
		}
		n, err := strconv.ParseInt(c.Literal(), 16, 64)
		if err != nil {
			return nil, err
		}
		if _, err = p.nextToken(1); err != nil {
			return nil, err
		}
		return ast.Number(n), nil
	case *token.StringLiteral:
		if _, err := p.nextToken(1); err != nil {
			return nil, err
		}
		return ast.String(c.Literal()), nil
	case *token.Keyword:
		switch c.Type() {
		case token.True:
			if _, err := p.nextToken(1); err != nil {
				return nil, err
			}
			return ast.Boolean(true), nil
		case token.False:
			if _, err := p.nextToken(1); err != nil {
				return nil, err
			}
			return ast.Boolean(false), nil
		case token.Nil:
			if _, err := p.nextToken(1); err != nil {
				return nil, err
			}
			return &ast.Nil{}, nil
		default:
			return p.parsePrefix1()
		}
	default:
		if current.Type() == token.Varying {
			if _, err := p.nextToken(1); err != nil {
				return nil, err
			}
			return ast.Vararg(current.String()), nil
		}
		if current.Type() == token.LeftBrace {
			return p.parseTable()
		}
		return p.parsePrefix1()
	}
}

func (p *Parser) parsePrefix1() (ast.Expression, error) {
	left, err := p.parsePrefix0()
	if err != nil {
		return nil, err
	}
	for {
		switch p.current.Type() {
		case token.Dot:
			if p.next.Type() != token.Identifier {
				return nil, errUnexpectedError(p.next)
			}
			left = &ast.TableAccess{
				Left:  left,
				Index: ast.Identifier(p.next.String()),
			}
			if _, err = p.nextToken(2); err != nil {
				return nil, err
			}
		case token.LeftBracket:
			if _, err = p.nextToken(1); err != nil {
				return nil, err
			}
			idx, err := p.parseExp12()
			if err != nil {
				return nil, err
			}
			left = &ast.TableAccess{
				Left:  left,
				Index: idx,
			}
			if p.current.Type() != token.RightBracket {
				return nil, errUnexpectedError(p.current)
			}
			if _, err = p.nextToken(1); err != nil {
				return nil, err
			}
		case token.LeftParenthesis:
			if p.next.Type() == token.RightParenthesis {
				left = &ast.FunctionCall{
					Function: left,
					Args:     nil,
				}
				if _, err = p.nextToken(2); err != nil {
					return nil, err
				}
				continue
			}
			if _, err = p.nextToken(1); err != nil {
				return nil, err
			}
			args, err := p.parseExpressions()
			if err != nil {
				return nil, err
			}
			if p.current.Type() != token.RightParenthesis {
				return nil, errUnexpectedError(p.current)
			}
			left = &ast.FunctionCall{
				Function: left,
				Args:     args,
			}
			if _, err = p.nextToken(1); err != nil {
				return nil, err
			}
		case token.LeftBrace:
			args, err := p.parseTable()
			if err != nil {
				return nil, err
			}
			left = &ast.FunctionCall{
				Function: left,
				Args:     args,
			}
		case token.String:
			left = &ast.FunctionCall{
				Function: left,
				Args:     ast.String(p.current.(*token.StringLiteral).Literal()),
			}
			if _, err = p.nextToken(1); err != nil {
				return nil, err
			}
		case token.Colon:
			if p.next.Type() != token.Identifier {
				return nil, errUnexpectedError(p.next)
			}
			id := p.next.String()
			if _, err = p.nextToken(2); err != nil {
				return nil, err
			}
			args, err := p.parseArguments()
			if err != nil {
				return nil, err
			}
			return &ast.FunctionCall{
				Function: ast.Identifier(id),
				Args:     args,
				Self:     left,
			}, nil
		default:
			return left, nil
		}
	}
}

func (p *Parser) parsePrefix0() (ast.Expression, error) {
	current := p.current
	switch p.current.Type() {
	case token.LeftParenthesis:
		if _, err := p.nextToken(1); err != nil {
			return nil, err
		}
		exp, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if p.current.Type() != token.RightParenthesis {
			return nil, errUnexpectedError(p.current)
		}
		if _, err := p.nextToken(1); err != nil {
			return nil, err
		}
		return exp, nil
	case token.Identifier:
		if _, err := p.nextToken(1); err != nil {
			return nil, err
		}
		return ast.Identifier(current.String()), nil
	default:
		return nil, errUnexpectedError(p.current)
	}
}

func (p *Parser) parseTable() (ast.Table, error) {
	// skip '{'
	if _, err := p.nextToken(1); err != nil {
		return nil, err
	}
	var pairs ast.Table
	i := 1
	for p.current.Type() != token.RightBrace {
		switch p.current.Type() {
		case token.LeftBracket:
			if _, err := p.nextToken(1); err != nil {
				return nil, err
			}
			k, err := p.parseExp12()
			if err != nil {
				return nil, err
			}
			if err = p.assertCurrentAndSkip(token.RightBracket); err != nil {
				return nil, err
			}
			if err = p.assertCurrentAndSkip(token.Assign); err != nil {
				return nil, err
			}
			v, err := p.parseExp12()
			if err != nil {
				return nil, err
			}
			pairs = append(pairs, &ast.Keypair{
				Key:   k,
				Value: v,
			})
		case token.Identifier:
			id := p.current.String()
			if p.next.Type() != token.Assign {
				return nil, errUnexpectedError(p.next)
			}
			if _, err := p.nextToken(2); err != nil {
				return nil, err
			}
			v, err := p.parseExp12()
			if err != nil {
				return nil, err
			}
			pairs = append(pairs, &ast.Keypair{
				Key:   ast.String(id),
				Value: v,
			})
		default:
			v, err := p.parseExp12()
			if err != nil {
				return nil, err
			}
			pairs = append(pairs, &ast.Keypair{
				Key:   ast.Number(i),
				Value: v,
			})
			i++
		}
		if p.current.Type() != token.Comma && p.current.Type() != token.Semicolon {
			break
		}
		if p.next.Type() == token.RightBrace {
			if _, err := p.nextToken(1); err != nil {
				return nil, err
			}
			break
		}
		if _, err := p.nextToken(1); err != nil {
			return nil, err
		}
	}
	if err := p.assertCurrentAndSkip(token.RightBrace); err != nil {
		return nil, err
	}
	return pairs, nil
}
