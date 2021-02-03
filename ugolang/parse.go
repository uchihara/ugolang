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
		panic(fmt.Sprintf("expect %v but got %v", tokenType, tokens[0].Type))
	}
}

func expectSign(sign string) {
	if !consumeSign(sign) {
		panic(fmt.Sprintf("expect %s but got %v", sign, tokens[0]))
	}
}

func prog() []Node {
	nodes := make([]Node, 0)
	for len(tokens) > 0 {
		node := stmt()
		nodes = append(nodes, *node)
	}
	return nodes
}

func stmt() *Node {
	if consume(TokenIf) {
		return if_()
	}
	if consume(TokenWhile) {
		return while_()
	}
	node := expr()
	expect(TokenEOL)
	return node
}

func expr() *Node {
	return assign()
}

func assign() *Node {
	node := eq()
	if consumeSign("=") {
		node = NewNode(NodeAssign, node, assign())
	}

	return node
}

func eq() *Node {
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
		panic(fmt.Sprintf("expect num but got %v", token))
	}
	tokens = tokens[1:]
	return NewNumNode(token.Num)
}

func if_() *Node {
	dprintf("if start\n")
	condNode := expr()
	expectSign("{")
	thenNode := stmt()
	expectSign("}")
	var elseNode *Node
	if consume(TokenElse) {
		expectSign("{")
		elseNode = stmt()
		expectSign("}")
	}
	return NewIfNode(condNode, thenNode, elseNode)
}

func while_() *Node {
	dprintf("while start\n")
	condNode := expr()
	expectSign("{")
	bodyNode := stmt()
	expectSign("}")
	return NewWhileNode(condNode, bodyNode)
}
