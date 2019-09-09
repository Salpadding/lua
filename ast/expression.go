package ast

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/Salpadding/lua/common"
	"github.com/Salpadding/lua/token"
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
	return fmt.Sprintf(`"%s"`, common.Escape(bytes.NewBufferString(string(s))))
}

func (s String) arguments() {}

type Identifier string

func (i Identifier) expression() {}

func (i Identifier) String() string {
	return string(i)
}

func (i Identifier) parameter() {}

type Vararg string

func (v Vararg) expression() {}

func (v Vararg) String() string {
	return string(v)
}

func (v Vararg) parameter() {}

type FunctionCall struct {
	Function Expression
	Args     Arguments
	Self     Expression
}

func (f *FunctionCall) statement() {}

func (f *FunctionCall) expression() {}

func (f *FunctionCall) String() string {
	if f.Args == nil {
		f.Args = Expressions{}
	}
	if f.Self == nil {
		return fmt.Sprintf("( %s ) ( %s )", f.Function.String(), f.Args.String())
	}
	return fmt.Sprintf("%s:%s ( %s )", f.Self.String(), f.Function.String(), f.Args.String())
}

type TableAccess struct {
	Left  Expression
	Index Expression
}

func (i *TableAccess) expression() {}

func (i *TableAccess) String() string {
	return fmt.Sprintf("(%s[ %s ])", i.Left.String(), i.Index.String())
}

type Keypair struct {
	Key   Expression
	Value Expression
}

func (k *Keypair) String() string {
	return fmt.Sprintf("%s = %s", k.Key, k.Value)
}

type Table []*Keypair

func (tb Table) expression() {}

func (tb Table) String() string {
	return fmt.Sprintf("{ %s }", joinComma(tb))
}

func (tb Table) arguments() {}
