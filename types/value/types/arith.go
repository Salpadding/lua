package types

type ArithmeticOperator int

const (
	Unknown ArithmeticOperator = iota
	Add
	Sub
	Mul
	Mod
	Pow
	Div
	IDiv
	BitwiseAnd
	BitwiseOr
	BitwiseXor
	ShiftLeft
	ShiftRight
	UnaryMinus
	BitwiseNot
)
