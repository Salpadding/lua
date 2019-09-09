package parser

import (
	"io"

	"github.com/Salpadding/lua/ast"
	"github.com/Salpadding/lua/lex"
	"github.com/Salpadding/lua/token"
)

/*
exp ::=  nil | false | true | Numeral | LiteralString | ‘...’ | functiondef |
	 prefixexp | tableconstructor | exp binop exp | unop exp
*/
/*
exp   ::= exp12
exp12 ::= exp11 {or exp11}
exp11 ::= exp10 {and exp10}
exp10 ::= exp9 {(‘<’ | ‘>’ | ‘<=’ | ‘>=’ | ‘~=’ | ‘==’) exp9}
exp9  ::= exp8 {‘|’ exp8}
exp8  ::= exp7 {‘~’ exp7}
exp7  ::= exp6 {‘&’ exp6}
exp6  ::= exp5 {(‘<<’ | ‘>>’) exp5}
exp5  ::= exp4 {‘..’ exp4}
exp4  ::= exp3 {(‘+’ | ‘-’) exp3}
exp3  ::= exp2 {(‘*’ | ‘/’ | ‘//’ | ‘%’) exp2}
exp2  ::= exp1 | (‘not’ | ‘#’ | ‘-’ | ‘~’) exp1
exp1  ::= exp0 {‘^’ exp2}
exp0  ::= nil | false | true | Numeral | LiteralString
		| ‘...’ | functiondef | prefixexp1 | tableconstructor
*/

/*
prefix1 ::= prefix0 ‘[’ exp ‘]’
	| prefix0 ‘.’ Name
	| prefix0 [‘:’ Name] args
prefix0 ::= Name
	| ‘(’ exp ‘)’

*/
type Parser struct {
	*lex.Lexer
	current token.Token
	next    token.Token
}

func (p *Parser) nextToken(count int) (token.Token, error) {
	for i := 0; i < count; i++ {
		p.current = p.next
		next, err := p.Lexer.NextToken()
		if err != nil {
			return nil, err
		}
		p.next = next
	}
	return p.current, nil
}

func New(reader io.RuneReader) (*Parser, error) {
	p := &Parser{
		Lexer: lex.New(reader),
	}
	if _, err := p.nextToken(2); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Parser) Parse() (*ast.Block, error) {
	return p.parseBlock()
}

func (p *Parser) parseStatements() ([]ast.Statement, error) {
	var res []ast.Statement
	for !p.isReturnOrKeyword(p.current) {
		s, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		res = append(res, s)
	}
	return res, nil
}

func (p *Parser) parseBlock() (*ast.Block, error) {
	statements, err := p.parseStatements()
	if err != nil {
		return nil, err
	}
	if p.current.Type() != token.Return {
		return &ast.Block{
			Statements: statements,
		}, nil
	}
	re, err := p.parseReturn()
	if err != nil {
		return nil, err
	}
	return &ast.Block{
		Statements: statements,
		Return:     re,
	}, nil
}

func (p *Parser) parseStatement() (ast.Statement, error) {
	switch p.current.Type() {
	case token.Break:
		if _, err := p.nextToken(1); err != nil {
			return nil, err
		}
		return ast.Break("break"), nil
	case token.Semicolon:
		if _, err := p.nextToken(1); err != nil {
			return nil, err
		}
		return ast.Empty(";"), nil
	case token.Label:
		if _, err := p.nextToken(1); err != nil {
			return nil, err
		}
		if p.current.Type() != token.Identifier {
			return nil, errUnexpectedError(p.current)
		}
		id := p.current.String()
		if p.next.Type() != token.Label {
			return nil, errUnexpectedError(p.next)
		}
		if _, err := p.nextToken(2); err != nil {
			return nil, err
		}
		return ast.Label(id), nil
	case token.Goto:
		if p.next.Type() != token.Identifier {
			return nil, errUnexpectedError(p.next)
		}
		id := p.next.String()
		if _, err := p.nextToken(2); err != nil {
			return nil, err
		}
		return ast.Goto(id), nil
	case token.Do:
		return p.parseDoBlockEnd()
	case token.While:
		return p.parseWhile()
	case token.Repeat:
		return p.parseRepeat()
	case token.If:
		return p.parseIf()
	case token.For:
		return p.parseFor()
	case token.Function:
		return p.parseFunction()
	case token.Local:
		if p.next.Type() == token.Function {
			if _, err := p.nextToken(1); err != nil {
				return nil, err
			}
			function, err := p.parseFunction()
			if err != nil {
				return nil, err
			}
			return &ast.LocalFunction{Function: function}, nil
		}
		return p.parseLocalAssign()
	default:
		call, err := p.parsePrefix1()
		if err != nil {
			return nil, err
		}
		c, ok := call.(*ast.FunctionCall)
		if ok {
			return c, nil
		}
		switch call.(type) {
		case ast.Identifier, *ast.TableAccess:
		default:
			return nil, errUnexpectedError(p.current)
		}
		if p.current.Type() == token.Comma {
			if _, err := p.nextToken(1); err != nil {
				return nil, err
			}
			assign, err := p.parseAssign()
			if err != nil {
				return nil, err
			}
			assign.Vars = append([]ast.Expression{call}, assign.Vars...)
			return assign, nil
		}
		if err = p.assertCurrentAndSkip(token.Assign); err != nil {
			return nil, err
		}
		exps, err := p.parseExpressions()
		if err != nil {
			return nil, err
		}
		return &ast.Assign{
			Vars:   []ast.Expression{call},
			Values: exps,
		}, nil
	}
}

func (p *Parser) parseExpression() (ast.Expression, error) {
	return p.parseExp12()
}
