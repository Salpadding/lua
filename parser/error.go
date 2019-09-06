package parser

import (
	"fmt"

	"github.com/Salpadding/lua/token"
)

func errUnexpectedError(tk token.Token) error {
	return fmt.Errorf("unexpected token %s found at line %d, column %d", tk.String(), tk.Line(), tk.Column())
}
