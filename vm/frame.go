package vm

import (
	"errors"

	types2 "github.com/Salpadding/lua/types"
	"github.com/Salpadding/lua/types/code"
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

type BinaryOperator func(a, b types2.Value) (types2.Value, bool)

type UnaryOperator func(a types2.Value) (types2.Value, bool)

var binaryOperators = map[types.ArithmeticOperator]BinaryOperator{
	types.Add:        types2.Add,
	types.Sub:        types2.Sub,
	types.Mul:        types2.Mul,
	types.IDiv:       types2.IDiv,
	types.Mod:        types2.Mod,
	types.Pow:        types2.Pow,
	types.Div:        types2.Div,
	types.BitwiseAnd: types2.BitwiseAnd,
	types.BitwiseXor: types2.BitwiseXor,
	types.BitwiseOr:  types2.BitwiseOr,
	types.ShiftLeft:  types2.ShiftLeft,
	types.ShiftRight: types2.ShiftRight,
}

var unaryOperators = map[types.ArithmeticOperator]UnaryOperator{
	types.UnaryMinus: types2.UnaryMinus,
	types.BitwiseNot: types2.BitwiseNot,
}

// Frame 是函数调用帧
type Frame struct {
	*Register
	proto    *types2.Prototype
	pc       int
	returned []types2.Value
}

func NewFrame(prototype *types2.Prototype) *Frame {
	return &Frame{
		Register: &Register{},
		proto:    prototype,
		pc:       0,
	}
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
			if err := f.Push(types2.GetNil()); err != nil {
				return err
			}
		}
	}
	return nil
}

func (f *Frame) TypeName(v types2.Value) string {
	return v.Type().String()
}

func (f *Frame) Type(idx int) types.Type {
	if !f.IsValid(idx) {
		return types.None
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

func (f *Frame) GetConst(idx int) (types2.Value, error) {
	if idx < 0 || idx >= len(f.proto.Constants) {
		return nil, errIndexOverFlow
	}
	return f.proto.Constants[idx], nil
}

func (f *Frame) Fetch() code.Instruction {
	i := f.proto.Code[f.pc]
	f.pc++
	return i
}

func (f *Frame) GetRK(rk int) (types2.Value, error) {
	if rk > 0xff {
		return f.GetConst(rk & 0xff)
	}
	return f.Get(rk), nil
}

func (f *Frame) execute() ([]types2.Value, error) {
	f.Register = NewRegister(0)
	for {
		ins := &Instruction{Instruction: f.Fetch()}
		if ins.Opcode().Type == code.Return {
			return f.returned, nil
		}
		if err := ins.execute(f); err != nil {
			return nil, err
		}
		//fmt.Printf("%s %s\n", ins.Opcode().Name, f)
	}
}
