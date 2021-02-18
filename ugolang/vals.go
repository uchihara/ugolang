package ugolang

import (
	"fmt"
)

// ValType dummy
type ValType int

const (
	// NumVal dummy
	NumVal ValType = iota + 1
	// StrVal dummy
	StrVal
)

func (v ValType) String() string {
	switch v {
	case NumVal:
		return "numVal"
	case StrVal:
		return "strVal"
	default:
		return fmt.Sprintf("unknown valType: %d", v)
	}
}

// Val dummy
type Val struct {
	Type ValType
	Num  int
	Str  string
}

// NewNumVal dummy
func NewNumVal(num int) *Val {
	return &Val{
		Type: NumVal,
		Num:  num,
	}
}

// NewStrVal dummy
func NewStrVal(str string) *Val {
	return &Val{
		Type: StrVal,
		Str:  str,
	}
}

// NewDefaultVal dummy
func NewDefaultVal(valType ValType) *Val {
	if valType == StrVal {
		return NewStrVal("")
	}
	return NewNumVal(0)
}

func (v Val) String() string {
	switch v.Type {
	case NumVal:
		return fmt.Sprintf("num(%d)", v.Num)
	case StrVal:
		return fmt.Sprintf("str(%s)", v.Str)
	default:
		return fmt.Sprintf("unknown type %d", v.Type)
	}
}

// Add dummy
func (v *Val) Add(other *Val) *Val {
	if v.Type == StrVal && v.Type == other.Type {
		return NewStrVal(v.Str + other.Str)
	}
	return NewNumVal(v.Num + other.Num)
}

// Sub dummy
func (v *Val) Sub(other *Val) *Val {
	return NewNumVal(v.Num - other.Num)
}

// Mul dummy
func (v *Val) Mul(other *Val) *Val {
	return NewNumVal(v.Num * other.Num)
}

// Eq dummy
func (v *Val) Eq(other *Val) bool {
	if v.Type == StrVal {
		return v.Str == other.Str
	}
	return v.Num == other.Num
}

// Ne dummy
func (v *Val) Ne(other *Val) bool {
	return !v.Eq(other)
}
