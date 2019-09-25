package vm

import (
	"bytes"
	"errors"

	"github.com/Salpadding/lua/types"
	"github.com/Salpadding/lua/types/code"
	"github.com/Salpadding/lua/types/value"
)

const (
	fieldsPerFlush = 50
)

var opMapping = map[code.Type]value.ArithmeticOperator{
	// binary operators
	code.Add:        value.Add,
	code.Sub:        value.Sub,
	code.Mul:        value.Mul,
	code.Div:        value.Div,
	code.IDiv:       value.IDiv,
	code.BitwiseAnd: value.BitwiseAnd,
	code.BitwiseOr:  value.BitwiseOr,
	code.BitwiseXor: value.BitwiseXor,
	code.Pow:        value.Pow,
	code.Mod:        value.Mod,
	code.ShiftLeft:  value.ShiftLeft,
	code.ShiftRight: value.ShiftRight,

	// unary operators
	code.BitwiseNot: value.BitwiseNot,
	code.UnaryMinus: value.UnaryMinus,
}

var cmpMapping = map[code.Type]value.Comparison{
	code.Equal:           value.Equal,
	code.LessThan:        value.LessThan,
	code.LessThanOrEqual: value.LessThanOrEqual,
}

type Instruction struct {
	code.Instruction
}

// R(A) := R(B)
func (ins *Instruction) move(f *Frame) error {
	dst, src, _ := ins.ABC()
	return f.Copy(dst, src)
}

func (ins *Instruction) jmp(vm *Frame) error {
	a, sBx := ins.AsBx()
	vm.AddPC(sBx)
	if a != 0 {
		return errors.New("todo")
	}
	return nil
}

func (ins *Instruction) execute(f *Frame) error {
	_, ok := opMapping[ins.Opcode().Type]
	if ok {
		return ins.arithmetic(f)
	}
	switch ins.Opcode().Type {
	case code.Move:
		return ins.move(f)
	case code.Jmp:
		return ins.jmp(f)
	case code.LoadK:
		return ins.loadK(f)
	case code.LoadKX:
		return ins.loadKx(f)
	case code.LoadBool:
		return ins.loadBool(f)
	case code.LoadNil:
		return ins.loadNil(f)
	case code.LogicalNot:
		return ins.not(f)
	case code.Len:
		return ins.len(f)
	case code.Concat:
		return ins.concat(f)
	case code.Equal, code.LessThan, code.LessThanOrEqual:
		op := cmpMapping[ins.Opcode().Type]
		return ins.compare(f, op)
	case code.Test:
		return ins.test(f)
	case code.TestSet:
		return ins.testSet(f)
	case code.ForLoop:
		return ins.forLoop(f)
	case code.ForPrep:
		return ins.forPrep(f)
	case code.NewTable:
		return ins.newTable(f)
	case code.GetTable:
		return ins.getTable(f)
	case code.SetTable:
		return ins.setTable(f)
	case code.SetList:
		return ins.setList(f)
	case code.Return:
		return ins.iReturn(f)
	case code.Closure:
		return ins.closure(f)
	case code.Call:
		return ins.call(f)
	case code.VarArg:
		return ins.varArgs(f)
	case code.TailCall:
		return ins.call(f)
	case code.Self:
		return ins.self(f)
	default:
		return nil
	}
}

// R(A), R(A+1), ..., R(A+B) := nil
func (ins *Instruction) loadNil(frame *Frame) error {
	a, b, _ := ins.ABC()
	for i := a; i <= a+b; i++ {
		if err := frame.Set(i, types.GetNil()); err != nil {
			return err
		}
	}
	return nil
}

// R(A) := (bool)B; if (C) pc++
func (ins *Instruction) loadBool(vm *Frame) error {
	a, b, c := ins.ABC()
	if err := vm.Set(a, types.Boolean(b != 0)); err != nil {
		return err
	}
	if c != 0 {
		vm.AddPC(1)
	}
	return nil
}

// R(A) := Kst(Bx)
func (ins *Instruction) loadK(vm *Frame) error {
	a, bx := ins.ABx()
	v, err := vm.GetConst(bx)
	if err != nil {
		return err
	}
	return vm.Set(a, v)
}

// R(A) := Kst(extra arg)
func (ins *Instruction) loadKx(vm *Frame) error {
	a, _ := ins.ABx()
	ax := vm.Fetch().Ax()
	v, err := vm.GetConst(ax)
	if err != nil {
		return err
	}
	return vm.Set(a, v)
}

func (ins *Instruction) arithmetic(vm *Frame) error {
	op, _ := opMapping[ins.Opcode().Type]
	_, ok := binaryOperators[op]
	if ok {
		return ins.binaryArithmetic(vm)
	}
	return ins.unaryArithmetic(vm)
}

// R(A) := RK(B) op RK(C)
func (ins *Instruction) binaryArithmetic(vm *Frame) error {
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
func (ins *Instruction) unaryArithmetic(vm *Frame) error {
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
func (ins *Instruction) len(vm *Frame) error {
	a, b, _ := ins.ABC()
	length, ok := types.Len(vm.Get(b))
	if !ok {
		return errInvalidOperand
	}
	return vm.Set(a, length)
}

// R(A) := R(B).. ... ..R(C)
func (ins *Instruction) concat(vm *Frame) error {
	a, b, c := ins.ABC()
	var str bytes.Buffer
	for i := b; i <= c; i++ {
		s, ok := vm.Get(i).ToString()
		if !ok {
			return errInvalidOperand
		}
		str.WriteString(s)
	}
	return vm.Set(a, types.String(str.String()))
}

// if ((RK(B) op RK(C)) ~= A) then pc++
func (ins *Instruction) compare(vm *Frame, comparison value.Comparison) error {
	var (
		cmp value.Comparison
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
	if comparison == value.Equal {
		cmp, ok = types.Equal(v1, v2)
	} else {
		cmp, ok = types.Compare(v1, v2)
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
func (ins *Instruction) not(vm *Frame) error {
	a, b, _ := ins.ABC()
	return vm.Set(a, !vm.Get(b).ToBoolean())
}

// if (R(B) <=> C) then R(A) := R(B) else pc++
func (ins *Instruction) testSet(vm *Frame) error {
	a, b, c := ins.ABC()
	if vm.Get(b).ToBoolean() == (c != 0) {
		return vm.Copy(a, b)
	}
	vm.AddPC(1)
	return nil
}

// if not (R(A) <=> C) then pc++
func (ins *Instruction) test(vm *Frame) error {
	a, _, c := ins.ABC()
	if vm.Get(a).ToBoolean() != (c != 0) {
		vm.AddPC(1)
	}
	return nil
}

// R(A)-=R(A+2); pc+=sBx
func (ins *Instruction) forPrep(vm *Frame) error {
	a, sBx := ins.AsBx()
	v, ok := types.Sub(vm.Get(a), vm.Get(a+2))
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
func (ins *Instruction) forLoop(vm *Frame) error {
	a, sBx := ins.AsBx()
	var expect value.Comparison
	v, ok := types.Add(vm.Get(a), vm.Get(a+2))
	if !ok {
		return errInvalidOperand
	}
	if err := vm.Set(a, v); err != nil {
		return err
	}
	num, _ := vm.Get(a + 2).ToFloat()
	v1, v2 := vm.Get(a), vm.Get(a+1)
	cmp, _ := types.Compare(v1, v2)
	if num >= 0 {
		expect = value.LessThanOrEqual
	} else {
		expect = value.GreaterThanOrEqual
	}
	if expect&cmp != 0 {
		vm.AddPC(sBx)
		return vm.Copy(a+3, a)
	}
	return nil
}

// R(A) := {} (size = B, C)
func (ins *Instruction) newTable(vm *Frame) error {
	a, _, _ := ins.ABC()
	return vm.Set(a, types.NewTable())
}

// R(A) [RK(B)] := RK(C)
func (ins *Instruction) setTable(vm *Frame) error {
	a, b, c := ins.ABC()
	v1, err := vm.GetRK(b)
	if err != nil {
		return err
	}
	v2, err := vm.GetRK(c)
	if err != nil {
		return err
	}
	tb, ok := vm.Get(a).(*types.Table)
	if !ok {
		return errInvalidOperand
	}
	return tb.Set(v1, v2)
}

// R(A)[(C-1)*FPF+i] := R(A+i), 1 <= i <= B
func (ins *Instruction) setList(f *Frame) error {
	a, b, c := ins.ABC()
	if c > 0 {
		c--
	} else {
		c = f.Fetch().Ax()
	}
	if b == 0 {
		b = f.GetTop() - a
	}

	tb, ok := f.Get(a).(*types.Table)
	if !ok {
		return errInvalidOperand
	}
	idx := c * fieldsPerFlush
	for j := 1; j <= b; j++ {
		idx++
		if err := tb.Set(types.Integer(idx), f.Get(a+j)); err != nil {
			return err
		}
	}
	return nil
}

// R(A) := R(B)[RK(C)]
func (ins *Instruction) getTable(vm *Frame) error {
	a, b, c := ins.ABC()
	v, err := vm.GetRK(c)
	if err != nil {
		return err
	}
	tb, ok := vm.Get(b).(*types.Table)
	if !ok {
		return errInvalidOperand
	}
	v, err = tb.Get(v)
	if err != nil {
		return err
	}
	return vm.Set(a, v)
}

// R(A) := closure(KPROTO[Bx])
func (ins *Instruction) closure(f *Frame) error {
	a, bx := ins.ABx()
	proto := f.proto.Prototypes[bx]
	fn := &types.Function{Prototype: proto}
	return f.Set(a, fn)
}

// return R(A), ... ,R(A+B-2)
func (ins *Instruction) iReturn(f *Frame) error {
	a, b, _ := ins.ABC()
	if b == 1 {
		return nil
	}
	f.returned = f.Slice(a, a+b-1)
	return nil
}

// R(A), ... ,R(A+C-2) := R(A)(R(A+1), ... ,R(A+B-1))
func (ins *Instruction) call(f *Frame) error {
	a, b, c := ins.ABC()
	args := f.Slice(a+1, a+b)
	var (
		values []types.Value
		err    error
	)
	switch fn := f.Get(a).(type) {
	case *types.Function:
		newFrame := f.vm.NewFrame(fn.Prototype)
		// 参数传递
		err = newFrame.PushN(int(fn.NumParams), args...)
		if err != nil {
			return err
		}
		if len(args) > int(fn.NumParams) && f.proto.IsVararg {
			newFrame.varArgs = args[fn.NumParams:]
		}
		values, err = newFrame.execute()
	case types.Native:
		values, err = fn(args...)
		if err != nil {
			return err
		}
	default:
		return errInvalidOperand
	}
	if err != nil || len(values) != c-1 {
		return errInvalidOperand
	}
	for i := range values {
		err = f.Set(a+i, values[i])
		if err != nil {
			return err
		}
	}
	return nil
}

// R(A), R(A+1), ..., R(A+B-2) = vararg
func (ins *Instruction) varArgs(f *Frame) error {
	a, b, _ := ins.ABC()
	if b == 1 {
		return nil
	}
	varArgsSize := b - 1
	if varArgsSize < 0 {
		varArgsSize = len(f.varArgs)
	}
	args := f.varArgs
	if varArgsSize < len(f.varArgs) {
		args = f.varArgs[:varArgsSize]
	}
	for i := 0; i < varArgsSize; i++ {
		if i >= len(args) {
			if err := f.Set(a+i, types.GetNil()); err != nil {
				return err
			}
			continue
		}
		if err := f.Set(a+i, args[i]); err != nil {
			return err
		}
	}
	return nil
}

// R(A+1) := R(B); R(A) := R(B)[RK(C)]
func (ins *Instruction) self(f *Frame) error {
	a, b, c := ins.ABC()

	if err := f.Copy(a+1, b); err != nil {
		return err
	}

	tb, ok := f.Get(b).(*types.Table)
	if !ok {
		return errInvalidOperand
	}
	k, err := f.GetRK(c)
	if err != nil {
		return err
	}
	v, err := tb.Get(k)
	if err != nil {
		return err
	}
	return f.Set(a, v)
}

// R(A) := UpValue[B]
func (ins *Instruction) getUpValue(f *Frame) error {
	a, b, _ := ins.ABC()
	return f.Copy(a, f.vm.upValueIndex(b))
}

// UpValue[B] := R(A)
func (ins *Instruction) setUpValue(f *Frame) error {
	a, b, _ := ins.ABC()
	return f.Copy(f.vm.upValueIndex(b), a)
}


// UpValue[A][RK(B)] := RK(C) 
func(ins *Instruction) getTableUpValue(f *Frame) error{
	a, b, c := ins.ABC()
	v, err := f.GetRK(c)
	if err != nil{
		return err
	}
	tb, ok := f.Get(f.vm.upValueIndex(b)).(*types.Table)
	if !ok{
		return errInvalidOperand
	}
	v, err = tb.Get(v)
	if err != nil{
		return err
	}
	return f.Set(a, v)
}
