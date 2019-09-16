package vm

import (
	"errors"
	"fmt"

	"github.com/Salpadding/lua/types/chunk"
	"github.com/Salpadding/lua/types/code"
	"github.com/Salpadding/lua/types/value"
	"github.com/Salpadding/lua/types/value/types"
)

const (
	LuaVersionMajor   = "5"
	LuaVersionMINOR   = "3"
	LuaVersionNUM     = 503
	LuaVersionRelease = "4"

	LuaVersion = "Lua " + LuaVersionMajor + "." + LuaVersionMINOR
	LuaRelease = LuaVersion + "." + LuaVersionRelease
)

type BinaryOperator func(a, b value.Value) (value.Value, bool)

type UnaryOperator func(a value.Value) (value.Value, bool)

var binaryOperators = map[types.ArithmeticOperator]BinaryOperator{
	types.Add:        value.Add,
	types.Sub:        value.Sub,
	types.Mul:        value.Mul,
	types.IDiv:       value.IDiv,
	types.Mod:        value.Mod,
	types.Pow:        value.Pow,
	types.Div:        value.Div,
	types.BitwiseAnd: value.BitwiseAnd,
	types.BitwiseXor: value.BitwiseXor,
	types.BitwiseOr:  value.BitwiseOr,
	types.ShiftLeft:  value.ShiftLeft,
	types.ShiftRight: value.ShiftRight,
}

var unaryOperators = map[types.ArithmeticOperator]UnaryOperator{
	types.UnaryMinus: value.UnaryMinus,
	types.BitwiseNot: value.BitwiseNot,
}

// LuaVM is lua state api implementation
type LuaVM struct {
	*Stack
	proto *chunk.Prototype
	pc    int
}

func (vm *LuaVM) Close() {}

func (vm *LuaVM) GetTop() int {
	return vm.Stack.top
}

func (vm *LuaVM) CheckStack(i int) {
	vm.Stack.Check(i)
}

func (vm *LuaVM) Pop(n int) ([]value.Value, error) {
	return vm.Stack.PopN(n)
}

func (vm *LuaVM) Copy(dst, src int) error {
	return vm.Set(dst, vm.Get(src))
}

func (vm *LuaVM) PushValue(idx int) error {
	return vm.Push(vm.Get(idx))
}

func (vm *LuaVM) Replace(idx int) error {
	v, err := vm.Stack.Pop()
	if err != nil {
		return err
	}
	return vm.Set(idx, v)
}

func (vm *LuaVM) Insert(idx int) error {
	return vm.Rotate(idx, 1)
}

func (vm *LuaVM) Remove(idx int) error {
	if err := vm.Rotate(idx, -1); err != nil {
		return err
	}
	if _, err := vm.Pop(1); err != nil {
		return err
	}
	return nil
}

func (vm *LuaVM) SetTop(idx int) error {
	newTop := vm.Stack.AbsIndex(idx)
	if newTop < 0 {
		return errors.New("stack underflow")
	}

	n := vm.Stack.top - newTop
	if n > 0 {
		_, err := vm.Stack.PopN(n)
		if err != nil {
			return err
		}
		return nil
	}
	if n < 0 {
		for i := 0; i > n; i-- {
			if err := vm.Stack.Push(value.GetNil()); err != nil {
				return err
			}
		}
	}
	return nil
}

func (vm *LuaVM) TypeName(v value.Value) string {
	return v.Type().String()
}

func (vm *LuaVM) Type(idx int) types.Type {
	if !vm.IsValid(idx) {
		return types.None
	}
	return vm.Get(idx).Type()
}

func (vm *LuaVM) Rotate(idx, n int) error {
	t := vm.Stack.top - 1           /* end of stack segment being rotated */
	p := vm.Stack.AbsIndex(idx) - 1 /* start of segment */
	var m int                       /* end of prefix */
	if n >= 0 {
		m = t - n
	} else {
		m = p - n - 1
	}
	if err := vm.Stack.reverse(p, m); err != nil {
		return err
	} /* reverse the prefix with length 'n' */
	if err := vm.Stack.reverse(m+1, t); err != nil {
		return err
	} /* reverse the suffix */
	if err := vm.Stack.reverse(p, t); err != nil {
		return err
	} /* reverse the entire segment */
	return nil
}

func (vm *LuaVM) Arithmetic(op types.ArithmeticOperator) error {
	switch op {
	case types.Add, types.Sub, types.Mul, types.IDiv, types.Mod,
		types.Pow, types.Div, types.BitwiseAnd, types.BitwiseXor,
		types.BitwiseOr, types.ShiftLeft, types.ShiftRight:
		a, err := vm.Stack.Pop()
		if err != nil {
			return err
		}
		b, err := vm.Stack.Pop()
		if err != nil {
			return err
		}
		operator := binaryOperators[op]
		v, ok := operator(a, b)
		if !ok {
			return errInvalidOperand
		}
		return vm.Push(v)
	case types.BitwiseNot, types.UnaryMinus:
		a, err := vm.Stack.Pop()
		if err != nil {
			return err
		}
		operator := unaryOperators[op]
		v, ok := operator(a)
		if !ok {
			return errInvalidOperand
		}
		return vm.Push(v)
	}
	return errInvalidOperand
}

func (vm *LuaVM) Len(idx int) error {
	val := vm.Get(idx)
	v, ok := value.Len(val)
	if !ok {
		return errInvalidOperand
	}
	return vm.Push(v)
}

func (vm *LuaVM) Concat(n int) error {
	if n == 0 {
		return vm.Push(value.String(""))
	}
	if n == 1 {
		return nil
	}
	for i := 1; i < n; i++ {
		s2, ok := vm.Get(-1).ToString()
		s1, ok2 := vm.Get(-2).ToString()
		if !ok || !ok2 {
			return errInvalidOperand
		}
		if _, err := vm.PopN(2); err != nil {
			return err
		}
		if err := vm.Push(value.String(s1 + s2)); err != nil {
			return err
		}
	}
	return nil
}

func (vm *LuaVM) AddPC(i int) {
	vm.pc += i
}

func (vm *LuaVM) GetConst(idx int) error {
	if idx < 0 || idx >= len(vm.proto.Constants) {
		return errIndexOverFlow
	}
	v := vm.proto.Constants[idx]
	return vm.Push(v)
}

func (vm *LuaVM) Fetch() code.Instruction {
	i := vm.proto.Code[vm.pc]
	vm.pc++
	return i
}

func (vm *LuaVM) GetRK(rk int) error {
	if rk > 0xff {
		return vm.GetConst(rk)
	}
	return vm.Push(value.Integer(rk + 1))
}

func(vm *LuaVM) execute() error{
	vm.Stack = NewStack(int(vm.proto.MaxStackSize))
	for{
		ins := &Instruction{Instruction: vm.proto.Code[vm.pc]}
		if ins.Opcode().Type == code.Return {
			break
		}
		if err := ins.execute(vm); err != nil{
			return err
		}
		fmt.Println(vm)
	}
	return nil
}
