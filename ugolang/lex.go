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

type tokenPair struct {
	keyword string
	token   TokenType
}

func matchTokens(tokenPairs []tokenPair, code string, idx int) (int, bool, TokenType) {
	for _, tokenPair := range tokenPairs {
		if matchLen, matched := matchToken(tokenPair.keyword, code, idx); matched {
			return matchLen, matched, tokenPair.token
		}
	}
	return 0, false, 0
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

func tokenize(code string) []*Token {
	tokenPairs := []tokenPair{
		{"if", TokenIf},
		{"else", TokenElse},
		{"while", TokenWhile},
		{"func", TokenFunc},
		{"call", TokenCall},
		{"return", TokenReturn},
		{"break", TokenBreak},
		{"continue", TokenContinue},
	}
	tokens := make([]*Token, 0)
	for i := 0; i < len(code); i++ {
		if matchLen, matched, token := matchTokens(tokenPairs, code, i); matched {
			tokens = append(tokens, NewToken(token))
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
			tokens = append(tokens, NewNumToken(int(num)))
			i += (matchLen - 1)
			continue
		}

		if matchLen, matched := matchPattern("^[A-Za-z0-9_]+", code, i); matched {
			tokens = append(tokens, NewIdentToken(code[i:i+matchLen]))
			i += (matchLen - 1)
			continue
		}

		signs := []string{"==", "!=", "<=", ">=", "<", ">", "=", "+", "-", "*", "(", ")", "{", "}", ","}
		if matchLen, matched := matchSigns(signs, code, i); matched {
			tokens = append(tokens, NewSignToken(code[i:i+matchLen]))
			i += (matchLen - 1)
			continue
		}

		if c == ';' {
			tokens = append(tokens, NewToken(TokenEOL))
			continue
		}
	}
	return tokens
}
