package ugolang

type funcMap map[string]*Node

var funcs funcMap = funcMap{}

func (f funcMap) Define(name string, body *Node) {
	f[name] = body
}

func (f funcMap) Defined(name string) bool {
	_, ok := f[name]
	return ok
}

// Frame dummy
type Frame struct {
	funcName string
}

// FuncStack dummy
type FuncStack []Frame

var funcStack FuncStack

func (s *FuncStack) push(funcName string) {
	*s = append(*s, Frame{funcName: funcName})
}

func (s *FuncStack) pop() Frame {
	frame := (*s)[len(*s)-1]
	*s = (*s)[0 : len(*s)-1]
	return frame
}

func (s FuncStack) count() int {
	return len(s)
}
