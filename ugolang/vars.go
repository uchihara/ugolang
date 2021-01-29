package ugolang

type varsType map[rune]int

var vars varsType = map[rune]int{}

func (v varsType) Get(name rune) int {
	return v[name]
}

func (v varsType) Set(name rune, val int) {
	v[name] = val
}
