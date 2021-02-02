package ugolang

import (
	"fmt"
)

// TokenType dummmy
type TokenType int

const (
	// TokenNum dummmy
	TokenNum TokenType = iota + 1
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
	// TokenEOL dummmy
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
	case TokenWhile:
		return "whileToken"
	case TokenEOL:
		return "eolToken"
	default:
		return fmt.Sprintf("unknown type: %d", t)
	}
}

// Token dummy
type Token struct {
	Type  TokenType
	Num   int
	Sign  string
	Ident rune
}

func (t Token) String() string {
	switch t.Type {
	case TokenNum:
		return fmt.Sprintf("num(%d)", t.Num)
	case TokenSign:
		return fmt.Sprintf("sign(%s)", t.Sign)
	case TokenIdent:
		return fmt.Sprintf("ident(%c)", t.Ident)
	case TokenIf:
		return "if"
	case TokenElse:
		return "else"
	case TokenWhile:
		return "while"
	case TokenEOL:
		return ";"
	default:
		return fmt.Sprintf("unknown type: %v", t.Type)
	}
}

// NewToken dummy
func NewToken(typ TokenType) *Token {
	return &Token{
		Type: typ,
	}
}

// NewNumToken dummy
func NewNumToken(num int) *Token {
	return &Token{
		Type: TokenNum,
		Num:  num,
	}
}

// NewSignToken dummy
func NewSignToken(sign string) *Token {
	return &Token{
		Type: TokenSign,
		Sign: sign,
	}
}

// NewIdentToken dummy
func NewIdentToken(ident rune) *Token {
	return &Token{
		Type:  TokenIdent,
		Ident: ident,
	}
}
