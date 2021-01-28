package ugolang

import (
  "fmt"
)

type stackType []int

var stack stackType

func (s *stackType) push(n int) {
  *s = append(*s, n)
}

func (s *stackType) pop() int {
  n := (*s)[len(*s)-1]
  *s = (*s)[0:len(*s)-1]
  return n
}

var codes []rune


func WriteCode(code string) {
  codes = make([]rune, 0)
  for _, c := range code {
    codes = append(codes, c)
  }
}

func Eval() int {
  expr()
  return stack.pop()
}

func consume(c rune) bool {
  if c == codes[0] {
    codes = codes[1:]
    return true
  }
  return false
}

func expect(c rune) {
  if !consume(c) {
    panic(fmt.Sprintf("expect %c but got %c", c, codes[0]))
  }
}

func expr() {
  add()
}

func add() {
  mul()
  for ; len(codes) > 0; {
    if consume('+') {
      mul()
      n2 := stack.pop()
      n1 := stack.pop()
      stack.push(n1+n2)
    } else {
      break
    }
  }
}

func mul() {
  pri()
  for ; len(codes) > 0; {
    if consume('*') {
      pri()
      n2 := stack.pop()
      n1 := stack.pop()
      stack.push(n1*n2)
    } else {
      break
    }
  }
}

func pri() {
  if consume('(') {
    expr()
    expect(')')
    return
  }

  num()
}

func num() {
  c := codes[0]
  if c < '0' || '9' < c {
    panic(fmt.Sprintf("expect num but got %c", c))
  }
  codes = codes[1:]
  stack.push(int(c - '0'))
}
