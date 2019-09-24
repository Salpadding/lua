package vm

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/Salpadding/lua/types/value"
)

type Register []value.Value

func NewRegister(cap int) *Register {
	res := make(Register, 0, cap)
	return &res
}

func (r *Register) Push(v value.Value) error {
	*r = append(*r, v)
	return nil
}

func (r *Register) Pop() (value.Value, error) {
	if len(*r) == 0 {
		return nil, errors.New("index overflow")
	}
	last := (*r)[len(*r)-1]
	*r = (*r)[:len(*r)-1]
	return last, nil
}

func (r *Register) PushN(n int, values ...value.Value) error {
	if n < 0 {
		n = len(values)
	}
	for i := 0; i < n; i++ {
		if i >= len(values) {
			if err := r.Push(value.GetNil()); err != nil {
				return err
			}
		}
		if err := r.Push(values[i]); err != nil {
			return err
		}
	}
	return nil
}

func (r *Register) PopN(n int) ([]value.Value, error) {
	values := make([]value.Value, n)
	var err error
	for i := n - 1; i >= 0; i-- {
		values[i], err = r.Pop()
		if err != nil {
			return nil, err
		}
	}
	return values, nil
}

func (r *Register) Check(n int) {
}

func (r *Register) AbsIndex(idx int) int {
	if idx >= 0 {
		return idx
	}
	idx += len(*r)
	return idx
}

func (r *Register) IsValid(idx int) bool {
	idx = r.AbsIndex(idx)
	return idx < len(*r) && idx >= 0
}

func (r *Register) Get(idx int) value.Value {
	idx = r.AbsIndex(idx)
	if r.IsValid(idx) {
		return (*r)[idx]
	}
	return value.GetNil()
}

func (r *Register) Set(idx int, v value.Value) error {
	idx = r.AbsIndex(idx)
	if idx < 0 {
		return errors.New("stack set fail, invalid index")
	}
	for idx >= len(*r) {
		if err := r.Push(value.GetNil()); err != nil {
			return err
		}
	}
	(*r)[idx] = v
	return nil
}

func (r *Register) reverse(from, to int) error {
	if !r.IsValid(from) || !r.IsValid(to) {
		return errors.New("reverse op fail, index overflow")
	}
	for from < to {
		(*r)[from], (*r)[to] = (*r)[to], (*r)[from]
		from++
		to--
	}
	return nil
}

func (r *Register) GetTop() int {
	return len(*r) - 1
}

func (r *Register) String() string {
	var buf bytes.Buffer
	for i := range *r {
		buf.WriteString(fmt.Sprintf("[%s]", (*r)[i]))
	}
	return buf.String()
}

func(r *Register) Slice(start, end int) []value.Value{
	res := make([]value.Value, end - start)
	for i := range res {
		res[i] = r.Get(start + i)
	}
	return res
}
