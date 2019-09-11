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

type Stack struct {
	slots []value.Value
	top   int
	pc    int
	prev  *Stack
}

func (s *Stack) Push(val value.Value) error {
	if s.top == len(s.slots) {
		return errors.New("stack over flow")
	}
	s.slots[s.top] = val
	s.top++
	return nil
}

func (s *Stack) PushNil() error {
	return s.Push(value.GetNil())
}

func (s *Stack) Pop() (value.Value, error) {
	if s.top == 0 {
		return nil, errors.New("stack underflow")
	}
	s.top--
	return s.slots[s.top], nil
}

func (s *Stack) PushN(values []value.Value, n int) error {
	if n < 0 {
		n = len(values)
	}
	for i := 0; i < n; i++ {
		if i >= len(values) {
			if err := s.Push(value.GetNil()); err != nil {
				return err
			}
		}
		if err := s.Push(values[i]); err != nil {
			return err
		}
	}
	return nil
}

func (s *Stack) PopN(n int) ([]value.Value, error) {
	values := make([]value.Value, n)
	var err error
	for i := n - 1; i >= 0; i-- {
		values[i], err = s.Pop()
		if err != nil {
			return nil, err
		}
	}
	return values, nil
}

func NewStack(size int) *Stack {
	return &Stack{
		slots: make([]value.Value, size),
		top:   0,
		pc:    0,
	}
}

func (s *Stack) Check(n int) {
	free := len(s.slots) - s.top
	if free >= n {
		return
	}
	slots := make([]value.Value, len(s.slots)+n)
	copy(slots, s.slots)
	s.slots = slots
}

func (s *Stack) AbsIndex(idx int) int {
	if idx >= 0 {
		return idx
	}
	idx += s.top
	idx++
	return idx
}

func (s *Stack) IsValid(idx int) bool {
	idx = s.AbsIndex(idx)
	return idx <= s.top && idx > 0
}

func (s *Stack) Get(idx int) value.Value {
	idx = s.AbsIndex(idx)
	if s.IsValid(idx) {
		return s.slots[idx-1]
	}
	return value.GetNil()
}

func (s *Stack) Set(idx int, val value.Value) error {
	idx = s.AbsIndex(idx)
	if s.IsValid(idx) {
		s.slots[idx-1] = val
		return nil
	}
	return errors.New("stack set fail, invalid index")
}

func (s *Stack) reverse(from, to int) error {
	if !s.IsValid(from) || !s.IsValid(to) {
		return errors.New("reverse op fail, index overflow")
	}
	slots := s.slots
	for from < to {
		slots[from], slots[to] = slots[to], slots[from]
		from++
		to--
	}
	return nil
}
