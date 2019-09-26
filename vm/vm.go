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
		return []types.Value{types.GetNil()}, errors.New("assertion fail")
	},
}

type LuaVM struct {
	main     *Frame       // 主函数帧栈
	registry *types.Table // lua 注册表
	global   *types.Table // 全局变量
}

func (vm *LuaVM) Load(rd io.Reader) error {
	proto, err := types.ReadPrototype(rd)
	if err != nil {
		return err
	}
	vm.main = &Frame{
		Register: &Register{},
		fn: &types.Function{
			Prototype: proto,
			UpValues:  make([]*types.ValuePointer, len(proto.UpValues)),
		},
		pc: 0,
		vm: vm,
	}
	for i := range vm.main.fn.UpValues{
		vm.main.fn.UpValues[i] = &types.ValuePointer{Value: types.GetNil()}
	}
	vm.registry = types.NewTable()
	// global
	vm.global = types.NewTable()
	for k, v := range natives {
		if err = vm.global.Set(k, v); err != nil {
			return err
		}
	}
	if err = vm.registry.Set(types.String("_ENV"), vm.global); err != nil {
		return err
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

func (vm *LuaVM) NewFrame(fn *types.Function) *Frame {
	return &Frame{
		vm:       vm,
		Register: &Register{},
		fn:       fn,
	}
}
