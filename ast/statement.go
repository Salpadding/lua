package ast

import (
	"bytes"
	"fmt"
	"reflect"
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
	g := toGeneral(b.Statements)
	if b.Return != nil {
		g = append(g, b.Return)
	}
	return join(g, "\n")
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
	return fmt.Sprintf("local %s = %s", joinComma(l.Identifiers), joinComma(l.Values))
}

type Assign struct {
	Vars   []Expression
	Values []Expression
}

func (a *Assign) statement() {}

func (a *Assign) String() string {
	return fmt.Sprintf("%s = %s", joinComma(a.Vars), joinComma(a.Values))
}

type Function struct {
	Name       Identifier
	Body       *Block
	Parameters []Parameter
}

func(f *Function) statement(){}

func(f *Function) String() string{
	return fmt.Sprintf("function %s (%s)\n%s\nend ", f.Name, joinComma(f.Parameters), f.Body)
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
	cons := fmt.Sprintf("if %s then\n %s\n", i.Consequence.Condition.String(), i.Consequence.Body.String())
	buf := bytes.NewBufferString(cons)
	if i.Alternatives != nil {
		for _, a := range i.Alternatives {
			buf.WriteString(fmt.Sprintf("elseif %s then\n %s\n", a.Condition.String(), a.Body.String()))
		}
	}
	if i.Else != nil {
		buf.WriteString(fmt.Sprintf("else\n %s\nend\n", i.Else.String()))
	}
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
	buf.WriteString(" do \n")
	buf.WriteString(f.Body.String())
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
	return fmt.Sprintf("for %s in %s do\n%s\nend\n", joinComma(f.NameList), joinComma(f.Expressions), f.Body.String())
}

func joinComma(i interface{}) string {
	return join(toGeneral(i), ", ")
}

func join(li []interface{}, sep string) string {
	res := make([]string, len(li))
	for i := range res {
		str, ok := li[i].(fmt.Stringer)
		if !ok {
			return ""
		}
		res[i] = str.String()
	}
	return strings.Join(res, sep)
}
func toGeneral(args interface{}) []interface{} {
	s := reflect.ValueOf(args)
	if s.Kind() != reflect.Slice {
		panic("toGeneral given a non-slice type")
	}
	ret := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}
	return ret
}
