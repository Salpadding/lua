package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Salpadding/lua/common"
)

type Statement interface {
	statement()
	String() string
}

type Block struct {
	Statements []Statement
	Return     *Return
}

func (b Block) statement() {}

func (b *Block) String() string {
	g := common.ToGeneral(b.Statements)
	if b.Return != nil {
		g = append(g, b.Return)
	}
	return common.Join(g, "\n")
}

// ;
type Empty string

func (e Empty) statement() {}

func (e Empty) String() string {
	return string(e)
}

type Break string

func (b Break) statement() {}

func (b Break) String() string {
	return string(b)
}

type Return struct {
	Values []Expression
}

func (r *Return) statement() {}

func (r *Return) String() string {
	res := make([]string, len(r.Values))
	for i := range res {
		res[i] = r.Values[i].String()
	}
	return fmt.Sprintf("return %s", strings.Join(res, ", "))
}

type Label string

func (l Label) statement() {}

func (l Label) String() string {
	return fmt.Sprintf(":: %s ::", string(l))
}

type Goto string

func (g Goto) statement() {}

func (g Goto) String() string {
	return fmt.Sprintf("goto %s", string(g))
}

type While struct {
	Condition Expression
	Body      *Block
}

func (w *While) statement() {}

func (w *While) String() string {
	return fmt.Sprintf("while %s do\n%s\nend", w.Condition.String(), common.Indent(2, w.Body.String()))
}

type Repeat struct {
	Condition Expression
	Body      *Block
}

func (r *Repeat) statement() {}

func (r *Repeat) String() string {
	return fmt.Sprintf("repeat \n %s \nuntil %s", r.Body.String(), r.Condition.String())
}

type LocalAssign struct {
	Identifiers []Identifier
	Values      []Expression
}

func (l *LocalAssign) statement() {}

func (l *LocalAssign) String() string {
	if len(l.Values) == 0 {
		return fmt.Sprintf("local %s", common.JoinComma(l.Identifiers))
	}
	return fmt.Sprintf("local %s = %s", common.JoinComma(l.Identifiers), common.JoinComma(l.Values))
}

type Assign struct {
	Vars   []Expression
	Values []Expression
}

func (a *Assign) statement() {}

func (a *Assign) String() string {
	return fmt.Sprintf("%s = %s", common.JoinComma(a.Vars), common.JoinComma(a.Values))
}

type Function struct {
	Name       Identifier
	Body       *Block
	Parameters []Parameter
}

func (f *Function) expression() {}

func (f *Function) statement() {}

func (f *Function) String() string {
	return fmt.Sprintf("function %s (%s)\n%s\nend\n", f.Name, common.JoinComma(f.Parameters), common.Indent(2, f.Body.String()))
}

type LocalFunction struct {
	*Function
}

type Branch struct {
	Condition Expression
	Body      *Block
}

type If struct {
	Consequence  *Branch
	Alternatives []*Branch
	Else         *Block
}

func (i *If) statement() {}

func (i *If) String() string {
	cons := fmt.Sprintf("if %s then\n%s\n", i.Consequence.Condition.String(), common.Indent(2, i.Consequence.Body.String()))
	buf := bytes.NewBufferString(cons)
	if i.Alternatives != nil {
		for _, a := range i.Alternatives {
			buf.WriteString(fmt.Sprintf("elseif %s then\n%s\n", a.Condition.String(), common.Indent(2, a.Body.String())))
		}
	}
	if i.Else != nil {
		buf.WriteString(fmt.Sprintf("else\n%s\n", common.Indent(2, i.Else.String())))
	}
	buf.WriteString("end")
	return buf.String()
}

type For struct {
	Name  Identifier
	Start Expression
	Stop  Expression
	Step  Expression
	Body  *Block
}

func (f *For) statement() {}

func (f *For) String() string {
	buf := bytes.NewBufferString(fmt.Sprintf("for %s = %s, %s", f.Name, f.Start, f.Stop))
	if f.Step != nil {
		buf.WriteString(", ")
		buf.WriteString(f.Step.String())
	}
	buf.WriteString(" do\n")
	buf.WriteString(common.Indent(2, f.Body.String()))
	buf.WriteString("\nend")
	return buf.String()
}

type ForIn struct {
	NameList    []Identifier
	Expressions Expressions
	Body        *Block
}

func (f *ForIn) statement() {}

func (f *ForIn) String() string {
	return fmt.Sprintf("for %s in %s do\n%s\nend", common.JoinComma(f.NameList), common.JoinComma(f.Expressions), common.Indent(2, f.Body.String()))
}
