package ugolang

import (
	"fmt"
)

type Ugolang struct {
	DumpTokens bool
	DumpNodes  bool
}

func NewUgolang() *Ugolang {
	return &Ugolang{}
}

var tokens []Token

func (u *Ugolang) Exec(code string) int {
	tokens = tokenize(code)
	if u.DumpTokens {
		fmt.Printf("tokens: %v\n", tokens)
	}
	nodes := prog()
	if u.DumpNodes {
		fmt.Printf("nodes: %v\n", nodes)
	}
	ret := 0
	for _, node := range nodes {
		eval(&node)
		dprintf("node=%v\n", node)
		ret = eval(&node)
	}
	return ret
}

func eval(node *Node) int {
	switch node.Type {
	case NodeNum:
		return node.Val
	case NodeAdd:
		l := eval(node.Lhs)
		r := eval(node.Rhs)
		return l + r
	case NodeMul:
		l := eval(node.Lhs)
		r := eval(node.Rhs)
		return l * r
	case NodeEq:
		l := eval(node.Lhs)
		r := eval(node.Rhs)
		if l == r {
			return 1
		} else {
			return 0
		}
	case NodeNe:
		l := eval(node.Lhs)
		r := eval(node.Rhs)
		if l != r {
			return 1
		} else {
			return 0
		}
	case NodeLe:
		l := eval(node.Lhs)
		r := eval(node.Rhs)
		if l <= r {
			return 1
		} else {
			return 0
		}
	case NodeLt:
		l := eval(node.Lhs)
		r := eval(node.Rhs)
		if l < r {
			return 1
		} else {
			return 0
		}
	case NodeAssign:
		val := eval(node.Rhs)
		vars.Set(node.Lhs.Ident, val)
		return val
	case NodeVar:
		return vars.Get(node.Ident)
	case NodeIf:
		cond := eval(node.Cond)
		if cond != 0 {
			return eval(node.Then)
		} else {
			if node.Else != nil {
				return eval(node.Else)
			}
		}
		return 0 // FIXME
	default:
		panic(fmt.Sprintf("unknown type: %d", node.Type))
	}
}
