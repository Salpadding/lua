package parser

import (
	"github.com/Salpadding/lua/ast"
	"github.com/Salpadding/lua/token"
)

// 解析赋值语句
func (p *Parser) parseAssign() (*ast.Assign, error) {
	var vars []ast.Expression
	for {
		variable, err := p.parsePrefix1()
		if err != nil {
			return nil, err
		}
		switch variable.(type) {
		case ast.Identifier, *ast.TableAccess:
		default:
			return nil, errUnexpectedError(p.current)
		}
		vars = append(vars, variable)
		if p.current.Type() != token.Comma {
			break
		}
		if _, err := p.nextToken(1); err != nil {
			return nil, err
		}
	}
	if p.current.Type() != token.Assign {
		return nil, errUnexpectedError(p.current)
	}
	if _, err := p.nextToken(1); err != nil {
		return nil, err
	}
	values, err := p.parseExpressions()
	if err != nil {
		return nil, err
	}
	return &ast.Assign{
		Vars:   vars,
		Values: values,
	}, nil
}

func (p *Parser) parseDoBlockEnd() (*ast.Block, error) {
	// skip do
	if _, err := p.nextToken(1); err != nil {
		return nil, err
	}
	blk, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	if err = p.assertType(p.current, token.End); err != nil {
		return nil, err
	}
	if _, err := p.nextToken(1); err != nil {
		return nil, err
	}
	return blk, nil
}

func (p *Parser) parseWhile() (*ast.While, error) {
	// skip while
	if _, err := p.nextToken(1); err != nil {
		return nil, err
	}
	condition, err := p.parseExp12()
	if err != nil {
		return nil, err
	}
	if err = p.assertType(p.current, token.Do); err != nil {
		return nil, err
	}
	blk, err := p.parseDoBlockEnd()
	if err != nil {
		return nil, err
	}
	return &ast.While{
		Condition: condition,
		Body:      blk,
	}, nil
}

func (p *Parser) parseRepeat() (*ast.Repeat, error) {
	// skip repeat
	if _, err := p.nextToken(1); err != nil {
		return nil, err
	}
	blk, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	if err = p.assertType(p.current, token.Until); err != nil {
		return nil, err
	}
	if _, err := p.nextToken(1); err != nil {
		return nil, err
	}
	cond, err := p.parseExp12()
	if err != nil {
		return nil, err
	}
	return &ast.Repeat{
		Condition: cond,
		Body:      blk,
	}, nil
}

func (p *Parser) parseLocalAssign() (ast.Statement, error) {
	// skip local
	if _, err := p.nextToken(1); err != nil {
		return nil, err
	}
	var ok bool
	as, err := p.parseAssign()
	if err != nil {
		return nil, err
	}
	ids := make([]ast.Identifier, len(as.Vars))
	for i := range ids {
		ids[i], ok = as.Vars[i].(ast.Identifier)
		if !ok {
			return nil, errUnexpectedError(p.current)
		}
	}
	return &ast.LocalAssign{
		Identifiers: ids,
		Values:      as.Values,
	}, nil
}

func (p *Parser) parseReturn() (*ast.Return, error) {
	// skip return
	if _, err := p.nextToken(1); err != nil {
		return nil, err
	}
	switch p.current.Type() {
	case token.EndOfFile, token.End, token.Else, token.ElseIf, token.Until:
		return &ast.Return{}, nil
	case token.Semicolon:
		if _, err := p.nextToken(1); err != nil {
			return nil, err
		}
		return &ast.Return{}, nil
	default:
		exps, err := p.parseExpressions()
		if err != nil {
			return nil, err
		}
		if p.current.Type() != token.Semicolon {
			return &ast.Return{Values: exps}, nil
		}
		if _, err := p.nextToken(1); err != nil {
			return nil, err
		}
		return &ast.Return{Values: exps}, nil
	}
}
