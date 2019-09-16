package types

type Comparison int

const(
	LessThan Comparison = 1 << iota
	Equal
	GreaterThan
	LessThanOrEqual = LessThan | Equal
	GreaterThanOrEqual = GreaterThan | Equal
)
