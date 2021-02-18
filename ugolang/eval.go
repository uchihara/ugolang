package ugolang

import (
	"fmt"
)

// Eval dummy
func Eval(nodes []*Node) *Val {
	InitFuncs()
	funcStack.reset()
	funcStack.push("main")
	var ret *Val
	for _, node := range nodes {
		dprintf("node=%v\n", node)
		var nodeType NodeType
		ret, nodeType = eval(node)
		if nodeType == NodeReturn {
			return ret
		}
	}
	return ret
}

func eval(node *Node) (ret *Val, nodeType NodeType) {
	dprintf("eval start node: %v\n", node)
	switch node.Type {
	case NodeVal:
		ret, nodeType = node.Val, 0
	case NodeDefVar:
		val := NewNumVal(0)
		if node.RHS != nil {
			val, _ = eval(node.RHS)
		}
		funcStack.peek().vars.Set(node.Ident, val)
		ret, nodeType = NewNumVal(0), 0
	case NodeAdd:
		l, _ := eval(node.LHS)
		r, _ := eval(node.RHS)
		ret, nodeType = l.Add(r), 0
	case NodeSub:
		l, _ := eval(node.LHS)
		r, _ := eval(node.RHS)
		ret, nodeType = l.Sub(r), 0
	case NodeMul:
		l, _ := eval(node.LHS)
		r, _ := eval(node.RHS)
		ret, nodeType = l.Mul(r), 0
	case NodeEq:
		l, _ := eval(node.LHS)
		r, _ := eval(node.RHS)
		if l.Eq(r) {
			ret, nodeType = NewNumVal(1), 0
			goto end
		}
		ret, nodeType = NewNumVal(0), 0
	case NodeNe:
		l, _ := eval(node.LHS)
		r, _ := eval(node.RHS)
		if l.Ne(r) {
			ret, nodeType = NewNumVal(1), 0
			goto end
		}
		ret, nodeType = NewNumVal(0), 0
	case NodeLe:
		l, _ := eval(node.LHS)
		r, _ := eval(node.RHS)
		if l.Num <= r.Num {
			ret, nodeType = NewNumVal(1), 0
			goto end
		}
		ret, nodeType = NewNumVal(0), 0
	case NodeLt:
		l, _ := eval(node.LHS)
		r, _ := eval(node.RHS)
		if l.Num < r.Num {
			ret, nodeType = NewNumVal(1), 0
			goto end
		}
		ret, nodeType = NewNumVal(0), 0
	case NodeAssign:
		val, _ := eval(node.RHS)
		funcStack.peek().vars.Set(node.LHS.Ident, val)
		ret, nodeType = val, 0
	case NodeVar:
		// bug check
		if !funcStack.peek().vars.Defined(node.Ident) {
			panic(fmt.Sprintf("fatal: undefined var %s found", node.Ident))
		}
		ret, nodeType = funcStack.peek().vars.Get(node.Ident), 0
	case NodeIf:
		cond, _ := eval(node.Cond)
		if cond.Num != 0 {
			ret, nodeType = eval(node.Then)
			goto end
		}
		if node.Else != nil {
			ret, nodeType = eval(node.Else)
			goto end
		}
		ret, nodeType = NewNumVal(0), 0 // FIXME
	case NodeWhile:
		for {
			cond, _ := eval(node.Cond)
			if cond.Num == 0 {
				break
			}
			ret, nodeType = eval(node.Body)
			switch nodeType {
			case NodeReturn:
				goto end
			case NodeBreak:
				ret, nodeType = NewNumVal(0), 0
				goto end
			case NodeContinue:
				continue
			}
		}
	case NodeFunc:
		funcName := node.Ident
		funcs.Define(funcName, node.Args, node.ValType, node.Body)
		ret, nodeType = NewNumVal(0), 0
	case NodeCall:
		funcName := node.Ident
		fn, ok := funcs.Lookup(funcName)
		if !ok {
			panic(fmt.Sprintf("call %s but is not defined", funcName))
		}
		vals := make([]*Val, 0)
		for _, param := range node.Params {
			r, _ := eval(param)
			vals = append(vals, r)
		}
		if fn.IsNative {
			return fn.CallNative(funcName, vals), 0
		}
		funcStack.push(funcName)
		fp := funcStack.peek()
		for i, val := range vals {
			argName := fn.Args[i]
			fp.vars.Set(argName, val)
		}
		body := fn.Body
		r, _ := eval(body)
		funcStack.pop()
		ret, nodeType = r, 0
	case NodeReturn:
		r, _ := eval(node.Expr)
		ret, nodeType = r, node.Type
	case NodeBreak:
		ret, nodeType = NewNumVal(0), node.Type
	case NodeContinue:
		ret, nodeType = NewNumVal(0), node.Type
	case NodeBlock:
		for _, stmt := range node.Statements {
			ret, nodeType = eval(stmt)
			switch nodeType {
			case NodeReturn:
				goto end
			case NodeBreak:
				goto end
			case NodeContinue:
				goto end
			}
		}
	default:
		panic(fmt.Sprintf("unknown type: %d", node.Type))
	}
end:
	dprintf("eval end,  nodeType: %v, ret: %v, new nodeType: %v\n", node.Type, ret, nodeType)
	return ret, nodeType
}
