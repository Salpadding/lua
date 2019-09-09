package parser

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/Salpadding/lua/ast"
)

func TestParseAssign(t *testing.T) {
	p, err := New(bytes.NewBufferString(`
	a, a["bb"], a.name = "aa", 1 + 333, call("a", "b")
`))
	if err != nil {
		t.Error(err)
	}
	s, err := p.parseAssign()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(s.String())
}

func TestParseLocalAssign(t *testing.T) {
	p, err := New(bytes.NewBufferString(`
	local a, ab, ce = "aa", 1 + 333, call("a", "b")
`))
	if err != nil {
		t.Error(err)
	}
	s, err := p.parseLocalAssign()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(s.String())
}

func TestParseSimples(t *testing.T) {
	p, err := New(bytes.NewBufferString(`
	;
	break
	:: label ::
	goto label
`))
	if err != nil {
		t.Error(err)
	}
	stmts, err := p.parseStatements()
	if err != nil {
		t.Error(err)
	}
	for _, st := range stmts {
		fmt.Println(st.String())
	}
}

func TestParseBlock(t *testing.T) {
	p, err := New(bytes.NewBufferString(`
	a = 1 + 2
	b = 1 + a
	return a + b, 1 ,2 ;
`))
	if err != nil {
		t.Error(err)
	}
	blk, err := p.parseBlock()
	if err != nil {
		t.Error(err)
	}
	for _, st := range blk.Statements {
		fmt.Println(st.String())
	}
	if blk.Return != nil {
		fmt.Println(blk.Return.String())
	}
}

func TestParseWhile(t *testing.T) {
	p, err := New(bytes.NewBufferString(`
	while a + 3 < 100
	do
		a = 1 + 2
		b = 1 + a
		return a + b, 1 ,2 ;
	end

`))
	if err != nil {
		t.Error(err)
	}
	blk, err := p.parseWhile()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(blk.String())
}

func TestParseRepeat(t *testing.T) {
	p, err := New(bytes.NewBufferString(`
	repeat
		a = 1 + 2
		b = 1 + a
		return a + b, 1 ,2 ;
	until true
`))
	if err != nil {
		t.Error(err)
	}
	blk, err := p.parseRepeat()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(blk.String())
}

func TestIf(t *testing.T) {
	p, err := New(bytes.NewBufferString(`
	if a < 10 then
		a = a + 1
	elseif b == 100 then
		a = a + 100
	else
		a = b + 1000
	end
`))
	if err != nil {
		t.Error(err)
	}
	blk, err := p.parseIf()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(blk.String())
}

func TestIf1(t *testing.T) {
	p, err := New(bytes.NewBufferString(`
	if a < 10 then
		a = a + 1
	elseif b == 100 then
		a = a + 100
	end
`))
	if err != nil {
		t.Error(err)
	}
	blk, err := p.parseIf()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(blk.String())
}

func TestIf2(t *testing.T) {
	p, err := New(bytes.NewBufferString(`
	if a < 10 then
		a = a + 1
	end
`))
	if err != nil {
		t.Error(err)
	}
	blk, err := p.parseIf()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(blk.String())
}

func TestAAA(t *testing.T) {
	joinList := func(li []interface{}, sep string) (string, bool) {
		res := make([]string, len(li))
		for i := range res {
			str, ok := li[i].(fmt.Stringer)
			if !ok {
				return "", ok
			}
			res[i] = str.String()
		}
		return strings.Join(res, sep), true
	}
	toGeneral := func(args interface{}) []interface{} {
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
	in := []ast.Identifier{"aaa", "bbb", "ccc"}
	str, ok := joinList(toGeneral(in), ", ")
	if !ok {
		t.Fail()
	}
	fmt.Println(str)
}

func TestFor1(t *testing.T) {
	p, err := New(bytes.NewBufferString(`
	for i = 1 + 1, 10, 1 * 2 do
		a = 1 + 2
		b = 1 + a
	end
`))
	if err != nil {
		t.Error(err)
	}
	blk, err := p.parseFor()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(blk.String())
}

func TestFor2(t *testing.T) {
	p, err := New(bytes.NewBufferString(`
	for k, v in {12, k = "v", [133] = "addff"} do
		a = 1 + 2
		b = 1 + a
	end
`))
	if err != nil {
		t.Error(err)
	}
	blk, err := p.parseFor()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(blk.String())
}

func TestFunction1(t *testing.T) {
	p, err := New(bytes.NewBufferString(`
	function main (a, b, c, ...)
		c = 2
		return a + b + c	
	end
`))
	if err != nil {
		t.Error(err)
	}
	blk, err := p.parseFunction()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(blk.String())
}

func TestFunction2(t *testing.T) {
	p, err := New(bytes.NewBufferString(`
	function main (...)
		c = 2
		return a + b + c	
	end
`))
	if err != nil {
		t.Error(err)
	}
	blk, err := p.parseFunction()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(blk.String())
}

func TestFunction3(t *testing.T) {
	p, err := New(bytes.NewBufferString(`
	function main ()
		c = 2
		return a + b + c	
	end
`))
	if err != nil {
		t.Error(err)
	}
	blk, err := p.parseFunction()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(blk.String())
}