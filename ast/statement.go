package ast

import (
	"fmt"
	"strings"
)

type Statement interface {
	statement()
	String() string
}

type Block struct {
	Statements []Statement
	Return     *Return
}

//func (b Block) statement() {}

func (b *Block) String() string {
	res := make([]string, len(b.Statements))
	for i := range res {
		res[i] = b.Statements[i].String()
	}
	if b.Return != nil {
		res = append(res, b.Return.String())
	}
	return strings.Join(res, "\n")
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
	return fmt.Sprintf("while %s \ndo\n %s \nend", w.Condition.String(), w.Body.String())
}

type Repeat struct {
	Condition Expression
	Body      *Block
}

func(r *Repeat) statement(){}

func(r *Repeat) String() string{
	return fmt.Sprintf("repeat \n %s \nuntil\n %s", r.Body.String(), r.Condition.String())
}

type LocalAssign struct {
	Identifiers []Identifier
	Values      []Expression
}

func (l *LocalAssign) statement() {}

func (l *LocalAssign) String() string {
	ids := make([]Expression, len(l.Identifiers))
	for i := range ids {
		ids[i] = l.Identifiers[i]
	}
	return (&Assign{
		Vars:   ids,
		Values: l.Values,
	}).String()
}

type Assign struct {
	Vars   []Expression
	Values []Expression
}

func (a *Assign) statement() {}

func (a *Assign) String() string {
	vars := make([]string, len(a.Vars))
	values := make([]string, len(a.Values))
	for i := range vars {
		vars[i] = a.Vars[i].String()
	}
	for i := range values {
		values[i] = a.Values[i].String()
	}
	return fmt.Sprintf("%s = %s", strings.Join(vars, ", "), strings.Join(values, ", "))
}

type Function struct {
	Name      string
	Body      Block
	Arguments []Expression
}

type LocalFunction struct {
	*Function
}
