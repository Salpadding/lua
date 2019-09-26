package vm

import (
	"errors"

	"github.com/Salpadding/lua/types"
	"github.com/Salpadding/lua/types/code"
	"github.com/Salpadding/lua/types/value"
)

const (
	LuaVersionMajor   = "5"
	LuaVersionMINOR   = "3"
	LuaVersionNUM     = 503
	LuaVersionRelease = "4"

	LuaVersion = "Lua " + LuaVersionMajor + "." + LuaVersionMINOR
	LuaRelease = LuaVersion + "." + LuaVersionRelease

	// 伪索引支持
	LuaMaxStack      = 1000000
	LuaRegistryIndex = -LuaMaxStack - 1000
)

type BinaryOperator func(a, b types.Value) (types.Value, bool)

type UnaryOperator func(a types.Value) (types.Value, bool)

var binaryOperators = map[value.ArithmeticOperator]BinaryOperator{
	value.Add:        types.Add,
	value.Sub:        types.Sub,
	value.Mul:        types.Mul,
	value.IDiv:       types.IDiv,
	value.Mod:        types.Mod,
	value.Pow:        types.Pow,
	value.Div:        types.Div,
	value.BitwiseAnd: types.BitwiseAnd,
	value.BitwiseXor: types.BitwiseXor,
	value.BitwiseOr:  types.BitwiseOr,
	value.ShiftLeft:  types.ShiftLeft,
	value.ShiftRight: types.ShiftRight,
}

var unaryOperators = map[value.ArithmeticOperator]UnaryOperator{
	value.UnaryMinus: types.UnaryMinus,
	value.BitwiseNot: types.BitwiseNot,
}

// Frame 是函数调用帧
type Frame struct {
	vm *LuaVM

	*Register // 寄存器
	fn           *types.Function
	openUpValues map[int]*types.ValuePointer
	pc           int
	varArgs      []types.Value
	returned     []types.Value
}

func (f *Frame) Close() {}

func (f *Frame) Copy(dst, src int) error {
	return f.Set(dst, f.Get(src))
}

func (f *Frame) Replace(idx int) error {
	if idx == f.GetTop() {
		return nil
	}
	v, err := f.Pop()
	if err != nil {
		return err
	}
	return f.Set(idx, v)
}

func (f *Frame) Insert(idx int) error {
	return f.Rotate(idx, 1)
}

func (f *Frame) Remove(idx int) error {
	if err := f.Rotate(idx, -1); err != nil {
		return err
	}
	if _, err := f.Pop(); err != nil {
		return err
	}
	return nil
}

func (f *Frame) SetTop(idx int) error {
	newTop := f.AbsIndex(idx)
	if newTop < 0 {
		return errors.New("stack underflow")
	}

	n := f.GetTop() - newTop
	if n > 0 {
		_, err := f.PopN(n)
		if err != nil {
			return err
		}
		return nil
	}
	if n < 0 {
		for i := 0; i < -n; i++ {
			if err := f.Push(types.GetNil()); err != nil {
				return err
			}
		}
	}
	return nil
}

func (f *Frame) TypeName(v types.Value) string {
	return v.Type().String()
}

func (f *Frame) Type(idx int) value.Type {
	if !f.IsValid(idx) {
		return value.None
	}
	return f.Get(idx).Type()
}

func (f *Frame) Rotate(idx, n int) error {
	t := f.GetTop()      /* end of stack segment being rotated */
	p := f.AbsIndex(idx) /* start of segment */
	var m int            /* end of prefix */
	if n >= 0 {
		m = t - n
	} else {
		m = p - n - 1
	}
	if err := f.reverse(p, m); err != nil {
		return err
	} /* reverse the prefix with length 'n' */
	if err := f.reverse(m+1, t); err != nil {
		return err
	} /* reverse the suffix */
	if err := f.reverse(p, t); err != nil {
		return err
	} /* reverse the entire segment */
	return nil
}

func (f *Frame) AddPC(i int) {
	f.pc += i
}

func (f *Frame) GetConst(idx int) (types.Value, error) {
	if idx < 0 || idx >= len(f.fn.Constants) {
		return nil, errIndexOverFlow
	}
	return f.fn.Constants[idx], nil
}

func (f *Frame) Fetch() code.Instruction {
	i := f.fn.Code[f.pc]
	f.pc++
	return i
}

func (f *Frame) GetRK(rk int) (types.Value, error) {
	if rk > 0xff {
		return f.GetConst(rk & 0xff)
	}
	return f.Get(rk), nil
}

func (f *Frame) execute() ([]types.Value, error) {
	for {
		ins := &Instruction{Instruction: f.Fetch()}
		//name := ins.Opcode().Name
		if err := ins.execute(f); err != nil {
			return nil, err
		}
		if ins.Opcode().Type == code.Return {
			return f.returned, nil
		}
		//fmt.Printf("%s %s\n", name, f)
	}
}
