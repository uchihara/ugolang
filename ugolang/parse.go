package ugolang

import (
	"fmt"
)

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

func expectIdent() (*Token, string, error) {
	token, ident, ok := consumeIdent()
	if !ok {
		return nil, "", NewCompileError(token.Pos(), fmt.Sprintf("%v expect ident but got %v", caller(), tokens[0]))
	}
	return token, ident, nil
}

func prog() ([]*Node, error) {
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
	nodes = append(nodes, NewCallNode(nil, "main", []*Node{}))
	funcStack.push("main")
	return nodes, nil
}

func funcStmt() (node *Node, ok bool, err error) {
	dprintf("func start\n")
	var ident string
	var argss []string
	var blockNode *Node
	token, isFunc := consume(TokenFunc)
	if !isFunc {
		goto end
	}
	_, ident, err = expectIdent()
	if err != nil {
		goto end
	}
	argss, err = args()
	if err != nil {
		goto end
	}
	blockNode, err = block()
	if err != nil {
		goto end
	}
	node = NewFuncNode(token.Pos(), ident, argss, blockNode)
	ok = true
end:
	dprintf("func end\n")
	return node, ok, err
}

func args() (args []string, err error) {
	dprintf("args start\n")
	_, err = expectSign("(")
	if err != nil {
		goto end
	}
	args = make([]string, 0)
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
		_, ident, err = expectIdent()
		if err != nil {
			goto end
		}
		args = append(args, ident)
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
