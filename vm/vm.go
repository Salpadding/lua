package vm

import (
	"errors"
	"github.com/Salpadding/lua/types/value"
)

const (
	LuaVersionMajor   = "5"
	LuaVersionMINOR   = "3"
	LuaVersionNUM     = 503
	LuaVersionRelease = "4"

	LuaVersion = "Lua " + LuaVersionMajor + "." + LuaVersionMINOR
	LuaRelease = LuaVersion + "." + LuaVersionRelease
)

/* option for multiple returns in 'lua_pcall' and 'lua_call' */
const LUA_MULTRET = -1

/* minimum Lua stack available to a C function */
const LUA_MINSTACK = 20

/*
** Pseudo-indices
** (-LUAI_MAXSTACK is the minimum valid index; we keep some free empty
** space after that to help overflow detection)
 */
const LUA_REGISTRYINDEX = -LUAI_MAXSTACK - 1000

/* predefined values in the registry */
const LUA_RIDX_MAINTHREAD int64 = 1
const LUA_RIDX_GLOBALS int64 = 2
const LUA_RIDX_LAST = LUA_RIDX_GLOBALS

// lua-5.3.4/src/lvm.c
/* limit for table tag-method chains (to avoid loops) */
const MAXTAGLOOP = 2000

// State is lua state api implementation
type State struct {
	stack *Stack
}

func(s *State) Close(){}

func(s *State) GetTop() int{
	return s.stack.top
}

type Stack struct {
	slots []value.Value
	top   int
	pc    int
	prev  *Stack
}

func (s *Stack) push(val value.Value) error {
	if s.top == len(s.slots) {
		return errors.New("stack over flow")
	}
	s.slots[s.top] = val
	s.top++
	return nil
}

func (s *Stack) pop() (value.Value, error) {
	if s.top == 0 {
		return nil, errors.New("stack underflow")
	}
	s.top--
	return s.slots[s.top], nil
}

func (s *Stack) pushN(values []value.Value, n int) error {
	if n < 0 {
		n = len(values)
	}
	for i := 0; i < n; i++ {
		if i >= len(values) {
			if err := s.push(value.Nil("nil")); err != nil {
				return err
			}
		}
		if err := s.push(values[i]); err != nil {
			return err
		}
	}
	return nil
}

func (s *Stack) popN(n int) ([]value.Value, error) {
	values := make([]value.Value, n)
	var err error
	for i := n - 1; i >= 0; i-- {
		values[i], err = s.pop()
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

func (s *Stack) absIndex(idx int) int {
	// zero or positive or pseudo
	if idx >= 0 || idx <= LUA_REGISTRYINDEX {
		return idx
	}
	// negative
	return idx + self.top + 1
}