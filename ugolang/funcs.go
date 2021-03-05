package ugolang

import (
	"fmt"
)

// FuncType dummy
type FuncType struct {
	Name       string
	Args       []*Node
	RetValType ValType
	Body       *Node
	Vars       *NestedVarsType
	IsNative   bool
	nativeFunc func(args []interface{}) *Val
}

type funcMap map[string]FuncType

var funcTable funcMap

type nativeFunc struct {
	name       string
	retValType ValType
	fn         func([]interface{}) *Val
}

var nativeFuncs = []nativeFunc{
	{
		name:       "printf",
		retValType: NumVal,
		fn: func(args []interface{}) *Val {
			if len(args) == 0 {
				return NewNumVal(0)
			}
			f, ok := args[0].(string)
			if !ok {
				return NewNumVal(0)
			}
			params := make([]interface{}, 0)
			if len(args) > 1 {
				params = args[1:]
			}
			fmt.Printf(f, params...)
			return NewNumVal(0)
		},
	},
	{
		name:       "sprintf",
		retValType: StrVal,
		fn: func(args []interface{}) *Val {
			if len(args) == 0 {
				return NewStrVal("")
			}
			f, ok := args[0].(string)
			if !ok {
				return NewStrVal("")
			}
			params := make([]interface{}, 0)
			if len(args) > 1 {
				params = args[1:]
			}
			r := fmt.Sprintf(f, params...)
			return NewStrVal(r)
		},
	},
}

// InitFuncs dummy
func InitFuncs() {
	funcTable = funcMap{
		".global": FuncType{
			Vars: NewNestedVars(0),
		},
	}
	for _, nativeFunc := range nativeFuncs {
		funcTable[nativeFunc.name] = FuncType{
			IsNative:   true,
			Name:       nativeFunc.name,
			RetValType: nativeFunc.retValType,
			nativeFunc: nativeFunc.fn,
		}
	}
}

// CallNative dummy
func (f FuncType) CallNative(name string, vals []*Val) *Val {
	fn := funcTable[name]
	params := make([]interface{}, 0)
	for _, v := range vals {
		if v.Type == NumVal {
			params = append(params, v.Num)
		} else {
			params = append(params, v.Str)
		}
	}
	return fn.nativeFunc(params)
}

func (f funcMap) Define(name string) {
	f[name] = FuncType{
		Name: name,
		Vars: NewNestedVars(1),
	}
}

func (f funcMap) Set(name string, args []*Node, retValType ValType, body *Node) {
	m := f[name]
	m.Args = args
	m.RetValType = retValType
	m.Body = body
	f[name] = m
}

func (f funcMap) Lookup(name string) (FuncType, bool) {
	fn, ok := f[name]
	return fn, ok
}

func (f funcMap) Defined(name string) bool {
	_, ok := f[name]
	return ok
}

// Frame dummy
type Frame struct {
	funcName string
	vars     *NestedVarsType
}

// FuncStack dummy
type FuncStack []Frame

var funcStack FuncStack

func (s FuncStack) String() string {
	str := ""
	for i, stack := range s {
		if str != "" {
			str += ","
		}
		str += fmt.Sprintf("%d: %v", i, stack)
	}
	return str
}

func (f Frame) String() string {
	return fmt.Sprintf("funcName: %s, var: %s", f.funcName, f.vars.LocalString())
}

func (s *FuncStack) reset() {
	(*s) = []Frame{}
}

func (s *FuncStack) push(funcName string) {
	frame := Frame{
		funcName: funcName,
		vars:     NewNestedVars(s.count()),
	}
	*s = append(*s, frame)
}

func (s *FuncStack) pop() Frame {
	frame := (*s)[len(*s)-1]
	*s = (*s)[0 : len(*s)-1]
	return frame
}

func (s FuncStack) peek() Frame {
	return s[len(s)-1]
}

func (s FuncStack) count() int {
	return len(s)
}
