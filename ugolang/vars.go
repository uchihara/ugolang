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
		return val
	}
	if val, ok := (*v.globals)[name]; ok {
		return val
	}
	return NewNumVal(0)
}

// Set dummy
func (v NestedVarsType) Set(name string, val *Val) {
	if v.depth == 0 {
		(*v.globals)[name] = val
	} else {
		(*v.locals)[name] = val
	}
}
