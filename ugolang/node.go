package ugolang

import (
	"fmt"
)

// NodeType dummy
type NodeType int

const (
	// NodeVal dummy
	NodeVal NodeType = iota + 1
	// NodeDefVar dummy
	NodeDefVar
	// NodeAdd dummy
	NodeAdd
	// NodeSub dummy
	NodeSub
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
	// NodeBreak dummy
	NodeBreak
	// NodeContinue dummy
	NodeContinue
	// NodeBlock dummy
	NodeBlock
)

func (n NodeType) String() string {
	switch n {
	case NodeVal:
		return "val"
	case NodeDefVar:
		return "defvar"
	case NodeAdd:
		return "add"
	case NodeSub:
		return "sub"
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
	case NodeBreak:
		return "break"
	case NodeContinue:
		return "continue"
	case NodeBlock:
		return "block"
	default:
		return "unknown"
	}
}

// Node dummy
type Node struct {
	TokenPos   *TokenPos
	Type       NodeType
	Val        *Val
	ValType    ValType
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
	case NodeVal:
		return fmt.Sprintf("val(%v)", n.Val)
	case NodeDefVar:
		return fmt.Sprintf("defvar(%s, %v, %v)", n.Ident, n.ValType, n.RHS)
	case NodeAdd:
		return fmt.Sprintf("add(%s, %s)", n.LHS.String(), n.RHS.String())
	case NodeSub:
		return fmt.Sprintf("sub(%s, %s)", n.LHS.String(), n.RHS.String())
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
	case NodeBreak:
		return fmt.Sprintf("break")
	case NodeContinue:
		return fmt.Sprintf("continue")
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
func NewNode(tokenPos *TokenPos, typ NodeType) *Node {
	return &Node{
		TokenPos: tokenPos,
		Type:     typ,
	}
}

// NewBinNode dummy
func NewBinNode(tokenPos *TokenPos, typ NodeType, lhs, rhs *Node) *Node {
	return &Node{
		TokenPos: tokenPos,
		Type:     typ,
		LHS:      lhs,
		RHS:      rhs,
	}
}

// NewValNode dummy
func NewValNode(tokenPos *TokenPos, val *Val) *Node {
	return &Node{
		TokenPos: tokenPos,
		Type:     NodeVal,
		Val:      val,
	}
}

// NewDefVarNode dummy
func NewDefVarNode(tokenPos *TokenPos, ident string, valType ValType, rhs *Node) *Node {
	return &Node{
		TokenPos: tokenPos,
		Type:     NodeDefVar,
		Ident:    ident,
		ValType:  valType,
		RHS:      rhs,
	}
}

// NewVarNode dummy
func NewVarNode(tokenPos *TokenPos, name string) *Node {
	return &Node{
		TokenPos: tokenPos,
		Type:     NodeVar,
		Ident:    name,
	}
}

// NewIfNode dummy
func NewIfNode(tokenPos *TokenPos, condNode, thenNode, elseNode *Node) *Node {
	return &Node{
		TokenPos: tokenPos,
		Type:     NodeIf,
		Cond:     condNode,
		Then:     thenNode,
		Else:     elseNode,
	}
}

// NewWhileNode dummy
func NewWhileNode(tokenPos *TokenPos, condNode, bodyNode *Node) *Node {
	return &Node{
		TokenPos: tokenPos,
		Type:     NodeWhile,
		Cond:     condNode,
		Body:     bodyNode,
	}
}

// NewFuncNode dummy
func NewFuncNode(tokenPos *TokenPos, name string, args []string, bodyNode *Node) *Node {
	return &Node{
		TokenPos: tokenPos,
		Type:     NodeFunc,
		Ident:    name,
		Args:     args,
		Body:     bodyNode,
	}
}

// NewCallNode dummy
func NewCallNode(tokenPos *TokenPos, name string, params []*Node) *Node {
	return &Node{
		TokenPos: tokenPos,
		Type:     NodeCall,
		Ident:    name,
		Params:   params,
	}
}

// NewReturnNode dummy
func NewReturnNode(tokenPos *TokenPos, expr *Node) *Node {
	return &Node{
		TokenPos: tokenPos,
		Type:     NodeReturn,
		Expr:     expr,
	}
}

// NewBlockNode dummy
func NewBlockNode(tokenPos *TokenPos) *Node {
	return &Node{
		TokenPos: tokenPos,
		Type:     NodeBlock,
	}
}
