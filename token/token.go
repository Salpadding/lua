package token

type Type int

var Operators = map[string]Type{
	"=":   Assign,
	"+":   Plus,
	"-":   Minus,
	"*":   Asterisk,
	`/`:   Divide,
	`//`:  IntegerDivide,
	"%":   Modular,
	".":   Dot,
	"<":   LessThan,
	"<=":  LessThanOrEqual,
	">":   GreaterThan,
	">=":  GreaterThanOrEqual,
	"==":  Equal,
	"~=":  NotEqual,
	"&":   BitwiseAnd,
	"|":   BitwiseOr,
	"~":   Wave,
	"and": LogicalAnd,
	"or":  LogicalOr,
	"not": LogicalNot,
	"^":   Power,
	"<<":  LeftShift,
	">>":  RightShift,
	"..":  Concat,
	"#":   Len,
}

var Delimiters = map[string]Type{
	"...": Varing,
	"::":  Label,
	",":   Comma,
	";":   Semicolon,
	"(":   LeftParenthesis,
	")":   RightParenthesis,
	"{":   LeftBrace,    // {
	"}":   RightBrace,   // }
	"[":   LeftBracket,  // [
	"]":   RightBracket, // ]
	":":   Colon,
}

var Keywords = map[string]Type{
	"break":    Break,
	"do":       Do,
	"else":     Else,
	"elseif":   ElseIf,
	"end":      End,
	"false":    False,
	"for":      For,
	"function": Function,
	"goto":     Goto,
	"if":       If,
	"in":       In,
	"local":    Local,
	"nil":      Nil,
	"repeat":   Repeat,
	"return":   Return,
	"then":     Then,
	"true":     True,
	"until":    Until,
	"while":    While,
}

type Token interface {
	Type() Type
	String() string
	Line() int
	Column() int
}

type Literal struct {
	t       Type
	literal string
	line    int
	column  int
}

func (l *Literal) Type() Type {
	return l.t
}

func (l *Literal) String() string {
	return l.literal
}

func (l *Literal) Line() int {
	return l.line
}

func (l *Literal) Column() int {
	return l.column
}

func NewLiteral(t Type, literal string, line, column int) *Literal {
	return &Literal{
		t:       t,
		literal: literal,
		line:    line,
		column:  column,
	}
}

type Operator struct {
	t      Type
	line   int
	column int
}

type Keyword struct {
	t      Type
	line   int
	column int
}

func NewKeyword(keyword string, line, column int) *Keyword {
	return &Keyword{
		t:      Keywords[keyword],
		line:   line,
		column: column,
	}
}

func (k *Keyword) Type() Type {
	return k.t
}

func (k *Keyword) Line() int {
	return k.line
}

func (k *Keyword) Column() int {
	return k.column
}

func (k *Keyword) String() string {
	for i, v := range Keywords {
		if v == k.t {
			return i
		}
	}
	return ""
}

func NewOperator(op string, line, column int) *Operator {
	return &Operator{
		t:      Operators[op],
		line:   line,
		column: column,
	}
}

func (o *Operator) Type() Type {
	return o.t
}

func (o *Operator) String() string {
	for k, v := range Operators {
		if v == o.t {
			return k
		}
	}
	return ""
}

func (o *Operator) Line() int {
	return o.line
}

func (o *Operator) Column() int {
	return o.column
}

type EOF string

func (e EOF) Type() Type {
	return EndOfFile
}

func (e EOF) String() string {
	return "EOF"
}

func (e EOF) Line() int {
	return 0
}

func (e EOF) Column() int {
	return 0
}


type Delimiter struct{
	t      Type
	line   int
	column int
}

func NewDelimiter(op string, line, column int) *Delimiter {
	return &Delimiter{
		t:      Delimiters[op],
		line:   line,
		column: column,
	}
}

func (d *Delimiter) Type() Type {
	return d.t
}

func (d *Delimiter) String() string {
	for k, v := range Delimiters {
		if v == d.t {
			return k
		}
	}
	return ""
}

func (d *Delimiter) Line() int {
	return d.line
}

func (d *Delimiter) Column() int {
	return d.column
}