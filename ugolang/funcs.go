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
