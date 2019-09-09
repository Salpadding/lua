package parser

import (
	"github.com/Salpadding/lua/ast"
	"github.com/Salpadding/lua/token"
)

func (p *Parser) isReturnOrKeyword(tk token.Token) bool {
	switch tk.Type() {
	case token.Return, token.EndOfFile, token.End, token.Else, token.ElseIf, token.Until:
		return true
	}
	return false
}

// 解析表达式列表
func (p *Parser) parseExpressions() (ast.Expressions, error) {
	var values ast.Expressions
	for {
		val, err := p.parseExp12()
		if err != nil {
			return nil, err
		}
		values = append(values, val)
		if p.current.Type() != token.Comma {
			break
		}
		if _, err := p.nextToken(1); err != nil {
			return nil, err
		}
	}
	return values, nil
}

// 解析标识符列表
func (p *Parser) parseIdentifiers() ([]ast.Identifier, error) {
	var vars []ast.Identifier
	for {
		id := p.current.String()
		if err := p.assertCurrentAndSkip(token.Identifier); err != nil {
			return nil, err
		}
		vars = append(vars, ast.Identifier(id))
		if p.current.Type() != token.Comma {
			break
		}
		if _, err := p.nextToken(1); err != nil {
			return nil, err
		}
	}
	return vars, nil
}

// 解析函数参数列表
func (p *Parser) parseArguments() (ast.Arguments, error) {
	current := p.current
	switch current.Type() {
	case token.String:
		if _, err := p.nextToken(1); err != nil {
			return nil, err
		}
		return ast.String(current.(*token.StringLiteral).Literal()), nil
	case token.LeftBrace:
		return p.parseTable()
	case token.LeftParenthesis:
		if _, err := p.nextToken(1); err != nil {
			return nil, err
		}
		res, err := p.parseExpressions()
		if err != nil {
			return nil, err
		}
		if err := p.assertCurrentAndSkip(token.RightParenthesis); err != nil {
			return nil, err
		}
		return res, nil
	default:
		return nil, errUnexpectedError(p.current)
	}
}
