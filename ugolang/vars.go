package ugolang

import (
	"fmt"
)

// VarsType dummy
type VarsType map[string]*Val

var globalVars VarsType = map[string]*Val{}

// NestedVarsType dummy
type NestedVarsType struct {
	depth   int
	globals *VarsType
	locals  *VarsType
}

func (v VarsType) String() string {
	str := fmt.Sprintf("# of vars: %d", len(v))
	for name, val := range v {
		str += ","
		str += fmt.Sprintf("%s=%v", name, val)
	}
	return str
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
		dprintf("var get locally %s:%v\n", name, val)
		return val
	}
	if val, ok := (*v.globals)[name]; ok {
		dprintf("var get globally %s:%v\n", name, val)
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
func (v NestedVarsType) Define(name string, valType ValType) {
	v.Set(name, NewDefaultVal(valType))
}

// Defined dummy
func (v NestedVarsType) Defined(name string) (ValType, bool) {
	if valType, ok := v.DefinedLocally(name); ok {
		return valType, ok
	}
	val, ok := (*v.globals)[name]
	dprintf("defined globally %s %v %v\n", name, ok, val)
	var valType ValType
	if ok {
		valType = val.Type
	}
	return valType, ok
}

// DefinedLocally dummy
func (v NestedVarsType) DefinedLocally(name string) (ValType, bool) {
	val, ok := (*v.locals)[name]
	dprintf("defined locally %s %v %v\n", name, ok, val)
	var valType ValType
	if ok {
		valType = val.Type
	}
	return valType, ok
}

// GlobalString dummy
func (v NestedVarsType) GlobalString() string {
	return v.globals.String()
}

// LocalString dummy
func (v NestedVarsType) LocalString() string {
	return v.locals.String()
}
