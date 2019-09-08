package ast

import (
	"fmt"
	"strings"
)

type Arguments interface {
	arguments()
	String() string
}

type Expressions []Expression

func (e Expressions) arguments() {}

func (e Expressions) String() string {
	res := make([]string, len(e))
	for i := range res {
		res[i] = e[i].String()
	}
	return fmt.Sprintf("( %s )", strings.Join(res, ", "))
}
