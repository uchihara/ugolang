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
	nodes = append(nodes, *NewCallNode("main"))
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
		l := eval(node.LHS)
		r := eval(node.RHS)
		return l + r
	case NodeMul:
		l := eval(node.LHS)
		r := eval(node.RHS)
		return l * r
	case NodeEq:
		l := eval(node.LHS)
		r := eval(node.RHS)
		if l == r {
			return 1
		}
		return 0
	case NodeNe:
		l := eval(node.LHS)
		r := eval(node.RHS)
		if l != r {
			return 1
		}
		return 0
	case NodeLe:
		l := eval(node.LHS)
		r := eval(node.RHS)
		if l <= r {
			return 1
		}
		return 0
	case NodeLt:
		l := eval(node.LHS)
		r := eval(node.RHS)
		if l < r {
			return 1
		}
		return 0
	case NodeAssign:
		val := eval(node.RHS)
		vars.Set(node.LHS.Ident, val)
		return val
	case NodeVar:
		return vars.Get(node.Ident)
	case NodeIf:
		cond := eval(node.Cond)
		if cond != 0 {
			return eval(node.Then)
		}
		if node.Else != nil {
			return eval(node.Else)
		}
		return 0 // FIXME
	case NodeWhile:
		r := 0
		for eval(node.Cond) != 0 {
			r = eval(node.Body)
		}
		return r
	case NodeFunc:
		funcName := node.Ident
		funcs.Define(funcName, node.Body)
		return 0
	case NodeCall:
		funcName := node.Ident
		if !funcs.Defined(funcName) {
			panic(fmt.Sprintf("call %s but is not defined", funcName))
		}
		body := funcs[funcName]
		return eval(body)
	case NodeBlock:
		ret := 0
		for _, stmt := range node.Statements {
			ret = eval(stmt)
		}
		return ret
	default:
		panic(fmt.Sprintf("unknown type: %d", node.Type))
	}
}
