package ugolang

import (
	"fmt"
)

var funcNameStack stackType

func consume(tokenType TokenType) (*Token, bool) {
	if len(tokens) == 0 {
		return nil, false
	}
	if tokenType == tokens[0].Type {
		token := tokens[0]
		tokens = tokens[1:]
		return token, true
	}
	return nil, false
}

func consumeSign(sign string) (*Token, bool) {
	if len(tokens) == 0 {
		return nil, false
	}
	if tokens[0].Type == TokenSign && tokens[0].Sign == sign {
		token := tokens[0]
		tokens = tokens[1:]
		return token, true
	}
	return nil, false
}

func consumeIdent() (*Token, string, bool) {
	if len(tokens) == 0 {
		return nil, "", false
	}
	token := tokens[0]
	if token.Type == TokenIdent {
		tokens = tokens[1:]
		return token, token.Ident, true
	}
	return nil, "", false
}

func expect(tokenType TokenType) error {
	if token, ok := consume(tokenType); !ok {
		return NewCompileError(token.Pos(), fmt.Sprintf("expect %v but got %v", tokenType, token))
	}
	return nil
}

func expectSign(sign string) (*Token, error) {
	token, ok := consumeSign(sign)
	if !ok {
		return nil, NewCompileError(token.Pos(), fmt.Sprintf("expect %s but got %v", sign, token))
	}
	return token, nil
}

func expectValType() (*Token, error) {
	token, ok := consume(TokenValType)
	if !ok {
		return nil, NewCompileError(token.Pos(), fmt.Sprintf("%v expect valType but got %v", caller(), tokens[0]))
	}
	return token, nil
}

func expectIdent() (*Token, string, error) {
	token, ident, ok := consumeIdent()
	if !ok {
		return nil, "", NewCompileError(token.Pos(), fmt.Sprintf("%v expect ident but got %v", caller(), tokens[0]))
	}
	return token, ident, nil
}

// evalValType dummy
func evalValType(node *Node) ValType {
	switch node.Type {
	case NodeVar:
		valType, ok := funcTable[funcNameStack.peek()].Vars.Defined(node.Ident)
		if !ok {
			panic(fmt.Sprintf("invalid state var %s is not defined", node.Ident))
		}
		return valType
	case NodeVal:
		return node.Val.Type
	case NodeAdd, NodeSub, NodeMul:
		l := evalValType(node.LHS)
		r := evalValType(node.RHS)
		if l != r {
			return 0
		}
		return l
	case NodeAssign:
		return evalValType(node.RHS)
	case NodeCall:
		funcType, ok := funcTable.Lookup(node.Ident)
		if !ok {
			panic(fmt.Sprintf("invalid state func %s is not defined", node.Ident))
		}
		return funcType.RetValType
	default:
		return 0
	}
}

func validateValTypes(nodes []*Node) error {
	for _, node := range nodes {
		err := validateValType(node)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateValType(node *Node) (err error) {
	dprintf("validate val type start node: %v\n", node)
	switch node.Type {
	// case NodeVal:
	case NodeDefVar:
		l := node.ValType
		if node.RHS != nil {
			r := evalValType(node.RHS)
			if l != r {
				err = NewCompileError(node.TokenPos, fmt.Sprintf("type mismatched %v and %v on declare", l, r))
				goto end
			}
		}
	case NodeAdd, NodeSub, NodeMul, NodeEq, NodeNe, NodeLe, NodeLt, NodeAssign:
		l := evalValType(node.LHS)
		r := evalValType(node.RHS)
		if l == 0 || r == 0 || l != r {
			err = NewCompileError(node.TokenPos, fmt.Sprintf("type mismatched %v and %v on binnode", l, r))
			goto end
		}
	// case NodeVar:
	// bug check
	case NodeIf:
		err = validateValType(node.Then)
		if err != nil {
			goto end
		}
		if node.Else != nil {
			err := validateValType(node.Else)
			if err != nil {
				return err
			}
		}
	case NodeWhile:
		err = validateValType(node.Body)
		if err != nil {
			goto end
		}
	case NodeFunc:
		funcNameStack.push(node.Ident)
		err = validateValType(node.Body)
		if err != nil {
			goto end
		}
		funcNameStack.pop()
	case NodeCall:
		funcName := node.Ident
		fn, ok := funcTable.Lookup(funcName)
		if !ok {
			err = NewCompileError(node.TokenPos, fmt.Sprintf("call %s but is not defined", funcName))
			goto end
		}
		if len(node.Params) != len(fn.Args) {
			err = NewCompileError(node.TokenPos, fmt.Sprintf("# of args mismatch on call function %s", funcName))
			goto end
		}
		valTypes := make([]ValType, 0)
		for _, param := range node.Params {
			valType := evalValType(param)
			valTypes = append(valTypes, valType)
		}
		funcNameStack.push(funcName)
		for i := range valTypes {
			valType := evalValType(fn.Args[i])
			if valType != valTypes[i] {
				err = NewCompileError(node.TokenPos, fmt.Sprintf("type mismatch %v and %v on call function", valType, valTypes[i]))
				goto end
			}
		}
		err = validateValType(fn.Body)
		if err != nil {
			goto end
		}
		funcNameStack.pop()
	case NodeReturn:
		funcName := funcNameStack.peek()
		funcType, ok := funcTable.Lookup(funcName)
		if !ok {
			panic(fmt.Sprintf("invalid state func %s is not defined", funcName))
		}
		valType := evalValType(node.Expr)
		if funcType.RetValType != valType {
			err = NewCompileError(node.TokenPos, fmt.Sprintf("type mismatch %v and %v on return", funcType.RetValType, valType))
			goto end
		}
	// case NodeBreak:
	// case NodeContinue:
	case NodeBlock:
		for _, stmt := range node.Statements {
			err = validateValType(stmt)
			if err != nil {
				goto end
			}
		}
	default:
		dprintf("do nothing for nodeType: %v\n", node.Type)
	}
end:
	dprintf("validate val type end,  nodeType: %v, err: %v\n", node.Type, err)
	return err
}
func prog() ([]*Node, error) {
	InitFuncs()
	funcNameStack.push(".global")

	nodes := make([]*Node, 0)
	for len(tokens) > 0 {
		node, ok, err := funcStmt()
		if err != nil {
			return nil, err
		}
		if !ok {
			node, err = stmt()
			if err != nil {
				return nil, err
			}
		}
		nodes = append(nodes, node)
	}

	err := validateValTypes(nodes)
	if err != nil {
		return nil, err
	}

	nodes = append(nodes, NewCallNode(nil, "main", []*Node{}))
	return nodes, nil
}

func funcStmt() (node *Node, ok bool, err error) {
	dprintf("func start\n")
	var ident string
	var argss []*Node
	var retValTypeToken *Token
	var blockNode *Node
	token, isFunc := consume(TokenFunc)
	if !isFunc {
		goto end
	}

	_, ident, err = expectIdent()
	if err != nil {
		goto end
	}
	funcNameStack.push(ident)
	funcTable.Define(ident)
	argss, err = args()
	if err != nil {
		goto end
	}
	retValTypeToken, err = expectValType()
	if err != nil {
		goto end
	}

	blockNode, err = block()
	if err != nil {
		goto end
	}
	node = NewFuncNode(token.Pos(), ident, argss, retValTypeToken.ValType, blockNode)
	funcTable.Set(ident, argss, retValTypeToken.ValType, blockNode)
	ok = true

	funcNameStack.pop()
end:
	dprintf("func end\n")
	return node, ok, err
}

func args() (args []*Node, err error) {
	dprintf("args start\n")
	_, err = expectSign("(")
	if err != nil {
		goto end
	}
	args = make([]*Node, 0)
	for len(tokens) > 0 {
		var ident string
		var ok bool
		if _, ok = consumeSign(")"); ok {
			break
		}
		if len(args) != 0 {
			_, err = expectSign(",")
			if err != nil {
				goto end
			}
		}
		var token *Token
		var valTypeToken *Token
		token, ident, err = expectIdent()
		if err != nil {
			goto end
		}
		valTypeToken, err = expectValType()
		if err != nil {
			goto end
		}
		if _, ok := funcTable[funcNameStack.peek()].Vars.DefinedLocally(ident); ok {
			err = NewCompileError(token.Pos(), fmt.Sprintf("redeclared variable found: %s", ident))
			goto end
		}
		funcTable[funcNameStack.peek()].Vars.Define(ident, valTypeToken.ValType)
		args = append(args, NewVarNode(token.Pos(), ident))
	}
end:
	dprintf("args end\n")
	return args, err
}

func params() (params []*Node, err error) {
	dprintf("params start\n")
	params = make([]*Node, 0)
	for len(tokens) > 0 {
		if _, ok := consumeSign(")"); ok {
			break
		}
		if len(params) != 0 {
			_, err = expectSign(",")
			if err != nil {
				goto end
			}
		}
		var exprNode *Node
		exprNode, err = expr()
		if err != nil {
			goto end
		}
		params = append(params, exprNode)
	}
end:
	dprintf("params end\n")
	return params, err
}

func block() (node *Node, err error) {
	dprintf("block start\n")
	var token *Token
	token, err = expectSign("{")
	if err != nil {
		goto end
	}
	node = NewBlockNode(token.Pos())
	for len(tokens) > 0 {
		var stmtNode *Node
		stmtNode, err = stmt()
		if err != nil {
			goto end
		}
		node.Statements = append(node.Statements, stmtNode)
		_, ok := consumeSign("}")
		if ok {
			goto end
		}
	}
	_, err = expectSign("}")
	if err != nil {
		goto end
	}
end:
	dprintf("block end\n")
	return node, err
}

func stmt() (node *Node, err error) {
	dprintf("stmt start\n")
	if token, ok := consume(TokenVar); ok {
		node, err = varStmt(token)
		goto end
	}
	if token, ok := consume(TokenReturn); ok {
		var exprNode *Node
		exprNode, err = expr()
		if err != nil {
			goto end
		}
		node = NewReturnNode(token.Pos(), exprNode)
		err = expect(TokenEOL)
		if err != nil {
			goto end
		}
		goto end
	}
	if token, ok := consume(TokenBreak); ok {
		node = NewNode(token.Pos(), NodeBreak)
		err = expect(TokenEOL)
		if err != nil {
			goto end
		}
		goto end
	}
	if token, ok := consume(TokenContinue); ok {
		node = NewNode(token.Pos(), NodeContinue)
		err = expect(TokenEOL)
		if err != nil {
			goto end
		}
		goto end
	}
	if token, ok := consume(TokenIf); ok {
		node, err = ifStmt(token)
		goto end
	}
	if token, ok := consume(TokenWhile); ok {
		node, err = whileStmt(token)
		goto end
	}
	node, err = expr()
	if err != nil {
		goto end
	}
	err = expect(TokenEOL)
	if err != nil {
		goto end
	}
end:
	dprintf("stmt end\n")
	return node, err
}

func expr() (node *Node, err error) {
	dprintf("expr start\n")
	node, err = assign()
	dprintf("expr end\n")
	return node, err
}

func assign() (node *Node, err error) {
	dprintf("assign start\n")
	node, err = eq()
	if err != nil {
		goto end
	}
	if token, ok := consumeSign("="); ok {
		if node.Type != NodeVar {
			err = NewCompileError(token.Pos(), fmt.Sprintf("expect lhs var but got %v", node.Type))
			goto end
		}

		_, ok := funcTable[funcNameStack.peek()].Vars.Defined(node.Ident)
		if !ok {
			err = NewCompileError(token.Pos(), fmt.Sprintf("undefined variable found: %s", node.Ident))
			goto end
		}
		var assignNode *Node
		assignNode, err = assign()
		if err != nil {
			goto end
		}
		node = NewBinNode(token.Pos(), NodeAssign, node, assignNode)

	}
end:
	dprintf("assign end\n")
	return node, err
}

func eq() (node *Node, err error) {
	dprintf("eq start\n")
	node, err = rel()
	if err != nil {
		goto end
	}
	for len(tokens) > 0 {
		var rhs *Node
		if token, ok := consumeSign("=="); ok {
			rhs, err = rel()
			if err != nil {
				goto end
			}
			node = NewBinNode(token.Pos(), NodeEq, node, rhs)
		} else if token, ok := consumeSign("!="); ok {
			rhs, err = rel()
			if err != nil {
				goto end
			}
			node = NewBinNode(token.Pos(), NodeNe, node, rhs)
		} else {
			break
		}
	}
end:
	dprintf("eq end\n")
	return node, err
}

func rel() (node *Node, err error) {
	dprintf("rel start\n")
	node, err = add()
	if err != nil {
		goto end
	}
	for len(tokens) > 0 {
		var lhs, rhs *Node
		if token, ok := consumeSign("<="); ok {
			rhs, err = rel()
			if err != nil {
				goto end
			}

			node = NewBinNode(token.Pos(), NodeLe, node, rhs)
		} else if token, ok := consumeSign("<"); ok {
			rhs, err = rel()
			if err != nil {
				goto end
			}
			node = NewBinNode(token.Pos(), NodeLt, node, rhs)
		} else if token, ok := consumeSign(">="); ok {
			lhs, err = rel()
			if err != nil {
				goto end
			}
			node = NewBinNode(token.Pos(), NodeLe, lhs, node)
		} else if token, ok := consumeSign(">"); ok {
			lhs, err = rel()
			if err != nil {
				goto end
			}
			node = NewBinNode(token.Pos(), NodeLt, lhs, node)
		} else {
			break
		}
	}
end:
	dprintf("rel end\n")
	return node, err
}

func add() (node *Node, err error) {
	dprintf("add start\n")
	node, err = mul()
	if err != nil {
		goto end
	}
	dprintf("add lhs: %v\n", node)
	for len(tokens) > 0 {
		var rhs *Node
		if token, ok := consumeSign("+"); ok {
			rhs, err = mul()
			if err != nil {
				goto end
			}
			node = NewBinNode(token.Pos(), NodeAdd, node, rhs)
			dprintf("add rhs: %v\n", node)
		} else if token, ok := consumeSign("-"); ok {
			rhs, err = mul()
			if err != nil {
				goto end
			}
			node = NewBinNode(token.Pos(), NodeSub, node, rhs)
			dprintf("sub rhs: %v\n", node)
		} else {
			break
		}
	}
end:
	dprintf("add end\n")
	return node, err
}

func mul() (node *Node, err error) {
	dprintf("mul start\n")
	var rhs *Node
	node, err = pri()
	if err != nil {
		goto end
	}
	dprintf("mul lhs: %v\n", node)
	for len(tokens) > 0 {
		if token, ok := consumeSign("*"); ok {
			rhs, err = pri()
			if err != nil {
				goto end
			}
			node = NewBinNode(token.Pos(), NodeMul, node, rhs)
			dprintf("mul rhs: %v\n", node)
		} else {
			break
		}
	}
end:
	dprintf("mul end\n")
	return node, err
}

func pri() (node *Node, err error) {
	dprintf("pri start\n")
	var token *Token
	var ident string
	var ok bool
	var paramss []*Node
	if token, ok = consumeSign("("); ok {
		node, err = expr()
		if err != nil {
			goto end
		}
		_, err = expectSign(")")
		goto end
	}

	token, ident, ok = consumeIdent()
	if ok {
		if token, ok = consumeSign("("); ok {
			paramss, err = params()
			if err != nil {
				goto end
			}
			node = NewCallNode(token.Pos(), ident, paramss)
			goto end
		}

		if _, ok := funcTable[funcNameStack.peek()].Vars.Defined(ident); !ok {
			err = NewCompileError(token.Pos(), fmt.Sprintf("undefined variable found: %s", ident))
			goto end
		}
		node = NewVarNode(token.Pos(), ident)
		goto end
	}

	node, err = val()
	if err != nil {
		goto end
	}
end:
	dprintf("pri end\n")
	return node, err
}

func val() (node *Node, err error) {
	dprintf("val start\n")
	var token *Token
	if len(tokens) == 0 {
		err = NewCompileError(nil, fmt.Sprintf("expect val but no more tokens"))
		goto end
	}
	token = tokens[0]
	if token.Type != TokenVal {
		err = NewCompileError(token.Pos(), fmt.Sprintf("expect val but no got %v", token))
		goto end
	}
	tokens = tokens[1:]
	node = NewValNode(token.Pos(), token.Val)
end:
	dprintf("val end\n")
	return node, err
}

func varStmt(varToken *Token) (node *Node, err error) {
	dprintf("var start\n")
	var ident string
	var rhs *Node
	var valTypeToken *Token
	_ = valTypeToken // FIXME
	_, ident, err = expectIdent()
	if err != nil {
		goto end
	}
	valTypeToken, err = expectValType()
	if err != nil {
		goto end
	}
	if _, ok := consumeSign("="); ok {
		rhs, err = expr()
		if err != nil {
			goto end
		}
	}

	node = NewDefVarNode(varToken.Pos(), ident, valTypeToken.ValType, rhs)

	err = expect(TokenEOL)
	if err != nil {
		goto end
	}
	funcTable[funcNameStack.peek()].Vars.Define(ident, valTypeToken.ValType)
end:
	dprintf("var end\n")
	return node, err
}

func ifStmt(ifToken *Token) (node *Node, err error) {
	dprintf("if start\n")
	var condNode, thenNode, elseNode *Node
	condNode, err = expr()
	if err != nil {
		goto end
	}
	thenNode, err = block()
	if err != nil {
		goto end
	}
	if _, ok := consume(TokenElse); ok {
		elseNode, err = block()
		if err != nil {
			goto end
		}
	}
	node = NewIfNode(ifToken.Pos(), condNode, thenNode, elseNode)
end:
	dprintf("if end\n")
	return node, err
}

func whileStmt(whileToken *Token) (node *Node, err error) {
	var condNode, bodyNode *Node
	dprintf("while start\n")
	condNode, err = expr()
	if err != nil {
		goto end
	}
	bodyNode, err = block()
	if err != nil {
		goto end
	}
	node = NewWhileNode(whileToken.Pos(), condNode, bodyNode)
end:
	dprintf("while end\n")
	return node, err
}
