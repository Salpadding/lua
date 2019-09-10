package vm

import "github.com/Salpadding/lua/types/code"

type Instruction uint32

func (ins Instruction) OpType() code.Type {
	return code.Type(ins & 0x3f)
}
