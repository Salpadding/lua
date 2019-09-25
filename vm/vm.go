package vm

import (
	"errors"
	"fmt"
	"io"

	"github.com/Salpadding/lua/types"
)

var natives = map[types.Value]types.Native{
	types.String("print"): func(args ...types.Value) (values []types.Value, e error) {
		for _, v := range args {
			fmt.Println(v)
		}
		return []types.Value{types.GetNil()}, nil
	},
	types.String("fail"): func(args ...types.Value) (values []types.Value, e error) {
		return []types.Value{types.GetNil()}, errors.New("fail")
	},
}

type LuaVM struct {
	main     *Frame       // 主函数帧栈
	registry *types.Table // lua 注册表
}

func (vm *LuaVM) upValueIndex(idx int) int {
	return LuaRegistryIndex - idx
}

func (vm *LuaVM) Load(rd io.Reader) error {
	proto, err := types.ReadPrototype(rd)
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

func (vm *LuaVM) NewFrame(proto *types.Prototype) *Frame {
	return &Frame{
		vm:       vm,
		Register: &Register{},
		proto:    proto,
	}
}
