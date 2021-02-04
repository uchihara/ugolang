package ugolang

type varsType map[string]int

func (v varsType) Get(name string) int {
	return v[name]
}

func (v varsType) Set(name string, val int) {
	v[name] = val
}
