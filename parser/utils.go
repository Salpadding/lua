package parser

import "github.com/Salpadding/lua/token"

func (p *Parser) isReturnOrKeyword(tk token.Token) bool {
	switch tk.Type() {
	case token.Return, token.EndOfFile, token.End, token.Else, token.ElseIf, token.Until:
		return true
	}
	return false
}
