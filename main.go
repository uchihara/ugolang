package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/uchihara/ugolang/ugolang"
)

func main() {
	var dumpTokens bool
	flag.BoolVar(&dumpTokens, "tokens", false, "dump tokens")
	var dumpNodes bool
	flag.BoolVar(&dumpNodes, "nodes", false, "dump nodes")
	var implicitMain bool
	flag.BoolVar(&implicitMain, "main", false, "implicit main func")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Printf("usage %s [options] <source>\n", os.Args[0])
		flag.PrintDefaults()
		return
	}

	ugo := ugolang.NewUgolang()
	ugo.DumpTokens = dumpTokens
	ugo.DumpNodes = dumpNodes
	code := flag.Arg(0)
	if implicitMain {
		code = "func main() { " + code + " }"
	}
	n, err := ugo.Exec(code)
	if err != nil {
		fmt.Printf("%s has error: %s\n", code, err)
		return
	}
	fmt.Printf("%s=%v\n", code, n)
}
