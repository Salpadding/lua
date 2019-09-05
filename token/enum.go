package token

const (
	EndOfFile Type = iota

	// Identifiers + literals
	Identifier // add, foobar, x, y, ...
	Number     // 1343456
	String

	Assign             // =
	Plus               // +
	Minus              // -
	Asterisk           // *
	Divide             // /
	IntegerDivide      // //
	Modular            // %
	Dot                // .
	LessThan           // <
	LessThanOrEqual    // <=
	GreaterThan        // >
	GreaterThanOrEqual // >=
	Equal              // ==
	NotEqual           // ~=
	BitwiseAnd         // &
	BitwiseOr          // |
	Wave               // ~
	LogicalAnd         // and
	LogicalOr          // or
	LogicalNot         // not
	Power              // ^
	LeftShift          // <<
	RightShift         // >>
	Concat             // ..
	Len // #

	// Delimiters
	Varing           // ...
	Label            // ::
	Comma            // ,
	Semicolon        // ;
	LeftParenthesis  // (
	RightParenthesis // )
	LeftBrace        // {
	RightBrace       // }
	LeftBracket      // [
	RightBracket     // ]
	Colon            // :

	// Keywords
	Break    // break
	Do       // do
	Else     // else
	ElseIf   // elseif
	End      // end
	False    // false
	For      // for
	Function // function
	Goto     // goto
	If       // if
	In       // in
	Local    // local
	Nil      // nil
	Repeat   // repeat
	Return   // return
	Then     // then
	True     // true
	Until    // until
	While    // while
)
