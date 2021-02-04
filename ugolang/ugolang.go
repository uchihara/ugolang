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
		var isReturn bool
		ret, isReturn = eval(&node)
		if isReturn {
			return ret
		}
	}

	return ret
}

func eval(node *Node) (int, bool) {
	switch node.Type {
	case NodeNum:
		return node.Val, false
	case NodeAdd:
		l, _ := eval(node.LHS)
		r, _ := eval(node.RHS)
		return l + r, false
	case NodeMul:
		l, _ := eval(node.LHS)
		r, _ := eval(node.RHS)
		return l * r, false
	case NodeEq:
		l, _ := eval(node.LHS)
		r, _ := eval(node.RHS)
		if l == r {
			return 1, false
		}
		return 0, false
	case NodeNe:
		l, _ := eval(node.LHS)
		r, _ := eval(node.RHS)
		if l != r {
			return 1, false
		}
		return 0, false
	case NodeLe:
		l, _ := eval(node.LHS)
		r, _ := eval(node.RHS)
		if l <= r {
			return 1, false
		}
		return 0, false
	case NodeLt:
		l, _ := eval(node.LHS)
		r, _ := eval(node.RHS)
		if l < r {
			return 1, false
		}
		return 0, false
	case NodeAssign:
		val, _ := eval(node.RHS)
		vars.Set(node.LHS.Ident, val)
		return val, false
	case NodeVar:
		return vars.Get(node.Ident), false
	case NodeIf:
		cond, _ := eval(node.Cond)
		if cond != 0 {
			return eval(node.Then)
		}
		if node.Else != nil {
			return eval(node.Else)
		}
		return 0, false // FIXME
	case NodeWhile:
		r := 0
		for {
			cond, _ := eval(node.Cond)
			if cond == 0 {
				break
			}
			var isReturn bool
			r, isReturn = eval(node.Body)
			if isReturn {
				return r, true
			}
		}
		return r, false
	case NodeFunc:
		funcName := node.Ident
		funcs.Define(funcName, node.Body)
		return 0, false
	case NodeCall:
		funcName := node.Ident
		if !funcs.Defined(funcName) {
			panic(fmt.Sprintf("call %s but is not defined", funcName))
		}
		funcStack.push(funcName)
		body := funcs[funcName]
		r, _ := eval(body)
		funcStack.pop()
		return r, false
	case NodeReturn:
		r, _ := eval(node.Expr)
		return r, true
	case NodeBlock:
		ret := 0
		for _, stmt := range node.Statements {
			var isReturn bool
			ret, isReturn = eval(stmt)
			if isReturn {
				return ret, true
			}
		}
		return ret, false
	default:
		panic(fmt.Sprintf("unknown type: %d", node.Type))
	}
}
