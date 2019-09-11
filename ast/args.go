package ast

import "github.com/Salpadding/lua/common"

type Arguments interface {
	arguments()
	String() string
}

type Expressions []Expression

func (e Expressions) arguments() {}

func (e Expressions) String() string {
	return common.JoinComma(e)
}
