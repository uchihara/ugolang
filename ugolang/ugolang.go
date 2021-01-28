package ugolang

import (
  "path/filepath"
  "fmt"
  "runtime"
)

type stackType []int

var stack stackType

func (s *stackType) push(n int) {
  *s = append(*s, n)
}

func (s *stackType) pop() int {
  n := (*s)[len(*s)-1]
  *s = (*s)[0:len(*s)-1]
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
  default:
    return "unknown"
  }
}

type Node struct {
  Type NodeType
  Val   int
  Ident rune
  Lhs   *Node
  Rhs   *Node
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

var codes []rune

func Exec(code string) int {
  codes = make([]rune, 0)
  for _, c := range code {
    codes = append(codes, c)
  }
  nodes := prog()
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
    return l+r
  case NodeMul:
    l := eval(node.Lhs)
    r := eval(node.Rhs)
    return l*r
  case NodeAssign:
    val := eval(node.Rhs)
    vars.Set(node.Lhs.Ident, val)
    return val
  case NodeVar:
    return vars.Get(node.Ident)
  default:
    panic(fmt.Sprintf("unknown type: %d", node.Type))
  }
}

func consume(c rune) bool {
  if c == codes[0] {
    codes = codes[1:]
    return true
  }
  return false
}

func consumeIdent() (rune, bool) {
  c := codes[0]
  if 'a' <= c && c <= 'z' {
    codes = codes[1:]
    return c, true
  }
  return '?', false
}

func expect(c rune) {
  if !consume(c) {
    panic(fmt.Sprintf("expect %c but got %c", c, codes[0]))
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
  for ; len(codes) > 0; {
    node := stmt()
    nodes = append(nodes, *node)
  }
  return nodes
}

func stmt() *Node {
  node := expr()
  expect(';')
  return node
}

func expr() *Node {
  return assign()
}

func assign() *Node {
  node := add()
  if consume('=') {
    node = NewNode(NodeAssign, node, assign())
  }

  return node
}

func add() *Node {
  dprintf("add start\n")
  node := mul()
  dprintf("add lhs: %v\n", node)
  for ; len(codes) > 0; {
    if consume('+') {
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
  for ; len(codes) > 0; {
    if consume('*') {
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
  if consume('(') {
    node := expr()
    expect(')')
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
  c := codes[0]
  if c < '0' || '9' < c {
    panic(fmt.Sprintf("expect num but got %c", c))
  }
  codes = codes[1:]
  return NewNumNode(int(c - '0'))
}
