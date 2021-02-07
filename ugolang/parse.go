package ugolang

import (
	"fmt"
)

func consume(tokenType TokenType) bool {
	if len(tokens) == 0 {
		return false
	}
	if tokenType == tokens[0].Type {
		tokens = tokens[1:]
		return true
	}
	return false
}

func consumeSign(sign string) bool {
	if len(tokens) == 0 {
		return false
	}
	if tokens[0].Type == TokenSign && tokens[0].Sign == sign {
		tokens = tokens[1:]
		return true
	}
	return false
}

func consumeIdent() (string, bool) {
	if len(tokens) == 0 {
		return "", false
	}
	token := tokens[0]
	if token.Type == TokenIdent {
		tokens = tokens[1:]
		return token.Ident, true
	}
	return "", false
}

func expect(tokenType TokenType) {
	if !consume(tokenType) {
		panic(fmt.Sprintf("%v expect %v but got %v", caller(), tokenType, tokens[0].Type))
	}
}

func expectSign(sign string) {
	if !consumeSign(sign) {
		panic(fmt.Sprintf("%v expect %s but got %v", caller(), sign, tokens[0]))
	}
}

func expectIdent() string {
	ident, ok := consumeIdent()
	if !ok {
		panic(fmt.Sprintf("%v expect ident but got %v", caller(), tokens[0]))
	}
	return ident
}

func prog() []*Node {
	nodes := make([]*Node, 0)
	for len(tokens) > 0 {
		node, ok := funcStmt()
		if !ok {
			node = stmt()
		}
		nodes = append(nodes, node)
	}
	nodes = append(nodes, NewCallNode("main", []*Node{}))
	funcStack.push("main")
	return nodes
}

func funcStmt() (node *Node, ok bool) {
	dprintf("func start\n")
	if !consume(TokenFunc) {
		node, ok = nil, false
		goto end
	}
	node = NewFuncNode(expectIdent(), args(), block())
	ok = true
end:
	dprintf("func end\n")
	return node, ok
}

func args() []string {
	dprintf("args start\n")
	expectSign("(")
	args := make([]string, 0)
	for len(tokens) > 0 {
		if consumeSign(")") {
			break
		}
		if len(args) != 0 {
			expectSign(",")
		}
		args = append(args, expectIdent())
	}
	dprintf("args end\n")
	return args
}

func params() []*Node {
	dprintf("params start\n")
	expectSign("(")
	params := make([]*Node, 0)
	for len(tokens) > 0 {
		if consumeSign(")") {
			break
		}
		if len(params) != 0 {
			expectSign(",")
		}
		params = append(params, expr())
	}
	dprintf("params end\n")
	return params
}

func block() *Node {
	dprintf("block start\n")
	node := NewBlockNode()
	expectSign("{")
	for len(tokens) > 0 {
		node.Statements = append(node.Statements, stmt())
		if consumeSign("}") {
			goto end
		}
	}
	expectSign("}")
end:
	dprintf("block end\n")
	return node
}

func stmt() (node *Node) {
	dprintf("stmt start\n")
	if consume(TokenReturn) {
		node = NewReturnNode(expr())
		expect(TokenEOL)
		goto end
	}
	if consume(TokenBreak) {
		node = NewNode(NodeBreak)
		expect(TokenEOL)
		goto end
	}
	if consume(TokenContinue) {
		node = NewNode(NodeContinue)
		expect(TokenEOL)
		goto end
	}
	if consume(TokenIf) {
		node = ifStmt()
		goto end
	}
	if consume(TokenWhile) {
		node = whileStmt()
		goto end
	}
	node = expr()
	expect(TokenEOL)
end:
	dprintf("stmt end\n")
	return node
}

func expr() *Node {
	dprintf("expr start\n")
	node := assign()
	dprintf("expr end\n")
	return node
}

func assign() *Node {
	dprintf("assign start\n")
	node := eq()
	if consumeSign("=") {
		node = NewBinNode(NodeAssign, node, assign())
	}
	dprintf("assign end\n")
	return node
}

func eq() *Node {
	dprintf("eq start\n")
	node := rel()
	for len(tokens) > 0 {
		if consumeSign("==") {
			node = NewBinNode(NodeEq, node, rel())
		} else if consumeSign("!=") {
			node = NewBinNode(NodeNe, node, rel())
		} else {
			break
		}
	}
	dprintf("eq end\n")
	return node
}

func rel() *Node {
	dprintf("rel start\n")
	node := add()
	for len(tokens) > 0 {
		if consumeSign("<=") {
			node = NewBinNode(NodeLe, node, rel())
		} else if consumeSign("<") {
			node = NewBinNode(NodeLt, node, rel())
		} else if consumeSign(">=") {
			node = NewBinNode(NodeLe, rel(), node)
		} else if consumeSign(">") {
			node = NewBinNode(NodeLt, rel(), node)
		} else {
			break
		}
	}
	dprintf("rel end\n")
	return node
}

func add() *Node {
	dprintf("add start\n")
	node := mul()
	dprintf("add lhs: %v\n", node)
	for len(tokens) > 0 {
		if consumeSign("+") {
			node = NewBinNode(NodeAdd, node, mul())
			dprintf("add rhs: %v\n", node)
		} else if consumeSign("-") {
			node = NewBinNode(NodeSub, node, mul())
			dprintf("sub rhs: %v\n", node)
		} else {
			break
		}
	}
	dprintf("add end\n")
	return node
}

func mul() *Node {
	dprintf("mul start\n")
	node := pri()
	dprintf("mul lhs: %v\n", node)
	for len(tokens) > 0 {
		if consumeSign("*") {
			node = NewBinNode(NodeMul, node, pri())
			dprintf("mul rhs: %v\n", node)
		} else {
			break
		}
	}
	dprintf("mul end\n")
	return node
}

func pri() (node *Node) {
	dprintf("pri start\n")
	var ident string
	var ok bool
	if consume(TokenCall) {
		node = NewCallNode(expectIdent(), params())
		goto end
	}

	if consumeSign("(") {
		node = expr()
		expectSign(")")
		goto end
	}

	ident, ok = consumeIdent()
	if ok {
		node = NewVarNode(ident)
		goto end
	}

	node = num()
end:
	dprintf("pri end\n")
	return node
}

func num() *Node {
	dprintf("num start\n")
	token := tokens[0]
	if token.Type != TokenNum {
		panic(fmt.Sprintf("%v expect num but got %v", caller(), token))
	}
	tokens = tokens[1:]
	node := NewNumNode(token.Num)
	dprintf("num end\n")
	return node
}

func ifStmt() *Node {
	dprintf("if start\n")
	condNode := expr()
	thenNode := block()
	var elseNode *Node
	if consume(TokenElse) {
		elseNode = block()
	}
	node := NewIfNode(condNode, thenNode, elseNode)
	dprintf("if end\n")
	return node
}

func whileStmt() *Node {
	dprintf("while start\n")
	condNode := expr()
	bodyNode := block()
	node := NewWhileNode(condNode, bodyNode)
	dprintf("while end\n")
	return node
}
