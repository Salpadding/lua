package vm

import "errors"

var(
	errInvalidOperand = errors.New("invalid operand found")
	errIndexOverFlow = errors.New("index overflow")
)
