package ugolang

import (
	"fmt"
	"strconv"
)

func matchToken(token, code string, idx int) (int, bool) {
	tokenLen := len(token)
	if idx+tokenLen+1 > len(code) {
		return 0, false
	}
	if code[idx:idx+tokenLen] != token {
		return 0, false
	}
	nextChar := code[idx+tokenLen]
	if '0' <= nextChar && nextChar <= '9' ||
		'a' <= nextChar && nextChar <= 'z' ||
		'A' <= nextChar && nextChar <= 'Z' ||
		nextChar == '_' {
		return 0, false
	}
	return len(token) - 1, true
}

func tokenize(code string) []Token {
	tokens := make([]Token, 0)
	for i := 0; i < len(code); i++ {
		if matchLen, matched := matchToken("if", code, i); matched {
			tokens = append(tokens, *NewToken(TokenIf))
			i += matchLen
			continue
		}

		if matchLen, matched := matchToken("else", code, i); matched {
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
			tokens = append(tokens, *NewSignToken(string(c)))
			continue
		}

		if c == ';' {
			tokens = append(tokens, *NewToken(TokenEOL))
			continue
		}
	}
	return tokens
}
