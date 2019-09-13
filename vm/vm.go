package vm

import (
	"errors"

	"github.com/Salpadding/lua/types/chunk"
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

