package ugolang

import (
	"fmt"
	"regexp"
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
	return len(token), true
}

func matchSign(sign, code string, idx int) (int, bool) {
	signLen := len(sign)
	if idx+signLen > len(code) {
		return 0, false
	}
	if code[idx:idx+signLen] != sign {
		return 0, false
	}
	return len(sign), true
}

func matchSigns(signs []string, code string, idx int) (int, bool) {
	for _, sign := range signs {
		if matchLen, matched := matchSign(sign, code, idx); matched {
			return matchLen, matched
		}
	}
	return 0, false
}

func matchPattern(pattern, code string, idx int) (int, bool) {
	re := regexp.MustCompile(pattern)
	loc := re.FindStringIndex(code[idx:len(code)])
	if len(loc) != 2 {
		return 0, false
	}
	return loc[1] - loc[0], true
}

func tokenize(code string) []Token {
	tokens := make([]Token, 0)
	for i := 0; i < len(code); i++ {
		if matchLen, matched := matchToken("if", code, i); matched {
			tokens = append(tokens, *NewToken(TokenIf))
			i += (matchLen - 1)
			continue
		}

		if matchLen, matched := matchToken("else", code, i); matched {
			tokens = append(tokens, *NewToken(TokenElse))
			i += (matchLen - 1)
			continue
		}

		if matchLen, matched := matchToken("while", code, i); matched {
			tokens = append(tokens, *NewToken(TokenWhile))
			i += (matchLen - 1)
			continue
		}

		if matchLen, matched := matchToken("func", code, i); matched {
			tokens = append(tokens, *NewToken(TokenFunc))
			i += (matchLen - 1)
			continue
		}

		if matchLen, matched := matchToken("call", code, i); matched {
			tokens = append(tokens, *NewToken(TokenCall))
			i += (matchLen - 1)
			continue
		}

		c := code[i]

		if c == ' ' {
			continue
		}

		if matchLen, matched := matchPattern("^[0-9]+", code, i); matched {
			numStr := code[i : i+matchLen]
			num, err := strconv.ParseInt(numStr, 10, 64)
			if err != nil {
				panic(fmt.Sprintf("invalid num format: %s", numStr))
			}
			tokens = append(tokens, *NewNumToken(int(num)))
			i += (matchLen - 1)
			continue
		}

		if matchLen, matched := matchPattern("^[A-Za-z0-9_]+", code, i); matched {
			tokens = append(tokens, *NewIdentToken(code[i : i+matchLen]))
			i += (matchLen - 1)
			continue
		}

		signs := []string{"==", "!=", "<=", ">=", "<", ">", "=", "+", "*", "(", ")", "{", "}"}
		if matchLen, matched := matchSigns(signs, code, i); matched {
			tokens = append(tokens, *NewSignToken(code[i : i+matchLen]))
			i += (matchLen - 1)
			continue
		}

		if c == ';' {
			tokens = append(tokens, *NewToken(TokenEOL))
			continue
		}
	}
	return tokens
}
