package ugolang

import (
	"fmt"
)

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
