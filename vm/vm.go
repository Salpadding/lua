package vm

import (
	"errors"
	"github.com/Salpadding/lua/types/value"
)

// State is lua state api implementation
type State struct {
	stack *Stack

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
