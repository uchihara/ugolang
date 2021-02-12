package ugolang

import (
	"fmt"
)

// CompileError dummy
type CompileError struct {
	line    int
	column  int
	message string
}

// NewCompileError dummy
func NewCompileError(pos *TokenPos, message string) *CompileError {
	return &CompileError{
		line:    pos.Line(),
		column:  pos.Column(),
		message: message,
	}
}

func (e *CompileError) Error() string {
	return fmt.Sprintf("compile error, line: %d, column: %d, message: %s", e.line, e.column, e.message)
}
