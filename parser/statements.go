package parser

import (
	"errors"

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
	if err := p.assertCurrentAndSkip(token.Assign); err != nil {
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
	if err = p.assertCurrentAndSkip(token.End); err != nil {
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
	if p.current.Type() != token.Do {
		return nil, errUnexpectedError(p.current)
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
	if err = p.assertCurrentAndSkip(token.Until); err != nil {
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

func (p *Parser) parseLocalAssign() (*ast.LocalAssign, error) {
	// skip local
	if _, err := p.nextToken(1); err != nil {
		return nil, err
	}
	ids, err := p.parseIdentifiers()
	if err != nil {
		return nil, err
	}
	if p.current.Type() != token.Assign {
		return &ast.LocalAssign{
			Identifiers: ids,
		}, nil
	}
	if _, err := p.nextToken(1); err != nil {
		return nil, err
	}
	exps, err := p.parseExpressions()
	if err != nil {
		return nil, err
	}
	return &ast.LocalAssign{
		Identifiers: ids,
		Values:      exps,
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

func (p *Parser) assertCurrentAndSkip(t token.Type) error {
	if p.current.Type() != t {
		return errors.New(errUnexpectedError(p.current).Error() + " " + t.String() + " expected")
	}
	if _, err := p.nextToken(1); err != nil {
		return err
	}
	return nil
}

func (p *Parser) parseIf() (*ast.If, error) {
	// skip if
	if _, err := p.nextToken(1); err != nil {
		return nil, err
	}
	cond, err := p.parseExp12()
	if err != nil {
		return nil, err
	}
	if err := p.assertCurrentAndSkip(token.Then); err != nil {
		return nil, err
	}
	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	res := &ast.If{
		Consequence: &ast.Branch{
			Condition: cond,
			Body:      body,
		},
		Alternatives: []*ast.Branch{},
	}
	for p.current.Type() == token.ElseIf {
		if _, err = p.nextToken(1); err != nil {
			return nil, err
		}
		cond, err = p.parseExp12()
		if err != nil {
			return nil, err
		}
		if err := p.assertCurrentAndSkip(token.Then); err != nil {
			return nil, err
		}
		body, err = p.parseBlock()
		if err != nil {
			return nil, err
		}
		res.Alternatives = append(res.Alternatives, &ast.Branch{
			Condition: cond,
			Body:      body,
		})
	}
	if p.current.Type() == token.End {
		if _, err = p.nextToken(1); err != nil {
			return nil, err
		}
		return res, nil
	}
	if err := p.assertCurrentAndSkip(token.Else); err != nil {
		return nil, err
	}
	body, err = p.parseBlock()
	if err != nil {
		return nil, err
	}
	res.Else = body
	if err := p.assertCurrentAndSkip(token.End); err != nil {
		return nil, err
	}
	return res, nil
}

func (p *Parser) parseFor() (ast.Statement, error) {
	// skip for
	if _, err := p.nextToken(1); err != nil {
		return nil, err
	}
	if p.next.Type() == token.Assign {
		return p.parseForNum()
	}
	return p.parseForIn()
}

func (p *Parser) parseForNum() (*ast.For, error) {
	id := p.current.String()
	if err := p.assertCurrentAndSkip(token.Identifier); err != nil {
		return nil, err
	}
	if err := p.assertCurrentAndSkip(token.Assign); err != nil {
		return nil, err
	}
	start, err := p.parseExp12()
	if err != nil {
		return nil, err
	}
	if err = p.assertCurrentAndSkip(token.Comma); err != nil {
		return nil, err
	}
	stop, err := p.parseExp12()
	if err != nil {
		return nil, err
	}
	stmt := &ast.For{
		Name:  ast.Identifier(id),
		Start: start,
		Stop:  stop,
		Step:  nil,
		Body:  nil,
	}
	if p.current.Type() != token.Comma && p.current.Type() != token.Do {
		return nil, errUnexpectedError(p.current)
	}
	if p.current.Type() != token.Comma {
		stmt.Body, err = p.parseDoBlockEnd()
		if err != nil {
			return nil, err
		}
		return stmt, nil
	}
	// skip ,
	if _, err = p.nextToken(1); err != nil {
		return nil, err
	}
	step, err := p.parseExp12()
	if err != nil {
		return nil, err
	}
	stmt.Step = step
	if p.current.Type() != token.Do {
		return nil, errUnexpectedError(p.current)
	}
	stmt.Body, err = p.parseDoBlockEnd()
	if err != nil {
		return nil, err
	}
	return stmt, nil
}

func (p *Parser) parseForIn() (*ast.ForIn, error) {
	namelist, err := p.parseExpressions()
	if err != nil {
		return nil, err
	}
	names := make([]ast.Identifier, len(namelist))
	for i, n := range namelist {
		id, ok := n.(ast.Identifier)
		if !ok {
			return nil, errUnexpectedError(p.current)
		}
		names[i] = id
	}
	if err = p.assertCurrentAndSkip(token.In); err != nil {
		return nil, err
	}
	exps, err := p.parseExpressions()
	if err != nil {
		return nil, err
	}
	if p.current.Type() != token.Do {
		return nil, errUnexpectedError(p.current)
	}
	body, err := p.parseDoBlockEnd()
	if err != nil {
		return nil, err
	}
	return &ast.ForIn{
		NameList:    names,
		Expressions: exps,
		Body:        body,
	}, nil
}

func (p *Parser) parseFunction() (*ast.Function, error) {
	if err := p.assertCurrentAndSkip(token.Function); err != nil {
		return nil, err
	}
	name := ast.Identifier(p.current.String())
	if err := p.assertCurrentAndSkip(token.Identifier); err != nil {
		return nil, err
	}
	parameters, err := p.parseParameters()
	if err != nil {
		return nil, err
	}
	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	if err := p.assertCurrentAndSkip(token.End); err != nil {
		return nil, err
	}
	return &ast.Function{
		Name:       name,
		Body:       body,
		Parameters: parameters,
	}, nil
}

func (p *Parser) parseParameters() ([]ast.Parameter, error) {
	var res []ast.Parameter
	if p.next.Type() == token.RightParenthesis {
		if _, err := p.nextToken(2); err != nil {
			return nil, err
		}
		return res, nil
	}
	if _, err := p.nextToken(1); err != nil {
		return nil, err
	}
	if p.next.Type() == token.Varying {
		if _, err := p.nextToken(2); err != nil {
			return nil, err
		}
		if err := p.assertCurrentAndSkip(token.RightParenthesis); err != nil {
			return nil, err
		}
		return []ast.Parameter{ast.Vararg("...")}, nil
	}
	for {
		id := ast.Identifier(p.current.String())
		if _, err := p.nextToken(1); err != nil {
			return nil, err
		}
		res = append(res, id)
		if p.current.Type() != token.Comma {
			break
		}
		if p.next.Type() == token.Varying {
			res = append(res, ast.Vararg("..."))
			if _, err := p.nextToken(2); err != nil {
				return nil, err
			}
			break
		}
		if _, err := p.nextToken(1); err != nil {
			return nil, err
		}
	}
	if err := p.assertCurrentAndSkip(token.RightParenthesis); err != nil {
		return nil, err
	}
	return res, nil
}
