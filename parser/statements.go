package parser

import (
	"github.com/Salpadding/lua/ast"
	"github.com/Salpadding/lua/token"
)

// 解析赋值语句
func (p *Parser) parseAssign() (ast.Statement, error) {
	var vars []ast.Expression
	var values []ast.Expression
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
		if _, err := p.nextToken(); err != nil {
			return nil, err
		}
	}
	if p.current.Type() != token.Assign {
		return nil, errUnexpectedError(p.current)
	}
	if _, err := p.nextToken(); err != nil {
		return nil, err
	}
	for {
		val, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if p.current.Type() != token.Comma {
			break
		}
		if _, err := p.nextToken(); err != nil {
			return nil, err
		}
		values = append(values, val)
	}
	return &ast.Assign{
		Vars:   vars,
		Values: values,
	}, nil
}

func (p *Parser) parseLocalAssign() (ast.Statement, error) {
	// skip local
	if _, err := p.nextToken(); err != nil {
		return nil, err
	}
	st, err := p.parseAssign()
	if err != nil {
		return nil, err
	}
	as, ok := st.(*ast.Assign)
	if !ok {
		return nil, errUnexpectedError(p.current)
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
