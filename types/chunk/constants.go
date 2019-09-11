package chunk

const (
	LuaSignature    = "\x1bLua"
	LuaVersion      = 0x53
	LuaFormat       = 0
	LuaData         = "\x19\x93\r\n\x1a\n"
	CIntSize        = 4
	CSizeTSize      = 8
	InstructionSize = 4
	LuaIntegerSize  = 8
	LuaNumberSize   = 8
	LuaCInt         = 0x5678
	LuaCNumber      = 370.5
)
