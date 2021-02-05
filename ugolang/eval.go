package ugolang

import (
	"fmt"
)

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
		funcStack.peek().locals.Set(node.LHS.Ident, val)
		return val, false
	case NodeVar:
		return funcStack.peek().locals.Get(node.Ident), false
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
		funcs.Define(funcName, node.Args, node.Body)
		return 0, false
	case NodeCall:
		funcName := node.Ident
		fn, ok := funcs.Lookup(funcName)
		if !ok {
			panic(fmt.Sprintf("call %s but is not defined", funcName))
		}
		vals := make([]int, 0)
		for _, param := range node.Params {
			r, _ := eval(param)
			vals = append(vals, r)
		}
		funcStack.push(funcName)
		fp := funcStack.peek()
		for i, val := range vals {
			argName := fn.Args[i]
			fp.locals[argName] = val
		}
		body := fn.Body
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
