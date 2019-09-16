package types

type Type int

const (
	None Type = iota
	Nil
	Boolean
	Number
	String
	Table
	Function
	Thread
	UserData
)

var m = map[Type]string{
	None:     "none",
	Nil:      "nil",
	Boolean:  "boolean",
	String:   "string",
	Table:    "table",
	Function: "function",
	Thread:   "thread",
	UserData: "userdata",
}

func (t Type) String() string {
	return m[t]
}
