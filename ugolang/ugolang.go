package ugolang

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
)

type stackType []int

var stack stackType

func (s *stackType) push(n int) {
	*s = append(*s, n)
}

func (s *stackType) pop() int {
	n := (*s)[len(*s)-1]
	*s = (*s)[0 : len(*s)-1]
	return n
}

type varsType map[rune]int

var vars varsType = map[rune]int{}

func (v varsType) Get(name rune) int {
	return v[name]
}

func (v varsType) Set(name rune, val int) {
	v[name] = val
}

type NodeType int

const (
	NodeNum NodeType = iota + 1
	NodeAdd
	NodeMul
	NodeAssign
	NodeVar
	NodeIf
	NodeElse
)

func (n NodeType) String() string {
	switch n {
	case NodeNum:
		return "num"
	case NodeAdd:
		return "add"
	case NodeMul:
		return "mul"
	case NodeAssign:
		return "assign"
	case NodeVar:
		return "var"
	case NodeIf:
		return "if"
	case NodeElse:
		return "else"
	default:
		return "unknown"
	}
}

type Node struct {
	Type  NodeType
	Val   int
	Ident rune
	Lhs   *Node
	Rhs   *Node
	Cond  *Node
	Then  *Node
	Else  *Node
}

func (n Node) String() string {
	switch n.Type {
	case NodeNum:
		return fmt.Sprintf("num(%d)", n.Val)
	case NodeAdd:
		return fmt.Sprintf("add(%s, %s)", n.Lhs.String(), n.Rhs.String())
	case NodeMul:
		return fmt.Sprintf("mul(%s, %s)", n.Lhs.String(), n.Rhs.String())
	case NodeAssign:
		return fmt.Sprintf("assign(%s, %s)", n.Lhs.String(), n.Rhs.String())
	case NodeVar:
		return fmt.Sprintf("var(%c)", n.Ident)
	case NodeIf:
		return fmt.Sprintf("if(%v, %v)", n.Cond, n.Then)
	case NodeElse:
		return fmt.Sprintf("else(%v)", n.Else)
	default:
		return fmt.Sprintf("unknown type: %d", n.Type)
	}
}

func NewNode(typ NodeType, lhs, rhs *Node) *Node {
	return &Node{
		Type: typ,
		Lhs:  lhs,
		Rhs:  rhs,
	}
}

func NewNumNode(val int) *Node {
	return &Node{
		Type: NodeNum,
		Val:  val,
	}
}

func NewVarNode(name rune) *Node {
	return &Node{
		Type:  NodeVar,
		Ident: name,
	}
}

func NewIfNode(condNode, thenNode, elseNode *Node) *Node {
	return &Node{
		Type: NodeIf,
		Cond: condNode,
		Then: thenNode,
		Else: elseNode,
	}
}

type TokenType int

const (
	TokenNum TokenType = iota + 1
	TokenSign
	TokenIdent
	TokenIf
	TokenElse
	TokenEOL
)

func (t TokenType) String() string {
	switch t {
	case TokenNum:
		return "numToken"
	case TokenSign:
		return "signToken"
	case TokenIdent:
		return "identToken"
	case TokenIf:
		return "ifToken"
	case TokenEOL:
		return "eolToken"
	default:
		return fmt.Sprintf("unknown type: %d", t)
	}
}

type Token struct {
	Type  TokenType
	Num   int
	Sign  rune
	Ident rune
}

func (t Token) String() string {
	switch t.Type {
	case TokenNum:
		return fmt.Sprintf("num(%d)", t.Num)
	case TokenSign:
		return fmt.Sprintf("sign(%c)", t.Sign)
	case TokenIdent:
		return fmt.Sprintf("ident(%c)", t.Ident)
	case TokenIf:
		return "if"
	case TokenElse:
		return "else"
	case TokenEOL:
		return ";"
	default:
		return fmt.Sprintf("unknown type: %v", t.Type)
	}
}

func NewToken(typ TokenType) *Token {
	return &Token{
		Type: typ,
	}
}

func NewNumToken(num int) *Token {
	return &Token{
		Type: TokenNum,
		Num:  num,
	}
}

func NewSignToken(sign rune) *Token {
	return &Token{
		Type: TokenSign,
		Sign: sign,
	}
}

func NewIdentToken(ident rune) *Token {
	return &Token{
		Type:  TokenIdent,
		Ident: ident,
	}
}

var tokens []Token

func Exec(code string) int {
	tokens = tokenize(code)
	nodes := prog()
	ret := 0
	for _, node := range nodes {
		eval(&node)
		dprintf("node=%v\n", node)
		ret = eval(&node)
	}
	return ret
}

func matchKeyword(kw, code string, idx int) (int, bool) {
		matchLen := len(kw) + 1
		if idx + matchLen <= len(code) && code[idx:idx+matchLen] == kw + " " {
			return len(kw), true
		}
		return 0, false
}

func tokenize(code string) []Token {
	tokens := make([]Token, 0)
	for i := 0; i < len(code); i++ {
		if matchLen, matched := matchKeyword("if", code, i); matched {
			tokens = append(tokens, *NewToken(TokenIf))
			i += matchLen
			continue
		}

		if matchLen, matched := matchKeyword("else", code, i); matched {
			tokens = append(tokens, *NewToken(TokenElse))
			i += matchLen
			continue
		}

		c := code[i]

		if c == ' ' {
			continue
		}

		if '0' <= c && c <= '9' {
			var j int = i + 1
			for ; '0' <= code[j] && code[j] <= '9' && j < len(code); j++ {
			}
			numStr := code[i:j]
			num, err := strconv.ParseInt(numStr, 10, 64)
			if err != nil {
				panic(fmt.Sprintf("invalid num format: %s", numStr))
			}
			tokens = append(tokens, *NewNumToken(int(num)))
			continue
		}

		if 'a' <= c && c <= 'z' {
			tokens = append(tokens, *NewIdentToken(rune(c)))
			continue
		}

		if c == '=' || c == '+' || c == '*' || c == '(' || c == ')' || c == '{' || c == '}' {
			tokens = append(tokens, *NewSignToken(rune(c)))
			continue
		}

		if c == ';' {
			tokens = append(tokens, *NewToken(TokenEOL))
			continue
		}
	}
	return tokens
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
