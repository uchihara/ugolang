package ugolang

import (
	"fmt"
)

// Ugolang dummy
type Ugolang struct {
	Debug      bool
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
func (u *Ugolang) Exec(code string) (*Val, error) {
	debug = u.Debug

	var err error
	tokens, err = tokenize(code)
	if err != nil {
		return nil, err
	}
	if u.DumpTokens {
		fmt.Printf("tokens: %v\n", tokens)
	}
	nodes, err := prog()
	if err != nil {
		return nil, err
	}

	if u.DumpNodes {
		fmt.Printf("nodes: %v\n", nodes)
	}
	return Eval(nodes), nil
}
