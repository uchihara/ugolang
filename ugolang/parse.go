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

func prog() []Node {
	nodes := make([]Node, 0)
	for len(tokens) > 0 {
		node := funcStmt()
		nodes = append(nodes, *node)
	}
	return nodes
}

func funcStmt() *Node {
	dprintf("func start\n")
	expect(TokenFunc)
	ident := expectIdent()
	args := args()
	node := NewFuncNode(ident, args, block())
	return node
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
	return params
}

func block() *Node {
	dprintf("block start\n")
	node := NewBlockNode()
	expectSign("{")
	for len(tokens) > 0 {
		node.Statements = append(node.Statements, stmt())
		if consumeSign("}") {
			return node
		}
	}
	expectSign("}")
	return node
}

func stmt() *Node {
	dprintf("stmt start\n")
	if consume(TokenReturn) {
		node := NewReturnNode(expr())
		expect(TokenEOL)
		return node
	}
	if consume(TokenIf) {
		return ifStmt()
	}
	if consume(TokenWhile) {
		return whileStmt()
	}
	node := expr()
	expect(TokenEOL)
	return node
}

func expr() *Node {
	dprintf("expr start\n")
	return assign()
}

func assign() *Node {
	dprintf("assign start\n")
	node := eq()
	if consumeSign("=") {
		node = NewNode(NodeAssign, node, assign())
	}

	return node
}

func eq() *Node {
	dprintf("eq start\n")
	node := rel()
	for len(tokens) > 0 {
		if consumeSign("==") {
			node = NewNode(NodeEq, node, rel())
		} else if consumeSign("!=") {
			node = NewNode(NodeNe, node, rel())
		} else {
			break
		}
	}
	return node
}

func rel() *Node {
	dprintf("rel start\n")
	node := add()
	for len(tokens) > 0 {
		if consumeSign("<=") {
			node = NewNode(NodeLe, node, rel())
		} else if consumeSign("<") {
			node = NewNode(NodeLt, node, rel())
		} else if consumeSign(">=") {
			node = NewNode(NodeLe, rel(), node)
		} else if consumeSign(">") {
			node = NewNode(NodeLt, rel(), node)
		} else {
			break
		}
	}
	return node
}

func add() *Node {
	dprintf("add start\n")
	node := mul()
	dprintf("add lhs: %v\n", node)
	for len(tokens) > 0 {
		if consumeSign("+") {
			node = NewNode(NodeAdd, node, mul())
			dprintf("add rhs: %v\n", node)
		} else if consumeSign("-") {
			node = NewNode(NodeSub, node, mul())
			dprintf("sub rhs: %v\n", node)
		} else {
			break
		}
	}
	return node
}

func mul() *Node {
	dprintf("mul start\n")
	node := pri()
	dprintf("mul lhs: %v\n", node)
	for len(tokens) > 0 {
		if consumeSign("*") {
			node = NewNode(NodeMul, node, pri())
			dprintf("mul rhs: %v\n", node)
		} else {
			break
		}
	}
	return node
}

func pri() *Node {
	dprintf("pri start\n")
	if consume(TokenCall) {
		ident := expectIdent()
		prms := params()
		return NewCallNode(ident, prms)
	}

	if consumeSign("(") {
		node := expr()
		expectSign(")")
		return node
	}

	ident, ok := consumeIdent()
	if ok {
		return NewVarNode(ident)
	}

	return num()
}

func num() *Node {
	dprintf("num start\n")
	token := tokens[0]
	if token.Type != TokenNum {
		panic(fmt.Sprintf("%v expect num but got %v", caller(), token))
	}
	tokens = tokens[1:]
	return NewNumNode(token.Num)
}

func ifStmt() *Node {
	dprintf("if start\n")
	condNode := expr()
	thenNode := block()
	var elseNode *Node
	if consume(TokenElse) {
		elseNode = block()
	}
	return NewIfNode(condNode, thenNode, elseNode)
}

func whileStmt() *Node {
	dprintf("while start\n")
	condNode := expr()
	bodyNode := block()
	return NewWhileNode(condNode, bodyNode)
}
