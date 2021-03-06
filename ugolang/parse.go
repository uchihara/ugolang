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

// guessValType dummy
func guessValType(node *Node) ValType {
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
		l := guessValType(node.LHS)
		r := guessValType(node.RHS)
		if l != r {
			return 0
		}
		return l
	case NodeAssign:
		return guessValType(node.RHS)
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
	defer dprintf("validate val type end,  nodeType: %v, err: %v\n", node.Type, err)

	switch node.Type {
	// case NodeVal:
	case NodeDefVar:
		l := node.ValType
		if node.RHS != nil {
			r := guessValType(node.RHS)
			if l != r {
				err = NewCompileError(node.TokenPos, fmt.Sprintf("type mismatched %v and %v on declare", l, r))
				if err != nil {
					return
				}
			}
		}
	case NodeAdd, NodeSub, NodeMul, NodeEq, NodeNe, NodeLe, NodeLt, NodeAssign:
		l := guessValType(node.LHS)
		r := guessValType(node.RHS)
		if l == 0 || r == 0 || l != r {
			err = NewCompileError(node.TokenPos, fmt.Sprintf("type mismatched %v and %v on binnode", l, r))
			if err != nil {
				return
			}
		}
	// case NodeVar:
	// bug check
	case NodeIf:
		err = validateValType(node.Then)
		if err != nil {
			return
		}
		if node.Else != nil {
			err = validateValType(node.Else)
			if err != nil {
				return
			}
		}
	case NodeWhile:
		err = validateValType(node.Body)
		if err != nil {
			return
		}
	case NodeFunc:
		funcNameStack.push(node.Ident)
		err = validateValType(node.Body)
		if err != nil {
			return
		}
		funcNameStack.pop()
	case NodeCall:
		funcName := node.Ident
		fn, ok := funcTable.Lookup(funcName)
		if !ok {
			err = NewCompileError(node.TokenPos, fmt.Sprintf("call %s but is not defined", funcName))
			return
		}
		if len(node.Params) != len(fn.Args) {
			err = NewCompileError(node.TokenPos, fmt.Sprintf("# of args mismatch on call function %s", funcName))
			return
		}
		valTypes := make([]ValType, 0)
		for _, param := range node.Params {
			valType := guessValType(param)
			valTypes = append(valTypes, valType)
		}
		funcNameStack.push(funcName)
		for i := range valTypes {
			valType := guessValType(fn.Args[i])
			if valType != valTypes[i] {
				err = NewCompileError(node.TokenPos, fmt.Sprintf("type mismatch %v and %v on call function", valType, valTypes[i]))
				return
			}
		}
		err = validateValType(fn.Body)
		if err != nil {
			return
		}
		funcNameStack.pop()
	case NodeReturn:
		funcName := funcNameStack.peek()
		funcType, ok := funcTable.Lookup(funcName)
		if !ok {
			panic(fmt.Sprintf("invalid state func %s is not defined", funcName))
		}
		valType := guessValType(node.Expr)
		if funcType.RetValType != valType {
			err = NewCompileError(node.TokenPos, fmt.Sprintf("type mismatch %v and %v on return", funcType.RetValType, valType))
			return
		}
	// case NodeBreak:
	// case NodeContinue:
	case NodeBlock:
		for _, stmt := range node.Statements {
			err = validateValType(stmt)
			if err != nil {
				return
			}
		}
	default:
		dprintf("do nothing for nodeType: %v\n", node.Type)
	}
	return
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

func funcStmt() (*Node, bool, error) {
	dprintf("func start\n")
	defer dprintf("func end\n")

	token, isFunc := consume(TokenFunc)
	if !isFunc {
		return nil, false, nil
	}

	_, ident, err := expectIdent()
	if err != nil {
		return nil, false, err
	}
	funcNameStack.push(ident)
	funcTable.Define(ident)
	argss, err := args()
	if err != nil {
		return nil, false, err
	}
	retValTypeToken, err := expectValType()
	if err != nil {
		return nil, false, err
	}

	blockNode, err := block()
	if err != nil {
		return nil, false, err
	}
	dprintf("blockNode: %v\n", blockNode)
	node := NewFuncNode(token.Pos(), ident, argss, retValTypeToken.ValType, blockNode)
	funcTable.Set(ident, argss, retValTypeToken.ValType, blockNode)

	funcNameStack.pop()
	return node, true, nil
}

func args() ([]*Node, error) {
	dprintf("args start\n")
	defer dprintf("args end\n")

	_, err := expectSign("(")
	if err != nil {
		return nil, err
	}
	args := make([]*Node, 0)
	for len(tokens) > 0 {
		if _, ok := consumeSign(")"); ok {
			break
		}
		if len(args) != 0 {
			_, err := expectSign(",")
			if err != nil {
				return nil, err
			}
		}
		token, ident, err := expectIdent()
		if err != nil {
			return nil, err
		}
		valTypeToken, err := expectValType()
		if err != nil {
			return nil, err
		}
		if _, ok := funcTable[funcNameStack.peek()].Vars.DefinedLocally(ident); ok {
			return nil, NewCompileError(token.Pos(), fmt.Sprintf("redeclared variable found: %s", ident))
		}
		funcTable[funcNameStack.peek()].Vars.Define(ident, valTypeToken.ValType)
		args = append(args, NewVarNode(token.Pos(), ident))
	}
	return args, err
}

func params() ([]*Node, error) {
	dprintf("params start\n")
	defer dprintf("params end\n")

	params := make([]*Node, 0)
	for len(tokens) > 0 {
		if _, ok := consumeSign(")"); ok {
			break
		}
		if len(params) != 0 {
			_, err := expectSign(",")
			if err != nil {
				return nil, err
			}
		}
		exprNode, err := expr()
		if err != nil {
			return nil, err
		}
		params = append(params, exprNode)
	}
	return params, nil
}

func block() (*Node, error) {
	dprintf("block start\n")
	defer dprintf("block end\n")

	token, err := expectSign("{")
	if err != nil {
		return nil, err
	}
	node := NewBlockNode(token.Pos())
	for len(tokens) > 0 {
		stmtNode, err := stmt()
		if err != nil {
			return nil, err
		}
		node.Statements = append(node.Statements, stmtNode)
		_, ok := consumeSign("}")
		if ok {
			return node, nil
		}
	}
	_, err = expectSign("}")
	if err != nil {
		return nil, err
	}
	return node, nil
}

func stmt() (*Node, error) {
	dprintf("stmt start\n")
	defer dprintf("stmt end\n")

	if token, ok := consume(TokenVar); ok {
		return varStmt(token)
	}
	if token, ok := consume(TokenReturn); ok {
		exprNode, err := expr()
		if err != nil {
			return nil, err
		}
		node := NewReturnNode(token.Pos(), exprNode)
		err = expect(TokenEOL)
		if err != nil {
			return nil, err
		}
		return node, nil
	}
	if token, ok := consume(TokenBreak); ok {
		node := NewNode(token.Pos(), NodeBreak)
		err := expect(TokenEOL)
		if err != nil {
			return nil, err
		}
		return node, nil
	}
	if token, ok := consume(TokenContinue); ok {
		node := NewNode(token.Pos(), NodeContinue)
		err := expect(TokenEOL)
		if err != nil {
			return nil, err
		}
		return node, nil
	}
	if token, ok := consume(TokenIf); ok {
		return ifStmt(token)
	}
	if token, ok := consume(TokenWhile); ok {
		return whileStmt(token)
	}
	node, err := expr()
	if err != nil {
		return nil, err
	}
	err = expect(TokenEOL)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func expr() (*Node, error) {
	dprintf("expr start\n")
	defer dprintf("expr end\n")
	return assign()
}

func assign() (*Node, error) {
	dprintf("assign start\n")
	defer dprintf("assign end\n")

	node, err := eq()
	if err != nil {
		return nil, err
	}
	if token, ok := consumeSign("="); ok {
		if node.Type != NodeVar {
			return nil, NewCompileError(token.Pos(), fmt.Sprintf("expect lhs var but got %v", node.Type))
		}

		_, ok := funcTable[funcNameStack.peek()].Vars.Defined(node.Ident)
		if !ok {
			return nil, NewCompileError(token.Pos(), fmt.Sprintf("undefined variable found: %s", node.Ident))
		}
		assignNode, err := assign()
		if err != nil {
			return nil, err
		}
		node = NewBinNode(token.Pos(), NodeAssign, node, assignNode)

	}
	return node, nil
}

func eq() (*Node, error) {
	dprintf("eq start\n")
	defer dprintf("eq end\n")

	node, err := rel()
	if err != nil {
		return nil, err
	}
	for len(tokens) > 0 {
		if token, ok := consumeSign("=="); ok {
			rhs, err := rel()
			if err != nil {
				return nil, err
			}
			node = NewBinNode(token.Pos(), NodeEq, node, rhs)
		} else if token, ok := consumeSign("!="); ok {
			rhs, err := rel()
			if err != nil {
				return nil, err
			}
			node = NewBinNode(token.Pos(), NodeNe, node, rhs)
		} else {
			break
		}
	}
	return node, nil
}

func rel() (*Node, error) {
	dprintf("rel start\n")
	defer dprintf("rel end\n")

	node, err := add()
	if err != nil {
		return nil, err
	}
	for len(tokens) > 0 {
		if token, ok := consumeSign("<="); ok {
			rhs, err := rel()
			if err != nil {
				return nil, err
			}

			node = NewBinNode(token.Pos(), NodeLe, node, rhs)
		} else if token, ok := consumeSign("<"); ok {
			rhs, err := rel()
			if err != nil {
				return nil, err
			}
			node = NewBinNode(token.Pos(), NodeLt, node, rhs)
		} else if token, ok := consumeSign(">="); ok {
			lhs, err := rel()
			if err != nil {
				return nil, err
			}
			node = NewBinNode(token.Pos(), NodeLe, lhs, node)
		} else if token, ok := consumeSign(">"); ok {
			lhs, err := rel()
			if err != nil {
				return nil, err
			}
			node = NewBinNode(token.Pos(), NodeLt, lhs, node)
		} else {
			break
		}
	}

	return node, nil
}

func add() (*Node, error) {
	dprintf("add start\n")
	defer dprintf("add end\n")

	node, err := mul()
	if err != nil {
		return nil, err
	}
	dprintf("add lhs: %v\n", node)
	for len(tokens) > 0 {
		if token, ok := consumeSign("+"); ok {
			rhs, err := mul()
			if err != nil {
				return nil, err
			}
			node = NewBinNode(token.Pos(), NodeAdd, node, rhs)
			dprintf("add rhs: %v\n", node)
		} else if token, ok := consumeSign("-"); ok {
			rhs, err := mul()
			if err != nil {
				return nil, err
			}
			node = NewBinNode(token.Pos(), NodeSub, node, rhs)
			dprintf("sub rhs: %v\n", node)
		} else {
			break
		}
	}
	return node, nil
}

func mul() (*Node, error) {
	dprintf("mul start\n")
	defer dprintf("mul end\n")

	node, err := pri()
	if err != nil {
		return nil, err
	}
	dprintf("mul lhs: %v\n", node)
	for len(tokens) > 0 {
		if token, ok := consumeSign("*"); ok {
			rhs, err := pri()
			if err != nil {
				return nil, err
			}
			node = NewBinNode(token.Pos(), NodeMul, node, rhs)
			dprintf("mul rhs: %v\n", node)
		} else {
			break
		}
	}
	return node, nil
}

func pri() (*Node, error) {
	dprintf("pri start\n")
	defer dprintf("pri end\n")

	if _, ok := consumeSign("("); ok {
		node, err := expr()
		if err != nil {
			return nil, err
		}
		_, err = expectSign(")")
		if err != nil {
			return nil, err
		}
		return node, err
	}

	token, ident, ok := consumeIdent()
	if ok {
		if token, ok = consumeSign("("); ok {
			paramss, err := params()
			if err != nil {
				return nil, err
			}
			return NewCallNode(token.Pos(), ident, paramss), nil
		}

		if _, ok := funcTable[funcNameStack.peek()].Vars.Defined(ident); !ok {
			return nil, NewCompileError(token.Pos(), fmt.Sprintf("undefined variable found: %s", ident))
		}
		return NewVarNode(token.Pos(), ident), nil
	}

	node, err := val()
	if err != nil {
		return nil, err
	}
	return node, nil
}

func val() (*Node, error) {
	dprintf("val start\n")
	defer dprintf("val end\n")

	if len(tokens) == 0 {
		return nil, NewCompileError(nil, fmt.Sprintf("expect val but no more tokens"))
	}
	token := tokens[0]
	if token.Type != TokenVal {
		return nil, NewCompileError(token.Pos(), fmt.Sprintf("expect val but no got %v", token))
	}
	tokens = tokens[1:]
	return NewValNode(token.Pos(), token.Val), nil
}

func varStmt(varToken *Token) (*Node, error) {
	dprintf("var start\n")
	defer dprintf("var end\n")

	_, ident, err := expectIdent()
	if err != nil {
		return nil, err
	}
	valTypeToken, err := expectValType()
	if err != nil {
		return nil, err
	}
	var rhs *Node
	if _, ok := consumeSign("="); ok {
		rhs, err = expr()
		if err != nil {
			return nil, err
		}
	}

	node := NewDefVarNode(varToken.Pos(), ident, valTypeToken.ValType, rhs)

	err = expect(TokenEOL)
	if err != nil {
		return nil, err
	}
	funcTable[funcNameStack.peek()].Vars.Define(ident, valTypeToken.ValType)
	return node, nil
}

func ifStmt(ifToken *Token) (*Node, error) {
	dprintf("if start\n")
	defer dprintf("if end\n")

	condNode, err := expr()
	if err != nil {
		return nil, err
	}
	thenNode, err := block()
	if err != nil {
		return nil, err
	}
	var elseNode *Node
	if _, ok := consume(TokenElse); ok {
		elseNode, err = block()
		if err != nil {
			return nil, err
		}
	}
	return NewIfNode(ifToken.Pos(), condNode, thenNode, elseNode), nil
}

func whileStmt(whileToken *Token) (*Node, error) {
	dprintf("while start\n")
	defer dprintf("while end\n")

	condNode, err := expr()
	if err != nil {
		return nil, err
	}
	bodyNode, err := block()
	if err != nil {
		return nil, err
	}
	return NewWhileNode(whileToken.Pos(), condNode, bodyNode), nil
}
