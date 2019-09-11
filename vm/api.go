package vm

type LuaState interface{
	GetTop() int
	AbsIndex(int) int
	CheckStack(int)
}