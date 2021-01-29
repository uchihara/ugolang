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
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Printf("usage %s [options] <source>\n", os.Args[0])
		flag.PrintDefaults()
		return
	}

	ugolang.DumpTokens = dumpTokens
	ugolang.DumpNodes = dumpNodes
	code := flag.Arg(0)
	n := ugolang.Exec(code)
	fmt.Printf("%s=%d\n", code, n)
}
