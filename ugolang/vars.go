package ugolang

// VarsType dummy
type VarsType map[string]int

var globalVars VarsType = map[string]int{}

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
func (v NestedVarsType) Get(name string) int {
	if val, ok := (*v.locals)[name]; ok {
		return val
	}
	return (*v.globals)[name]
}

// Set dummy
func (v NestedVarsType) Set(name string, val int) {
	if v.depth == 0 {
		(*v.globals)[name] = val
	} else {
		(*v.locals)[name] = val
	}
}
