package vm

import (
	"github.com/Salpadding/lua/types/code"
	"github.com/Salpadding/lua/types/value"
)

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
func(ins *Instruction) loadBool(vm *LuaVM) error{
	a, b, c := ins.ABC()
	a++
	if err := vm.Push(value.Boolean(b != 0)); err != nil{
		return err
	}
	if err := vm.Replace(a); err != nil{
		return err
	}
	if c != 0{
		vm.AddPC(1)
	}
	return nil
}

// R(A) := Kst(Bx)
func(ins *Instruction) loadK(vm *LuaVM) error{
	a, bx := ins.ABx()
	a ++
	if err := vm.GetConst(bx); err != nil{
		return err
	}
	return vm.Replace(a)
}
