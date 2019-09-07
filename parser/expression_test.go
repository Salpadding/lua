package parser

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/Salpadding/lua/ast"
	"github.com/Salpadding/lua/token"
)

func TestParse0(t *testing.T) {
	p, err := New(bytes.NewBufferString(`
	12 "abc" true false nil
`))
	if err != nil {
		t.Error(err)
	}
	for p.current.Type() != token.EndOfFile {
		exp, err := p.parseExp1()
		if err != nil {
			t.Error(err)
		}
		_, ok := exp.(ast.Expression)
		if !ok {
			t.Fail()
		}
		fmt.Println(exp.String())
	}
}

func TestParse1(t *testing.T) {
	p, err := New(bytes.NewBufferString(`
	1 ^ 2 ^ 3
`))
	if err != nil {
		t.Error(err)
	}
	for p.current.Type() != token.EndOfFile {
		exp, err := p.parseExp1()
		if err != nil {
			t.Error(err)
		}
		_, ok := exp.(ast.Expression)
		if !ok {
			t.Fail()
		}
		fmt.Println(exp.String())
	}
}

func TestParse2(t *testing.T) {
	p, err := New(bytes.NewBufferString(`
	not not true
`))
	if err != nil {
		t.Error(err)
	}
	for p.current.Type() != token.EndOfFile {
		exp, err := p.parseExp2()
		if err != nil {
			t.Error(err)
		}
		_, ok := exp.(ast.Expression)
		if !ok {
			t.Fail()
		}
		fmt.Println(exp.String())
	}
}

func TestParse3(t *testing.T) {
	p, err := New(bytes.NewBufferString(`
	1^200^2*3000 / 100 // 3 % 24 
`))
	if err != nil {
		t.Error(err)
	}
	for p.current.Type() != token.EndOfFile {
		exp, err := p.parseExp3()
		if err != nil {
			t.Error(err)
		}
		_, ok := exp.(ast.Expression)
		if !ok {
			t.Fail()
		}
		fmt.Println(exp.String())
	}
}

func TestParse4(t *testing.T) {
	p, err := New(bytes.NewBufferString(`
	1^200^2*3000 / 100 + 1 - 300 - -3
`))
	if err != nil {
		t.Error(err)
	}
	for p.current.Type() != token.EndOfFile {
		exp, err := p.parseExp4()
		if err != nil {
			t.Error(err)
		}
		_, ok := exp.(ast.Expression)
		if !ok {
			t.Fail()
		}
		fmt.Println(exp.String())
	}
}

func TestParse5(t *testing.T) {
	p, err := New(bytes.NewBufferString(`
	1^200^2 .. 1 .. "abc"
`))
	if err != nil {
		t.Error(err)
	}
	for p.current.Type() != token.EndOfFile {
		exp, err := p.parseExp5()
		if err != nil {
			t.Error(err)
		}
		_, ok := exp.(ast.Expression)
		if !ok {
			t.Fail()
		}
		fmt.Println(exp.String())
	}
}

func TestParse6(t *testing.T) {
	p, err := New(bytes.NewBufferString(`
	2 << 1 + 2 +3 .. "abc"
`))
	if err != nil {
		t.Error(err)
	}
	for p.current.Type() != token.EndOfFile {
		exp, err := p.parseExp6()
		if err != nil {
			t.Error(err)
		}
		_, ok := exp.(ast.Expression)
		if !ok {
			t.Fail()
		}
		fmt.Println(exp.String())
	}
}

func TestParse7(t *testing.T) {
	p, err := New(bytes.NewBufferString(`
	12 & 13 + 14
`))
	if err != nil {
		t.Error(err)
	}
	for p.current.Type() != token.EndOfFile {
		exp, err := p.parseExp7()
		if err != nil {
			t.Error(err)
		}
		_, ok := exp.(ast.Expression)
		if !ok {
			t.Fail()
		}
		fmt.Println(exp.String())
	}
}

func TestParse8(t *testing.T) {
	p, err := New(bytes.NewBufferString(`
	16 ~ ~ 12 & 13 + 14
`))
	if err != nil {
		t.Error(err)
	}
	for p.current.Type() != token.EndOfFile {
		exp, err := p.parseExp8()
		if err != nil {
			t.Error(err)
		}
		_, ok := exp.(ast.Expression)
		if !ok {
			t.Fail()
		}
		fmt.Println(exp.String())
	}
}

func TestParse9(t *testing.T) {
	p, err := New(bytes.NewBufferString(`
	true & false | 1 & 3
`))
	if err != nil {
		t.Error(err)
	}
	for p.current.Type() != token.EndOfFile {
		exp, err := p.parseExp9()
		if err != nil {
			t.Error(err)
		}
		_, ok := exp.(ast.Expression)
		if !ok {
			t.Fail()
		}
		fmt.Println(exp.String())
	}
}

func TestParse10(t *testing.T) {
	p, err := New(bytes.NewBufferString(`
	1 < 3 + 2 2 <= 4 + 9 1 == 1 2 ~= 3  2 > 9 * 99
`))
	if err != nil {
		t.Error(err)
	}
	for p.current.Type() != token.EndOfFile {
		exp, err := p.parseExp10()
		if err != nil {
			t.Error(err)
		}
		_, ok := exp.(ast.Expression)
		if !ok {
			t.Fail()
		}
		fmt.Println(exp.String())
	}
}

func TestParse11(t *testing.T) {
	p, err := New(bytes.NewBufferString(`
	1 < 3 and 1 <= 1000 and false
`))
	if err != nil {
		t.Error(err)
	}
	for p.current.Type() != token.EndOfFile {
		exp, err := p.parseExp11()
		if err != nil {
			t.Error(err)
		}
		_, ok := exp.(ast.Expression)
		if !ok {
			t.Fail()
		}
		fmt.Println(exp.String())
	}
}

func TestParse12(t *testing.T) {
	p, err := New(bytes.NewBufferString(`
	1 < 3 and 1 <= 1000 and false or true and false
`))
	if err != nil {
		t.Error(err)
	}
	for p.current.Type() != token.EndOfFile {
		exp, err := p.parseExp12()
		if err != nil {
			t.Error(err)
		}
		_, ok := exp.(ast.Expression)
		if !ok {
			t.Fail()
		}
		fmt.Println(exp.String())
	}
}

func TestParseGrouped(t *testing.T) {
	p, err := New(bytes.NewBufferString(`
	( 1 +  - 2 ^ 16 - 13 * 2 / 8) * 3
`))
	if err != nil {
		t.Error(err)
	}
	for p.current.Type() != token.EndOfFile {
		exp, err := p.parseExpression()
		if err != nil {
			t.Error(err)
		}
		_, ok := exp.(*ast.InfixExpression)
		if !ok {
			t.Fail()
		}
		fmt.Println(exp.String())
	}
}

func TestParseIndex(t *testing.T) {
	p, err := New(bytes.NewBufferString(`
	identifier["abc"]
	"abcddeff".length
`))
	if err != nil {
		t.Error(err)
	}
	for p.current.Type() != token.EndOfFile {
		exp, err := p.parseExpression()
		if err != nil {
			t.Error(err)
		}
		_, ok := exp.(*ast.TableAccess)
		if !ok {
			t.Fail()
		}
		fmt.Println(exp.String())
	}
}

func TestParseFunctionCall(t *testing.T) {
	p, err := New(bytes.NewBufferString(`
	(12).add(1, 2, 3, 100, ...).name()
	call("arg", ...).isOk.assert("true")
	abscd.len()
`))
	if err != nil {
		t.Error(err)
	}
	for p.current.Type() != token.EndOfFile {
		exp, err := p.parseExpression()
		if err != nil {
			t.Error(err)
		}
		_, ok := exp.(*ast.FunctionCall)
		if !ok {
			t.Fail()
		}
		fmt.Println(exp.String())
	}
}
