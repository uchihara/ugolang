package ugolang

import (
	"fmt"
	"regexp"
	"strconv"
)

func matchToken(token, code string) (int, bool) {
	tokenLen := len(token)
	if tokenLen+1 > len(code) {
		return 0, false
	}
	if code[0:tokenLen] != token {
		return 0, false
	}
	nextChar := code[tokenLen]
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
	fn      func(line, col int) *Token
}

func matchTokens(tokenPairs []tokenPair, code string) (int, bool, int) {
	for idx, tokenPair := range tokenPairs {
		if matchLen, matched := matchToken(tokenPair.keyword, code); matched {
			return matchLen, matched, idx
		}
	}
	return 0, false, 0
}

func matchSign(sign, code string) (int, bool) {
	signLen := len(sign)
	if signLen > len(code) {
		return 0, false
	}
	if code[0:signLen] != sign {
		return 0, false
	}
	return len(sign), true
}

func matchSigns(signs []string, code string) (int, bool) {
	for _, sign := range signs {
		if matchLen, matched := matchSign(sign, code); matched {
			return matchLen, matched
		}
	}
	return 0, false
}

func matchPattern(pattern, code string) (int, bool) {
	re := regexp.MustCompile(pattern)
	loc := re.FindStringIndex(code[0:])
	if len(loc) != 2 {
		return 0, false
	}
	return loc[1] - loc[0], true
}

func matchString(code string) (int, bool, string, error) {
	if len(code) == 0 {
		return 0, false, "", nil
	}
	if code[0] != '"' {
		return 0, false, "", nil
	}
	matchedStr := make([]rune, 0)

	pos := 1
	for {
		if pos >= len(code) {
			return 0, false, "", fmt.Errorf("unquoted string")
		}
		c := rune(code[pos])
		if c == '"' {
			pos++
			break
		}
		if c == '\\' {
			if pos >= len(code) {
				return 0, false, "", fmt.Errorf("unexpected escape character")
			}
			pos++
			var escaped rune
			switch code[pos] {
			case 'n':
				escaped = '\n'
			case 'r':
				escaped = '\r'
			case 't':
				escaped = '\t'
			case 'b':
				escaped = '\b'
			case '"':
				escaped = '"'
			case '\\':
				escaped = '\\'
			default:
				return 0, false, "", fmt.Errorf("unexpected escape character: %c", code[pos])
			}
			matchedStr = append(matchedStr, escaped)
			pos++
			continue
		}
		matchedStr = append(matchedStr, c)
		pos++
	}
	return pos, true, string(matchedStr), nil
}

func tokenize(code string) ([]*Token, error) {
	tokenPairs := []tokenPair{
		{"var", func(line, col int) *Token { return NewToken(line, col, TokenVar) }},
		{"int", func(line, col int) *Token { return NewValTypeToken(line, col, NumVal) }},
		{"string", func(line, col int) *Token { return NewValTypeToken(line, col, StrVal) }},
		{"if", func(line, col int) *Token { return NewToken(line, col, TokenIf) }},
		{"else", func(line, col int) *Token { return NewToken(line, col, TokenElse) }},
		{"while", func(line, col int) *Token { return NewToken(line, col, TokenWhile) }},
		{"func", func(line, col int) *Token { return NewToken(line, col, TokenFunc) }},
		{"return", func(line, col int) *Token { return NewToken(line, col, TokenReturn) }},
		{"break", func(line, col int) *Token { return NewToken(line, col, TokenBreak) }},
		{"continue", func(line, col int) *Token { return NewToken(line, col, TokenContinue) }},
	}
	line := 1
	col := 1
	pos := 0
	tokens := make([]*Token, 0)
	for pos < len(code) {
		c := code[pos]
		if c == ' ' || c == '\t' {
			pos++
			col++
			continue
		}

		if c == '\n' {
			line++
			pos++
			col = 1
			continue
		}

		if c == ';' {
			pos++
			col++
			tokens = append(tokens, NewToken(line, col, TokenEOL))
			continue
		}

		if matchLen, matched, matchedIdx := matchTokens(tokenPairs, code[pos:]); matched {
			token := tokenPairs[matchedIdx].fn(line, col)
			tokens = append(tokens, token)
			pos += matchLen
			col += matchLen
			continue
		}

		if matchLen, matched := matchPattern("^[0-9]+", code[pos:]); matched {
			numStr := code[pos : pos+matchLen]
			num, err := strconv.ParseInt(numStr, 10, 64)
			if err != nil {
				return nil, NewCompileError(NewTokenPos(line, col), fmt.Sprintf("invalid num format: %s", numStr))
			}
			tokens = append(tokens, NewValToken(line, col, NewNumVal(int(num))))
			pos += matchLen
			col += matchLen
			continue
		}

		if matchLen, matched := matchPattern("^[A-Za-z0-9_]+", code[pos:]); matched {
			tokens = append(tokens, NewIdentToken(line, col, code[pos:pos+matchLen]))
			pos += matchLen
			col += matchLen
			continue
		}

		signs := []string{"==", "!=", "<=", ">=", "<", ">", "=", "+", "-", "*", "(", ")", "{", "}", ","}
		if matchLen, matched := matchSigns(signs, code[pos:]); matched {
			tokens = append(tokens, NewSignToken(line, col, code[pos:pos+matchLen]))
			pos += matchLen
			col += matchLen
			continue
		}

		if matchLen, matched, str, err := matchString(code[pos:]); matched || err != nil {
			if err != nil {
				return nil, NewCompileError(NewTokenPos(line, col), err.Error())
			}
			tokens = append(tokens, NewValToken(line, col, NewStrVal(str)))
			pos += matchLen
			col += matchLen
			continue
		}

		return nil, NewCompileError(NewTokenPos(line, col), fmt.Sprintf("unknown character found: %c", c))
	}
	return tokens, nil
}
