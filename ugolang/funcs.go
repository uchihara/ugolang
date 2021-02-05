package ugolang

// FuncType dummy
type FuncType struct {
	Name string
	Args []string
	Body *Node
}

type funcMap map[string]FuncType

var funcs funcMap = funcMap{}

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
	locals   varsType
}

// FuncStack dummy
type FuncStack []Frame

var funcStack FuncStack

func (s *FuncStack) push(funcName string) {
	frame := Frame{
		funcName: funcName,
		locals:   varsType{},
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
