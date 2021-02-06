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

var tokens []Token

// Exec dummy
func (u *Ugolang) Exec(code string) int {
	tokens = tokenize(code)
	if u.DumpTokens {
		fmt.Printf("tokens: %v\n", tokens)
	}
	nodes := prog()
	nodes = append(nodes, *NewCallNode("main", []*Node{}))
	funcStack.push("main")
	if u.DumpNodes {
		fmt.Printf("nodes: %v\n", nodes)
	}
	ret := 0
	for _, node := range nodes {
		dprintf("node=%v\n", node)
		var nodeType NodeType
		ret, nodeType = eval(&node)
		if nodeType == NodeReturn {
			return ret
		}
	}

	return ret
}
