package ugolang

import (
	"fmt"
)

// NodeType dummy
type NodeType int

const (
	// NodeNum dummy
	NodeNum NodeType = iota + 1
	// NodeAdd dummy
	NodeAdd
	// NodeMul dummy
	NodeMul
	// NodeEq dummy
	NodeEq
	// NodeNe dummy
	NodeNe
	// NodeLt dummy
	NodeLt
	// NodeLe dummy
	NodeLe
	// NodeAssign dummy
	NodeAssign
	// NodeVar dummy
	NodeVar
	// NodeIf dummy
	NodeIf
	// NodeElse dummy
	NodeElse
	// NodeWhile dummy
	NodeWhile
)

func (n NodeType) String() string {
	switch n {
	case NodeNum:
		return "num"
	case NodeAdd:
		return "add"
	case NodeMul:
		return "mul"
	case NodeEq:
		return "eq"
	case NodeNe:
		return "ne"
	case NodeLe:
		return "le"
	case NodeLt:
		return "lt"
	case NodeAssign:
		return "assign"
	case NodeVar:
		return "var"
	case NodeIf:
		return "if"
	case NodeElse:
		return "else"
	case NodeWhile:
		return "while"
	default:
		return "unknown"
	}
}

// Node dummy
type Node struct {
	Type  NodeType
	Val   int
	Ident rune
	LHS   *Node
	RHS   *Node
	Cond  *Node
	Then  *Node
	Else  *Node
	Body  *Node
}

func (n Node) String() string {
	switch n.Type {
	case NodeNum:
		return fmt.Sprintf("num(%d)", n.Val)
	case NodeAdd:
		return fmt.Sprintf("add(%s, %s)", n.LHS.String(), n.RHS.String())
	case NodeMul:
		return fmt.Sprintf("mul(%s, %s)", n.LHS.String(), n.RHS.String())
	case NodeEq:
		return fmt.Sprintf("eq(%s, %s)", n.LHS.String(), n.RHS.String())
	case NodeNe:
		return fmt.Sprintf("ne(%s, %s)", n.LHS.String(), n.RHS.String())
	case NodeLe:
		return fmt.Sprintf("le(%s, %s)", n.LHS.String(), n.RHS.String())
	case NodeLt:
		return fmt.Sprintf("lt(%s, %s)", n.LHS.String(), n.RHS.String())
	case NodeAssign:
		return fmt.Sprintf("assign(%s, %s)", n.LHS.String(), n.RHS.String())
	case NodeVar:
		return fmt.Sprintf("var(%c)", n.Ident)
	case NodeIf:
		return fmt.Sprintf("if(%v, %v, %v)", n.Cond, n.Then, n.Else)
	case NodeElse:
		return fmt.Sprintf("else(%v)", n.Else)
	case NodeWhile:
		return fmt.Sprintf("while(%v, %v)", n.Cond, n.Body)
	default:
		return fmt.Sprintf("unknown type: %d", n.Type)
	}
}

// NewNode dummy
func NewNode(typ NodeType, lhs, rhs *Node) *Node {
	return &Node{
		Type: typ,
		LHS:  lhs,
		RHS:  rhs,
	}
}

// NewNumNode dummy
func NewNumNode(val int) *Node {
	return &Node{
		Type: NodeNum,
		Val:  val,
	}
}

// NewVarNode dummy
func NewVarNode(name rune) *Node {
	return &Node{
		Type:  NodeVar,
		Ident: name,
	}
}

// NewIfNode dummy
func NewIfNode(condNode, thenNode, elseNode *Node) *Node {
	return &Node{
		Type: NodeIf,
		Cond: condNode,
		Then: thenNode,
		Else: elseNode,
	}
}

// NewWhileNode dummy
func NewWhileNode(condNode, bodyNode *Node) *Node {
	return &Node{
		Type: NodeWhile,
		Cond: condNode,
		Body: bodyNode,
	}
}
