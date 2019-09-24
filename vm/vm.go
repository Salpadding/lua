package vm

import (
	"io"

	"github.com/Salpadding/lua/types/chunk"
)

type LuaVM struct {
	main *Frame
}

func (vm *LuaVM) Load(rd io.Reader) error {
	proto, err := chunk.ReadPrototype(rd)
	if err != nil {
		return err
	}
	vm.main = &Frame{
		Register: &Register{},
		proto:    proto,
		pc:       0,
	}
	return nil
}

func (vm *LuaVM) Execute() error {
	_, err := vm.main.execute()
	if err != nil {
		return err
	}
	return nil
}

