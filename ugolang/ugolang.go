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

type NodeType int

const (
  NodeNum NodeType = iota + 1
  NodeAdd
  NodeMul
)

func (n NodeType) String() string {
  switch n {
  case NodeNum:
    return "num"
  case NodeAdd:
    return "add"
  case NodeMul:
    return "mul"
  default:
    return "unknown"
  }
}

type Node struct {
  Type NodeType
  Val  int
  Lhs  *Node
  Rhs  *Node
}

func (n Node) String() string {
  switch n.Type {
  case NodeNum:
    return fmt.Sprintf("num(%d)", n.Val)
  case NodeAdd:
    return fmt.Sprintf("add(%s, %s)", n.Lhs.String(), n.Rhs.String())
  case NodeMul:
    return fmt.Sprintf("mul(%s, %s)", n.Lhs.String(), n.Rhs.String())
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

var codes []rune

func Exec(code string) int {
  codes = make([]rune, 0)
  for _, c := range code {
    codes = append(codes, c)
  }
  node := expr()
  dprintf("node=%v\n", node)
  n := eval(node)
  return n
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

func expr() *Node {
  return add()
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
