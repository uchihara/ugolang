package ugolang

import (
	"fmt"
)

// FuncType dummy
type FuncType struct {
	IsNative   bool
	Name       string
	Args       []string
	Body       *Node
	nativeFunc func(args []interface{}) *Val
}

type funcMap map[string]FuncType

var funcs funcMap = funcMap{
	"printf": {
		Name:     "printf",
		IsNative: true,
		nativeFunc: func(args []interface{}) *Val {
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
	"sprintf": {
		Name:     "sprintf",
		IsNative: true,
		nativeFunc: func(args []interface{}) *Val {
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

// CallNative dummy
func (f FuncType) CallNative(name string, vals []*Val) *Val {
	fn := funcs[name]
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

func (f funcMap) Define(name string, args []string, body *Node) {
	f[name] = FuncType{
		Name: name,
		Args: args,
		Body: body,
	}
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
