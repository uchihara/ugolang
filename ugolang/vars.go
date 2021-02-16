package ugolang

// VarsType dummy
type VarsType map[string]*Val

var globalVars VarsType = map[string]*Val{}

// NestedVarsType dummy
type NestedVarsType struct {
	depth   int
	globals *VarsType
	locals  *VarsType
}

// NewNestedVars dummy
func NewNestedVars(depth int) *NestedVarsType {
	return &NestedVarsType{
		depth:   depth,
		globals: &globalVars,
		locals:  &VarsType{},
	}
}

// Get dummy
func (v NestedVarsType) Get(name string) *Val {
	if val, ok := (*v.locals)[name]; ok {
		dprintf("var get locally %s:%vn", name, val)
		return val
	}
	if val, ok := (*v.globals)[name]; ok {
		dprintf("var get globally %s:%vn", name, val)
		return val
	}
	dprintf("var no get %s\n", name)
	return NewNumVal(0)
}

// Set dummy
func (v NestedVarsType) Set(name string, val *Val) {
	dprintf("var set name: %s, val: %v at depth %d\n", name, val, v.depth)
	if v.depth == 0 {
		(*v.globals)[name] = val
	} else {
		(*v.locals)[name] = val
	}
}

// Define dummy
func (v NestedVarsType) Define(name string) {
	v.Set(name, NewNumVal(0))
}

// Defined dummy
func (v NestedVarsType) Defined(name string) bool {
	if ok := v.DefinedLocally(name); ok {
		return ok
	}
	_, ok := (*v.globals)[name]
	dprintf("defined globally %s %v\n", name, ok)
	return ok
}

// DefinedLocally dummy
func (v NestedVarsType) DefinedLocally(name string) bool {
	_, ok := (*v.locals)[name]
	dprintf("defined locally %s %v\n", name, ok)
	return ok
}
