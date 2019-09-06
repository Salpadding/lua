package ast

import (
	"bytes"
	"fmt"
	"github.com/Salpadding/lua/common"
	"github.com/Salpadding/lua/token"
	"strconv"
	"strings"
)

type Expression interface {
	expression()
	String() string
}

type PrefixExpression struct {
	Operator *token.Operator
	Right    Expression
}

func (p *PrefixExpression) expression() {}

func (p *PrefixExpression) String() string {
	return fmt.Sprintf("(%s %s)", p.Operator.String(), p.Right.String())
}

type InfixExpression struct {
	Operator *token.Operator
	Left     Expression
	Right    Expression
}

func (e *InfixExpression) expression() {}

func (e *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", e.Left.String(), e.Operator.String(), e.Right.String())
}

type Number float64

func (n Number) expression() {}

func (n Number) String() string {
	return strconv.FormatFloat(float64(n), 'f', -1, 64)
}

type Nil struct{}

func (n *Nil) expression() {}

func (n *Nil) String() string {
	return "nil"
}

type Boolean bool

func (b Boolean) expression() {}

func (b Boolean) String() string {
	if b {
		return "true"
	}
	return "false"
}

type String string

func (s String) expression() {}

func (s String) String() string {
	return `"` + common.Escape(bytes.NewBufferString(string(s))) + `"`
}

type Identifier string

func (i Identifier) expression() {}

func (i Identifier) String() string {
	return string(i)
}

type Vararg string

func (v Vararg) expression() {}

func (v Vararg) String() string {
	return string(v)
}

type FunctionCall struct {
	Function Expression
	Args     []Expression
}

func (f *FunctionCall) expression() {}

func (f *FunctionCall) String() string {
	args := make([]string, len(f.Args))
	for i := range args {
		args[i] = f.Args[i].String()
	}
	return fmt.Sprintf("( %s ) ( %s )", f.Function.String(), strings.Join(args, ", "))
}
