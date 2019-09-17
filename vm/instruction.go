package vm

import (
	"errors"

	"github.com/Salpadding/lua/types/code"
	"github.com/Salpadding/lua/types/value"
	"github.com/Salpadding/lua/types/value/types"
)

var opMapping = map[code.Type]types.ArithmeticOperator{
	// binary operators
	code.Add:        types.Add,
	code.Sub:        types.Sub,
	code.Mul:        types.Mul,
	code.Div:        types.Div,
	code.IDiv:       types.IDiv,
	code.BitwiseAnd: types.BitwiseAnd,
	code.BitwiseOr:  types.BitwiseOr,
	code.BitwiseXor: types.BitwiseXor,
	code.Pow:        types.Pow,
	code.Mod:        types.Mod,
	code.ShiftLeft:  types.ShiftLeft,
	code.ShiftRight: types.ShiftRight,

	// unary operators
	code.BitwiseNot: types.BitwiseNot,
	code.UnaryMinus: types.UnaryMinus,
}

var cmpMapping = map[code.Type]types.Comparison{
	code.Equal:           types.Equal,
	code.LessThan:        types.LessThan,
	code.LessThanOrEqual: types.LessThanOrEqual,
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

func (ins *Instruction) jmp(vm *LuaVM) error {
	a, sBx := ins.AsBx()
	vm.AddPC(sBx)
	if a != 0 {
		return errors.New("todo")
	}
	return nil
}

func (ins *Instruction) execute(vm *LuaVM) error {
	switch ins.Opcode().Type {
	case code.Move:
		return ins.move(vm)
	case code.Jmp:
		return ins.jmp(vm)
	case code.LoadK:
		return ins.loadK(vm)
	case code.LoadKX:
		return ins.loadKx(vm)
	case code.LoadBool:
		return ins.loadBool(vm)
	case code.LoadNil:
		return ins.loadNil(vm)
	case code.Add, code.Sub, code.Mul, code.Mod,
		code.Pow, code.Div, code.IDiv, code.BitwiseAnd,
		code.BitwiseOr, code.BitwiseXor, code.ShiftLeft, code.ShiftRight:
		return ins.binaryArithmetic(vm)
	case code.BitwiseNot, code.UnaryMinus:
		return ins.unaryArithmetic(vm)
	case code.LogicalNot:
		return ins.not(vm)
	case code.Len:
		return ins.len(vm)
	case code.Concat:
		return ins.concat(vm)
	case code.Equal, code.LessThan, code.LessThanOrEqual:
		op := cmpMapping[ins.Opcode().Type]
		return ins.compare(vm, op)
	case code.Test:
		return ins.test(vm)
	case code.TestSet:
		return ins.testSet(vm)
	case code.ForLoop:
		return ins.forLoop(vm)
	case code.ForPrep:
		return ins.forPrep(vm)
	default:
		return nil
	}
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

// R(A) := RK(B) op RK(C)
func (ins *Instruction) binaryArithmetic(vm *LuaVM) error {
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

// R(A) := op R(B)
func (ins *Instruction) unaryArithmetic(vm *LuaVM) error {
	a, b, _ := ins.ABC()
	a++
	b++
	if err := vm.Push(vm.Get(b)); err != nil {
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

// R(A) := length of R(B)
func (ins *Instruction) len(vm *LuaVM) error {
	a, b, _ := ins.ABC()
	a++
	b++
	if err := vm.Len(b); err != nil {
		return err
	}
	return vm.Replace(a)
}

// R(A) := R(B).. ... ..R(C)
func (ins *Instruction) concat(vm *LuaVM) error {
	a, b, c := ins.ABC()
	a++
	b++
	c++
	n := c - b + 1
	vm.CheckStack(n)
	for i := b; i <= c; i++ {
		if err := vm.Push(vm.Get(i)); err != nil {
			return err
		}
	}
	if err := vm.Concat(n); err != nil {
		return err
	}
	return vm.Replace(a)
}

// if ((RK(B) op RK(C)) ~= A) then pc++
func (ins *Instruction) compare(vm *LuaVM, comparison types.Comparison) error {
	var (
		cmp types.Comparison
		ok  bool
	)
	a, b, c := ins.ABC()
	if err := vm.GetRK(b); err != nil {
		return err
	}
	if err := vm.GetRK(c); err != nil {
		return err
	}

	v1 := vm.Get(-1)
	v2 := vm.Get(-2)
	if comparison == types.Equal {
		cmp, ok = value.Equal(v1, v2)
	} else {
		cmp, ok = value.Compare(v1, v2)
	}
	if !ok {
		return errInvalidOperand
	}
	ok = cmp&comparison != 0
	if ok != (a != 0) {
		vm.AddPC(1)
	}
	_, err := vm.Pop(2)
	return err
}

// R(A) := not R(B)
func (ins *Instruction) not(vm *LuaVM) error {
	a, b, _ := ins.ABC()
	a++
	b++
	if err := vm.Push(vm.Get(b).ToBoolean()); err != nil {
		return err
	}
	return vm.Replace(a)
}

// if (R(B) <=> C) then R(A) := R(B) else pc++
func (ins *Instruction) testSet(vm *LuaVM) error {
	a, b, c := ins.ABC()
	a++
	b++
	if vm.Get(b).ToBoolean() == (c != 0) {
		return vm.Copy(a, b)
	}
	vm.AddPC(1)
	return nil
}

// if not (R(A) <=> C) then pc++
func (ins *Instruction) test(vm *LuaVM) error {
	a, _, c := ins.ABC()
	a++
	if vm.Get(a).ToBoolean() != (c != 0) {
		vm.AddPC(1)
	}
	return nil
}

// R(A)-=R(A+2); pc+=sBx
func (ins *Instruction) forPrep(vm *LuaVM) error {
	a, sBx := ins.AsBx()
	a++

	if err := vm.Push(vm.Get(a)); err != nil {
		return err
	}
	if err := vm.Push(vm.Get(a+2)); err != nil {
		return err
	}
	if err := vm.Arithmetic(types.Add); err != nil {
		return err
	}
	if err := vm.Replace(a); err != nil {
		return err
	}
	vm.AddPC(sBx)
	return nil
}

// R(A)+=R(A+2);
// if R(A) <?= R(A+1) then {
//   pc+=sBx; R(A+3)=R(A)
// }
func (ins *Instruction) forLoop(vm *LuaVM) error {
	a, sBx := ins.AsBx()
	a++

	if err := vm.Push(vm.Get(a + 2)); err != nil {
		return err
	}
	if err := vm.Push(vm.Get(a)); err != nil {
		return err
	}
	if err := vm.Arithmetic(types.Add); err != nil {
		return err
	}
	if err := vm.Replace(a); err != nil {
		return err
	}
	num, _ := vm.Get(a + 2).ToFloat()
	v1, v2 := vm.Get(a), vm.Get(a+1)
	cmp, _ := value.Compare(v1, v2)
	if num >= 0 && types.LessThanOrEqual&cmp != 0 || num < 0 && types.GreaterThanOrEqual&cmp != 0 {
		vm.AddPC(sBx)
		return vm.Copy(a+3, a)
	}
	return nil
}
