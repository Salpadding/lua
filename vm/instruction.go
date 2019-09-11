package vm

import (
	"github.com/Salpadding/lua/types/code"
)

type Instruction struct {
	code.Instruction
}

// R(A) := R(B)
func (ins *Instruction) move(vm *LuaVM) error {
	a, b, _ := ins.ABC()
	a += 1
	b += 1
	return vm.Copy(a, b)
}

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
	return nil
}
