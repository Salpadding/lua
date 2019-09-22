package vm

import (
	"bytes"
	"errors"

	"github.com/Salpadding/lua/types/code"
	"github.com/Salpadding/lua/types/value"
	"github.com/Salpadding/lua/types/value/types"
)

const (
	fieldsPerFlush = 50
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
	_, ok := opMapping[ins.Opcode().Type]
	if ok {
		return ins.arithmetic(vm)
	}
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
	case code.NewTable:
		return ins.newTable(vm)
	case code.GetTable:
		return ins.getTable(vm)
	case code.SetTable:
		return ins.setTable(vm)
	case code.SetList:
		return ins.setList(vm)
	default:
		return nil
	}
}

// R(A), R(A+1), ..., R(A+B) := nil
func (ins *Instruction) loadNil(vm *LuaVM) error {
	a, b, _ := ins.ABC()
	for i := a; i <= a+b; i++ {
		if err := vm.Set(i, value.GetNil()); err != nil {
			return err
		}
	}
	return nil
}

// R(A) := (bool)B; if (C) pc++
func (ins *Instruction) loadBool(vm *LuaVM) error {
	a, b, c := ins.ABC()
	if err := vm.Set(a, value.Boolean(b != 0)); err != nil {
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
	v, err := vm.GetConst(bx)
	if err != nil {
		return err
	}
	return vm.Set(a, v)
}

// R(A) := Kst(extra arg)
func (ins *Instruction) loadKx(vm *LuaVM) error {
	a, _ := ins.ABx()
	ax := vm.Fetch().Ax()
	v, err := vm.GetConst(ax)
	if err != nil {
		return err
	}
	return vm.Set(a, v)
}

func (ins *Instruction) arithmetic(vm *LuaVM) error {
	op, _ := opMapping[ins.Opcode().Type]
	_, ok := binaryOperators[op]
	if ok {
		return ins.binaryArithmetic(vm)
	}
	return ins.unaryArithmetic(vm)
}

// R(A) := RK(B) op RK(C)
func (ins *Instruction) binaryArithmetic(vm *LuaVM) error {
	a, b, c := ins.ABC()
	op, _ := opMapping[ins.Opcode().Type]
	fn, _ := binaryOperators[op]
	v1, err := vm.GetRK(b)
	if err != nil {
		return err
	}
	v2, err := vm.GetRK(c)
	if err != nil {
		return err
	}
	v, ok := fn(v1, v2)
	if !ok {
		return errInvalidOperand
	}
	return vm.Set(a, v)
}

// R(A) := op R(B)
func (ins *Instruction) unaryArithmetic(vm *LuaVM) error {
	a, b, _ := ins.ABC()
	op, _ := opMapping[ins.Opcode().Type]
	fn, _ := unaryOperators[op]
	v1, err := vm.GetRK(b)
	if err != nil {
		return err
	}
	v, ok := fn(v1)
	if !ok {
		return errInvalidOperand
	}
	return vm.Set(a, v)
}

// R(A) := length of R(B)
func (ins *Instruction) len(vm *LuaVM) error {
	a, b, _ := ins.ABC()
	length, ok := value.Len(vm.Get(b))
	if !ok {
		return errInvalidOperand
	}
	return vm.Set(a, length)
}

// R(A) := R(B).. ... ..R(C)
func (ins *Instruction) concat(vm *LuaVM) error {
	a, b, c := ins.ABC()
	var str bytes.Buffer
	for i := b; i <= c; i++ {
		s, ok := vm.Get(i).ToString()
		if !ok {
			return errInvalidOperand
		}
		str.WriteString(s)
	}
	return vm.Set(a, value.String(str.String()))
}

// if ((RK(B) op RK(C)) ~= A) then pc++
func (ins *Instruction) compare(vm *LuaVM, comparison types.Comparison) error {
	var (
		cmp types.Comparison
		ok  bool
	)
	a, b, c := ins.ABC()
	v1, err := vm.GetRK(b)
	if err != nil {
		return err
	}
	v2, err := vm.GetRK(c)
	if err != nil {
		return err
	}
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
	return nil
}

// R(A) := not R(B)
func (ins *Instruction) not(vm *LuaVM) error {
	a, b, _ := ins.ABC()
	return vm.Set(a, !vm.Get(b).ToBoolean())
}

// if (R(B) <=> C) then R(A) := R(B) else pc++
func (ins *Instruction) testSet(vm *LuaVM) error {
	a, b, c := ins.ABC()
	if vm.Get(b).ToBoolean() == (c != 0) {
		return vm.Copy(a, b)
	}
	vm.AddPC(1)
	return nil
}

// if not (R(A) <=> C) then pc++
func (ins *Instruction) test(vm *LuaVM) error {
	a, _, c := ins.ABC()
	if vm.Get(a).ToBoolean() != (c != 0) {
		vm.AddPC(1)
	}
	return nil
}

// R(A)-=R(A+2); pc+=sBx
func (ins *Instruction) forPrep(vm *LuaVM) error {
	a, sBx := ins.AsBx()
	v, ok := value.Sub(vm.Get(a), vm.Get(a+2))
	if !ok {
		return errInvalidOperand
	}
	if err := vm.Set(a, v); err != nil {
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
	var expect types.Comparison
	v, ok := value.Add(vm.Get(a), vm.Get(a+2))
	if !ok {
		return errInvalidOperand
	}
	if err := vm.Set(a, v); err != nil {
		return err
	}
	num, _ := vm.Get(a + 2).ToFloat()
	v1, v2 := vm.Get(a), vm.Get(a+1)
	cmp, _ := value.Compare(v1, v2)
	if num >= 0 {
		expect = types.LessThanOrEqual
	} else {
		expect = types.GreaterThanOrEqual
	}
	if expect&cmp != 0 {
		vm.AddPC(sBx)
		return vm.Copy(a+3, a)
	}
	return nil
}

// R(A) := {} (size = B, C)
func (ins *Instruction) newTable(vm *LuaVM) error {
	a, _, _ := ins.ABC()
	return vm.Set(a, value.NewTable())
}

// R(A) [RK(B)] := RK(C)
func (ins *Instruction) setTable(vm *LuaVM) error {
	a, b, c := ins.ABC()
	v1, err := vm.GetRK(b)
	if err != nil {
		return err
	}
	v2, err := vm.GetRK(c)
	if err != nil {
		return err
	}
	tb, ok := vm.Get(a).(*value.Table)
	if !ok {
		return errInvalidOperand
	}
	return tb.Set(v1, v2)
}

// R(A)[(C-1)*FPF+i] := R(A+i), 1 <= i <= B
func (ins *Instruction) setList(vm *LuaVM) error {
	a, b, c := ins.ABC()
	if c > 0 {
		c--
	} else {
		c = vm.Fetch().Ax()
	}
	tb, ok := vm.Get(a).(*value.Table)
	if !ok {
		return errInvalidOperand
	}
	idx := c * fieldsPerFlush
	for j := 1; j <= b; j++ {
		idx++
		if err := tb.Set(value.Integer(idx), vm.Get(a+j)); err != nil {
			return err
		}
	}
	return nil
}

// R(A) := R(B)[RK(C)]
func(ins *Instruction) getTable(vm *LuaVM) error{
	a, b, c := ins.ABC()
	v, err := vm.GetRK(c)
	if err != nil{
		return err
	}
	tb, ok := vm.Get(b).(*value.Table)
	if !ok {
		return errInvalidOperand
	}
	v, err = tb.Get(v)
	if err != nil{
		return err
	}
	return vm.Set(a, v)
}


