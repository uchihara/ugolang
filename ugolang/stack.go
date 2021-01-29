package ugolang

type stackType []int

var stack stackType

func (s *stackType) push(n int) {
	*s = append(*s, n)
}

func (s *stackType) pop() int {
	n := (*s)[len(*s)-1]
	*s = (*s)[0 : len(*s)-1]
	return n
}
