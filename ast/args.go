package ast

type Arguments interface {
	arguments()
	String() string
}

type Expressions []Expression

func (e Expressions) arguments() {}

func (e Expressions) String() string {
	return joinComma(e)
}
