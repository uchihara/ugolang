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
	// NodeFunc dummy
	NodeFunc
	// NodeCall dummy
	NodeCall
	// NodeReturn dummy
	NodeReturn
	// NodeBlock dummy
	NodeBlock
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
	case NodeFunc:
		return "func"
	case NodeCall:
		return "call"
	case NodeReturn:
		return "return"
	case NodeBlock:
		return "block"
	default:
		return "unknown"
	}
}

// Node dummy
type Node struct {
	Type       NodeType
	Val        int
	Ident      string
	LHS        *Node
	RHS        *Node
	Cond       *Node
	Then       *Node
	Else       *Node
	Body       *Node
	Expr       *Node
	Statements []*Node
	Args       []string
	Params     []*Node
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
		return fmt.Sprintf("var(%s)", n.Ident)
	case NodeIf:
		return fmt.Sprintf("if(%v, %v, %v)", n.Cond, n.Then, n.Else)
	case NodeElse:
		return fmt.Sprintf("else(%v)", n.Else)
	case NodeWhile:
		return fmt.Sprintf("while(%v, %v)", n.Cond, n.Body)
	case NodeFunc:
		return fmt.Sprintf("func(%s, %v, %v)", n.Ident, n.Args, n.Body)
	case NodeCall:
		return fmt.Sprintf("call(%s, %v)", n.Ident, n.Params)
	case NodeReturn:
		return fmt.Sprintf("return(%v)", n.Expr)
	case NodeBlock:
		s := ""
		for _, stmt := range n.Statements {
			if len(s) != 0 {
				s += ", "
			}
			s += stmt.String()
		}
		return fmt.Sprintf("block(%s)", s)
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
func NewVarNode(name string) *Node {
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

// NewFuncNode dummy
func NewFuncNode(name string, args []string, bodyNode *Node) *Node {
	return &Node{
		Type:  NodeFunc,
		Ident: name,
		Args:  args,
		Body:  bodyNode,
	}
}

// NewCallNode dummy
func NewCallNode(name string, params []*Node) *Node {
	return &Node{
		Type:   NodeCall,
		Ident:  name,
		Params: params,
	}
}

// NewReturnNode dummy
func NewReturnNode(expr *Node) *Node {
	return &Node{
		Type: NodeReturn,
		Expr: expr,
	}
}

// NewBlockNode dummy
func NewBlockNode() *Node {
	return &Node{
		Type: NodeBlock,
	}
}
