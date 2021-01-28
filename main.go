package main

import (
  "fmt"
  "os"

  "github.com/uchihara/ugolang/ugolang"
)

func main() {
  code := os.Args[1]
  n := ugolang.Exec(code)
  fmt.Printf("%s=%d\n", code, n)
}
