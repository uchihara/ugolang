package ugolang

import (
	"fmt"
	"path/filepath"
	"runtime"
)

var debug bool

type callerInfo struct {
	funcName string
	funcLine int
}

func (c callerInfo) String() string {
	return fmt.Sprintf("%s:%d", c.funcName, c.funcLine)
}

func caller() callerInfo {
	pc, _, _, _ := runtime.Caller(2)
	fn := runtime.FuncForPC(pc)
	fileName, fileLine := fn.FileLine(pc)
	return callerInfo{filepath.Base(fileName), fileLine}
}

func dprintf(f string, param ...interface{}) {
	if !debug {
		return
	}
	depth := 0
	for i := 0; ; i++ {
		_, _, _, ok := runtime.Caller(i)
		if !ok {
			break
		}
		depth++
	}
	pc, _, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	fileName, fileLine := fn.FileLine(pc)
	fmt.Printf("%s:%4d: ", filepath.Base(fileName), fileLine)

	for i := 8; i < depth; i++ {
		fmt.Printf(" ")
	}

	fmt.Printf(f, param...)
}
