package vm

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/Salpadding/lua/types/value"
)

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

func (s *Stack) String() string {
	var buf bytes.Buffer
	for i := 0; i < s.top; i++ {
		buf.WriteString(fmt.Sprintf("[%s]", s.slots[i]))
	}
	return buf.String()
}
