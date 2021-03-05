package ugolang

type stackType []string

var stack stackType

func (s *stackType) push(str string) {
	*s = append(*s, str)
}

func (s *stackType) pop() string {
	str := (*s)[len(*s)-1]
	*s = (*s)[0 : len(*s)-1]
	return str
}

func (s *stackType) peek() string {
	return (*s)[len(*s)-1]
}
