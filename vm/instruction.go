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
	code.UNM:  types.UnaryMinus,
}

var cmpMapping = map[code.Type]types.Comparison{
	code.EQ: types.Equal,
	code.LT: types.LessThan,
	code.LE: types.LessThanOrEqual,
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

func (ins *Instruction) execute(vm *LuaVM) error {
	switch ins.Opcode().Type {
	case code.MOVE:
		return ins.move(vm)
	case code.LOADK:
		return ins.loadK(vm)
	case code.LOADKX:
		return ins.loadKx(vm)
	case code.LOADBOOL:
		return ins.loadBool(vm)
	case code.LOADNIL:
		return ins.loadNil(vm)
	case code.ADD, code.SUB, code.MUL, code.MOD,
		code.POW, code.DIV, code.IDIV, code.BAND,
		code.BOR, code.BXOR, code.SHL, code.SHR:
		return ins.binaryArithmetic(vm)
	case code.BNOT, code.UNM:
		return ins.unaryArithmetic(vm)
	case code.NOT:
		return ins.not(vm)
	case code.LEN:
		return ins.len(vm)
	case code.CONCAT:
		return ins.concat(vm)
	case code.EQ, code.LT, code.LE:
		op := cmpMapping[ins.Opcode().Type]
		return ins.compare(vm, op)
	case code.TEST:
		return ins.test(vm)
	case code.TESTSET:
		return ins.testSet(vm)
	case code.FORLOOP:
		return ins.forLoop(vm)
	case code.FORPREP:
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
	if err := vm.Push(value.Integer(b)); err != nil {
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

func (ins *Instruction) len(vm *LuaVM) error {
	a, b, _ := ins.ABC()
	a++
	b++
	if err := vm.Len(b); err != nil {
		return err
	}
	return vm.Replace(a)
}

func (ins *Instruction) concat(vm *LuaVM) error {
	a, b, c := ins.ABC()
	a++
	b++
	c++
	n := c - b + 1
	vm.CheckStack(n)
	for i := b; i <= c; i++ {
		if err := vm.Push(value.Integer(i)); err != nil {
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

func (ins *Instruction) forPrep(vm *LuaVM) error {
	a, sBx := ins.AsBx()
	a++

	if err := vm.Push(value.Integer(a)); err != nil {
		return err
	}
	if err := vm.Push(value.Integer(a + 2)); err != nil {
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

func (ins *Instruction) forLoop(vm *LuaVM) error {
	a, sBx := ins.AsBx()
	a++

	if err := vm.Push(value.Integer(a + 2)); err != nil {
		return err
	}
	if err := vm.Push(value.Integer(a)); err != nil {
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
