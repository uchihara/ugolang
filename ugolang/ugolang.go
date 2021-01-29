package ugolang

import (
	"fmt"
	"path/filepath"
	"runtime"
)

type Ugolang struct {
	DumpTokens bool
	DumpNodes  bool
}

func NewUgolang() *Ugolang {
	return &Ugolang{
	}
}

var tokens []Token

func (u *Ugolang) Exec(code string) int {
	tokens = tokenize(code)
	if u.DumpTokens {
		fmt.Printf("tokens: %v\n", tokens)
	}
	nodes := prog()
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
		l := eval(node.Lhs)
		r := eval(node.Rhs)
		return l + r
	case NodeMul:
		l := eval(node.Lhs)
		r := eval(node.Rhs)
		return l * r
	case NodeAssign:
		val := eval(node.Rhs)
		vars.Set(node.Lhs.Ident, val)
		return val
	case NodeVar:
		return vars.Get(node.Ident)
	case NodeIf:
		cond := eval(node.Cond)
		if cond != 0 {
			return eval(node.Then)
		} else {
			if node.Else != nil {
				return eval(node.Else)
			}
		}
		return 0 // FIXME
	default:
		panic(fmt.Sprintf("unknown type: %d", node.Type))
	}
}

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

func consumeSign(sign rune) bool {
	if len(tokens) == 0 {
		return false
	}
	if tokens[0].Type == TokenSign && tokens[0].Sign == sign {
		tokens = tokens[1:]
		return true
	}
	return false
}

func consumeIdent() (rune, bool) {
	if len(tokens) == 0 {
		return '?', false
	}
	token := tokens[0]
	if token.Type == TokenIdent {
		tokens = tokens[1:]
		return token.Ident, true
	}
	return '?', false
}

func expect(tokenType TokenType) {
	if !consume(tokenType) {
		panic(fmt.Sprintf("expect %v but got %v", tokenType, tokens[0].Type))
	}
}

func expectSign(sign rune) {
	if !consumeSign(sign) {
		panic(fmt.Sprintf("expect %c but got %v", sign, tokens[0]))
	}
}

func dprintf(f string, param ...interface{}) {
	if true {
		return
	}
	depth := 0
	for i := 1; ; i++ {
		_, _, _, ok := runtime.Caller(i)
		if !ok {
			break
		}
		depth++
	}
	pc, _, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	fileName, fileLine := fn.FileLine(pc)
	fmt.Printf("%s:%d: ", filepath.Base(fileName), fileLine)

	for i := 8; i < depth; i++ {
		fmt.Printf(" ")
	}

	fmt.Printf(f, param...)
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
	node := expr()
	expect(TokenEOL)
	return node
}

func expr() *Node {
	return assign()
}

func assign() *Node {
	node := add()
	if consumeSign('=') {
		node = NewNode(NodeAssign, node, assign())
	}

	return node
}

func add() *Node {
	dprintf("add start\n")
	node := mul()
	dprintf("add lhs: %v\n", node)
	for len(tokens) > 0 {
		if consumeSign('+') {
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
		if consumeSign('*') {
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
	if consumeSign('(') {
		node := expr()
		expectSign(')')
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
	expectSign('{')
	thenNode := stmt()
	expectSign('}')
	var elseNode *Node
	if consume(TokenElse) {
		expectSign('{')
		elseNode = stmt()
		expectSign('}')
	}
	return NewIfNode(condNode, thenNode, elseNode)
}
