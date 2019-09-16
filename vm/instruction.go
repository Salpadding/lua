package vm

import (
	"github.com/Salpadding/lua/types/code"
	"github.com/Salpadding/lua/types/value"
	"github.com/Salpadding/lua/types/value/types"
)

var opMapping = map[code.Type]types.ArithmeticOperator{
	// binary operators
	code.ADD:  types.Add,
	code.SUB:  types.Sub,
	code.MUL:  types.Mul,
	code.DIV:  types.Div,
	code.IDIV: types.IDiv,
	code.BAND: types.BitwiseAnd,
	code.BOR:  types.BitwiseOr,
	code.BXOR: types.BitwiseXor,
	code.POW:  types.Pow,
	code.MOD:  types.Mod,
	code.SHL:  types.ShiftLeft,
	code.SHR:  types.ShiftRight,

	// unary operators
	code.BNOT: types.BitwiseNot,
	code.UNM: types.UnaryMinus,
}

type Instruction struct {
	code.Instruction
}

// R(A) := R(B)
func (ins *Instruction) move(vm *LuaVM) error {
	dst, src, _ := ins.ABC()
	src += 1
	dst += 1
	return vm.Copy(dst, src)
}

// R(A), R(A+1), ..., R(A+B) := nil
func (ins *Instruction) loadNil(vm *LuaVM) error {
	a, b, _ := ins.ABC()
	a += 1
	if err := vm.PushNil(); err != nil {
		return err
	}
	for i := a; i <= a+b; i++ {
		if err := vm.Copy(i, -1); err != nil {
			return err
		}
	}
	if _, err := vm.Stack.Pop(); err != nil {
		return err
	}
	return nil
}

// R(A) := (bool)B; if (C) pc++
func (ins *Instruction) loadBool(vm *LuaVM) error {
	a, b, c := ins.ABC()
	a++
	if err := vm.Push(value.Boolean(b != 0)); err != nil {
		return err
	}
	if err := vm.Replace(a); err != nil {
		return err
	}
	if c != 0 {
		vm.AddPC(1)
	}
	return nil
}

// R(A) := Kst(Bx)
func (ins *Instruction) loadK(vm *LuaVM) error {
	a, bx := ins.ABx()
	a++
	if err := vm.GetConst(bx); err != nil {
		return err
	}
	return vm.Replace(a)
}

// R(A) := Kst(extra arg)
func (ins *Instruction) loadKx(vm *LuaVM) error {
	a, _ := ins.ABx()
	a++
	ax := vm.Fetch().Ax()
	if err := vm.GetConst(ax); err != nil {
		return err
	}
	return vm.Replace(a)
}

func (ins *Instruction) binaryArith(vm *LuaVM) error {
	a, b, c := ins.ABC()
	a++
	if err := vm.GetRK(b); err != nil {
		return err
	}
	if err := vm.GetRK(c); err != nil {
		return err
	}
	op, ok := opMapping[ins.Opcode().Type]
	if !ok {
		return errInvalidOperand
	}
	if err := vm.Arithmetic(op); err != nil {
		return err
	}
	return vm.Replace(a)
}

func (ins *Instruction) unaryArith(vm *LuaVM) error {
	a, b, _ := ins.ABC()
	a ++
	b ++
	if err := vm.Push(value.Integer(b)); err != nil{
		return err
	}
	op, ok := opMapping[ins.Opcode().Type]
	if !ok {
		return errInvalidOperand
	}
	if err := vm.Arithmetic(op); err != nil{
		return err
	}
	return vm.Replace(a)
}

func(ins *Instruction) len(vm *LuaVM) error{
	a, b, _ := ins.ABC()
	a++
	b++
	if err := vm.Len(b); err != nil{
		return err
	}
	return vm.Replace(a)
}

func(ins *Instruction) concat(vm *LuaVM) error{
	a, b, c := ins.ABC()
	a ++
	b++
	c++
	n := c-b+1
	vm.CheckStack(n)
	for i := b; i <= c; i++{
		if err := vm.Push(value.Integer(i)); err != nil{
			return err
		}
	}
	if err := vm.Concat(n); err != nil{
		return err
	}
	return vm.Replace(a)
}