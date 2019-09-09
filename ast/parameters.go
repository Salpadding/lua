package ast

type Parameter interface {
	Expression
	parameter()
}
