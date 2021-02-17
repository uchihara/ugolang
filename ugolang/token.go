package ugolang

import (
	"fmt"
)

// TokenType dummmy
type TokenType int

const (
	// TokenVal dummmy
	TokenVal TokenType = iota + 1
	// TokenVar dummy
	TokenVar
	// TokenValType dummy
	TokenValType
	// TokenSign dummmy
	TokenSign
	// TokenIdent dummmy
	TokenIdent
	// TokenIf dummmy
	TokenIf
	// TokenElse dummmy
	TokenElse
	// TokenWhile dummmy
	TokenWhile
	// TokenFunc dummy
	TokenFunc
	// TokenReturn dummy
	TokenReturn
	// TokenBreak dummy
	TokenBreak
	// TokenContinue dummy
	TokenContinue
	// TokenEOL dummmy
	TokenEOL
)

func (t TokenType) String() string {
	switch t {
	case TokenVal:
		return "valToken"
	case TokenVar:
		return "varToken"
	case TokenValType:
		return "valTypeToken"
	case TokenSign:
		return "signToken"
	case TokenIdent:
		return "identToken"
	case TokenIf:
		return "ifToken"
	case TokenWhile:
		return "whileToken"
	case TokenFunc:
		return "funcToken"
	case TokenReturn:
		return "returnToken"
	case TokenBreak:
		return "breakToken"
	case TokenContinue:
		return "continueToken"
	case TokenEOL:
		return "eolToken"
	default:
		return fmt.Sprintf("unknown type: %d", t)
	}
}

// TokenPos dummy
type TokenPos struct {
	line   int
	column int
}

// NewTokenPos dummy
func NewTokenPos(line, column int) *TokenPos {
	return &TokenPos{
		line:   line,
		column: column,
	}
}

// Line dummy
func (p *TokenPos) Line() int {
	if p == nil {
		return 0
	}
	return p.line
}

// Column dummy
func (p *TokenPos) Column() int {
	if p == nil {
		return 0
	}
	return p.column
}

// Token dummy
type Token struct {
	pos     *TokenPos
	Type    TokenType
	Val     *Val
	ValType ValType
	Sign    string
	Ident   string
}

// Pos dummy
func (t *Token) Pos() *TokenPos {
	if t == nil {
		return nil
	}
	return t.pos
}

func (t *Token) String() string {
	switch t.Type {
	case TokenVal:
		return fmt.Sprintf("val(%v)", t.Val)
	case TokenVar:
		return fmt.Sprintf("var")
	case TokenValType:
		return fmt.Sprintf("valType(%v)", t.ValType)
	case TokenSign:
		return fmt.Sprintf("sign(%s)", t.Sign)
	case TokenIdent:
		return fmt.Sprintf("ident(%s)", t.Ident)
	case TokenIf:
		return "if"
	case TokenElse:
		return "else"
	case TokenWhile:
		return "while"
	case TokenFunc:
		return "func"
	case TokenReturn:
		return "return"
	case TokenBreak:
		return "break"
	case TokenContinue:
		return "continue"
	case TokenEOL:
		return ";"
	default:
		return fmt.Sprintf("unknown type: %v", t.Type)
	}
}

// NewToken dummy
func NewToken(line, column int, typ TokenType) *Token {
	return &Token{
		pos:  NewTokenPos(line, column),
		Type: typ,
	}
}

// NewValToken dummy
func NewValToken(line, column int, val *Val) *Token {
	return &Token{
		pos:  NewTokenPos(line, column),
		Type: TokenVal,
		Val:  val,
	}
}

// NewVarToken dummy
func NewVarToken(line, column int, ident string) *Token {
	return &Token{
		pos:   NewTokenPos(line, column),
		Type:  TokenVar,
		Ident: ident,
	}
}

// NewValTypeToken dummy
func NewValTypeToken(line, column int, valType ValType) *Token {
	return &Token{
		pos:     NewTokenPos(line, column),
		Type:    TokenValType,
		ValType: valType,
	}
}

// NewSignToken dummy
func NewSignToken(line, column int, sign string) *Token {
	return &Token{
		pos:  NewTokenPos(line, column),
		Type: TokenSign,
		Sign: sign,
	}
}

// NewIdentToken dummy
func NewIdentToken(line, column int, ident string) *Token {
	return &Token{
		pos:   NewTokenPos(line, column),
		Type:  TokenIdent,
		Ident: ident,
	}
}
