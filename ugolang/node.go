package ugolang

import (
	"fmt"
)

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
		return fmt.Sprintf("if(%v, %v, %v)", n.Cond, n.Then, n.Else)
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
