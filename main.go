package main

import (
  "fmt"
  "os"

  "github.com/uchihara/ugolang/ugolang"
)

func main() {
  var codes []rune
  code := os.Args[1]
  ugolang.WriteCode(code)
  n := ugolang.Eval()
  fmt.Printf("%s=%d\n", code, n)
}
