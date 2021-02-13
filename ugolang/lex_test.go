package ugolang

import (
	"testing"
)

func TestMatchString(t *testing.T) {
	type test struct {
		str        string
		matchedLen int
		matched    bool
		matchedStr string
		wantError  bool
	}

	tts := []test{
		{
			str:        ``,
			matchedLen: 0,
			matched:    false,
			matchedStr: ``,
		},
		{
			str:        `a`,
			matchedLen: 0,
			matched:    false,
			matchedStr: ``,
		},
		{
			str:        `"a"`,
			matchedLen: 3,
			matched:    true,
			matchedStr: `a`,
		},
		{
			str:       `"a`,
			wantError: true,
		},
		{
			str:        `"a\n"`,
			matchedLen: 5,
			matched:    true,
			matchedStr: "a\n",
		},
		{
			str:       `"a\z"`,
			wantError: true,
		},
		{
			str:        `"a\""`,
			matchedLen: 5,
			matched:    true,
			matchedStr: `a"`,
		},
		{
			str:        `"a\\"`,
			matchedLen: 5,
			matched:    true,
			matchedStr: `a\`,
		},
	}
	for _, tt := range tts {
		matchedLen, matched, matchedStr, err := matchString(tt.str)
		if matchedLen != tt.matchedLen {
			t.Errorf("%s expect matchedLen %d but got %d", tt.str, tt.matchedLen, matchedLen)
		}
		if matched != tt.matched {
			t.Errorf("%s expect matched %v but got %v", tt.str, tt.matched, matched)
		}
		if matchedStr != tt.matchedStr {
			t.Errorf("%s expect matchedStr %s but got %s", tt.str, tt.matchedStr, matchedStr)
		}
		if (err != nil) != tt.wantError {
			t.Errorf("%s expect error is %v but got %v", tt.str, tt.wantError, err)
		}
	}
}
