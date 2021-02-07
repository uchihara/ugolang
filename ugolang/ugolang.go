package ugolang

import (
	"fmt"
)

// Ugolang dummy
type Ugolang struct {
	DumpTokens bool
	DumpNodes  bool
}

// NewUgolang dummy
func NewUgolang() *Ugolang {
	funcStack.reset()
	return &Ugolang{}
}

var tokens []*Token

// Exec dummy
func (u *Ugolang) Exec(code string) int {
	tokens = tokenize(code)
	if u.DumpTokens {
		fmt.Printf("tokens: %v\n", tokens)
	}
	nodes := prog()
	if u.DumpNodes {
		fmt.Printf("nodes: %v\n", nodes)
	}
	return Eval(nodes)
}
