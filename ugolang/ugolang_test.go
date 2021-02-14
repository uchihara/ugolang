package ugolang

import (
	"testing"
)

func TestUgolang(t *testing.T) {
	type test struct {
		code      string
		want      *Val
		wantError bool
	}
	tts := []test{
		{
			code:      ";",
			wantError: true,
		},
		{
			code: "func main() { 1; }",
			want: NewNumVal(1),
		},
		{
			code: "func main() { 1;1+2; }",
			want: NewNumVal(3),
		},
		{
			code: "func main() { a; }",
			want: NewNumVal(0),
		},
		{
			code: "func main() { a=1; }",
			want: NewNumVal(1),
		},
		{
			code: "func main() { a=1;a+2; }",
			want: NewNumVal(3),
		},
		{
			code: "func main() { 1+2; }",
			want: NewNumVal(3),
		},
		{
			code: "func main() { 1+2+3; }",
			want: NewNumVal(6),
		},
		{
			code: "func main() { 1+2*3; }",
			want: NewNumVal(7),
		},
		{
			code: "func main() { (1+2)*3; }",
			want: NewNumVal(9),
		},
		{
			code: "func main() { if 1 { 1; } }",
			want: NewNumVal(1),
		},
		{
			code: "func main() { if 0 { 1; } }",
			want: NewNumVal(0),
		},
		{
			code: "func main() { a=0;if a { 1; } }",
			want: NewNumVal(0),
		},
		{
			code: "func main() { a=2;if a { 1; } }",
			want: NewNumVal(1),
		},
		{
			code: "func main() { if 0 { 1; } else { 2; } }",
			want: NewNumVal(2),
		},
		{
			code: "func main() { if 1 { 1; } else { 2; } }",
			want: NewNumVal(1),
		},
		{
			code: "func main() { a=1; if a { b=1; } else { b=2; } a+b; }",
			want: NewNumVal(2),
		},
		{
			code: "func main() { a=0; if a { b=1; } else { b=2; } a+b; }",
			want: NewNumVal(2),
		},
		{
			code: "func main() { a=0;if a{b=1;}else{b=2;}a+b; }",
			want: NewNumVal(2),
		},
		{
			code: "func main() { 1 == 1; }",
			want: NewNumVal(1),
		},
		{
			code: "func main() { 1 != 1; }",
			want: NewNumVal(0),
		},
		{
			code: "func main() { a=1; if a == 1 { 2; } else { 3; } }",
			want: NewNumVal(2),
		},
		{
			code: "func main() { a=0; if a == 1 { 2; } else { 3; } }",
			want: NewNumVal(3),
		},
		{
			code: "func main() { if 0<=0 { 1; } else { 0; } }",
			want: NewNumVal(1),
		},
		{
			code: "func main() { if 0<0 { 1; } else { 0; } }",
			want: NewNumVal(0),
		},
		{
			code: "func main() { if 0>=0 { 1; } else { 0; } }",
			want: NewNumVal(1),
		},
		{
			code: "func main() { if 0>0 { 1; } else { 0; } }",
			want: NewNumVal(0),
		},
		{
			code: "func main() { if 0<=1 { 1; } else { 0; } }",
			want: NewNumVal(1),
		},
		{
			code: "func main() { if 0<1 { 1; } else { 0; } }",
			want: NewNumVal(1),
		},
		{
			code: "func main() { if 0>=1 { 1; } else { 0; } }",
			want: NewNumVal(0),
		},
		{
			code: "func main() { if 0>1 { 1; } else { 0; } }",
			want: NewNumVal(0),
		},
		{
			code: "func main() { a=0; while a<2 { a=a+1; } a+1; }",
			want: NewNumVal(3),
		},
		{
			code: "func main() { aaa=1;b_=2;c123=3;aaa+b_+c123; }",
			want: NewNumVal(6),
		},
		{
			code: "func main() { 123; }",
			want: NewNumVal(123),
		},
		{
			code: "func main() { (123+456)*2; }",
			want: NewNumVal(1158),
		},
		{
			code: "func foo() { 1; } func main() { call foo(); }",
			want: NewNumVal(1),
		},
		{
			code: "func foo() { 1; } func main() { call foo() + 2; }",
			want: NewNumVal(3),
		},
		{
			code: "func foo() { if 1 { return 1; } return 2; } func main() { call foo(); }",
			want: NewNumVal(1),
		},
		{
			code: "func foo() { if 0 { return 1; } return 2; } func main() { call foo(); }",
			want: NewNumVal(2),
		},
		{
			code: "func foo() { a=2; } func main() { a=1; call foo(); a; }",
			want: NewNumVal(1),
		},
		{
			code: "func foo() { return a; } func main() { a=1; call foo(); }",
			want: NewNumVal(0),
		},
		{
			code: "func foo(a) { a; } func main() { call foo(1); }",
			want: NewNumVal(1),
		},
		{
			code: "func foo(a, b) { a+b; } func main() { call foo(1, 2); }",
			want: NewNumVal(3),
		},
		{
			code: "func foo(a) { call bar(a+1); } func bar(a) { a+1; } func main() { call foo(1); }",
			want: NewNumVal(3),
		},
		{
			code: "func main() { 2-1; }",
			want: NewNumVal(1),
		},
		{
			code: "func main() { 1-1; }",
			want: NewNumVal(0),
		},
		{
			code: "func main() { 1-2; }",
			want: NewNumVal(-1),
		},
		{
			code: "func fib(n) { if n < 2 { n; } else { call fib(n-2) + call fib(n-1); } } func main() { call fib(0); }",
			want: NewNumVal(0),
		},
		{
			code: "func fib(n) { if n < 2 { n; } else { call fib(n-2) + call fib(n-1); } } func main() { call fib(1); }",
			want: NewNumVal(1),
		},
		{
			code: "func fib(n) { if n < 2 { n; } else { call fib(n-2) + call fib(n-1); } } func main() { call fib(2); }",
			want: NewNumVal(1),
		},
		{
			code: "func fib(n) { if n < 2 { n; } else { call fib(n-2) + call fib(n-1); } } func main() { call fib(3); }",
			want: NewNumVal(2),
		},
		{
			code: "func fib(n) { if n < 2 { n; } else { call fib(n-2) + call fib(n-1); } } func main() { call fib(4); }",
			want: NewNumVal(3),
		},
		{
			code: "func fib(n) { if n < 2 { n; } else { call fib(n-2) + call fib(n-1); } } func main() { call fib(5); }",
			want: NewNumVal(5),
		},
		{
			code: "func fib(n) { if n < 2 { n; } else { call fib(n-2) + call fib(n-1); } } func main() { call fib(6); }",
			want: NewNumVal(8),
		},
		{
			code: "a=1; func foo(a) { call bar(a+1); } func bar(b) { a+b; } func main() { a=2; call foo(a); }",
			want: NewNumVal(4),
		},
		{
			code: "func main() { a=0; while 1 { if a == 1 { break; } a=a+1; } a; }",
			want: NewNumVal(1),
		},
		{
			code: "func main() { a=0; while 1 { if a == 2 { break; } a=a+1; continue; a=a+10; } a; }",
			want: NewNumVal(2),
		},
		{
			code: `
func main() {
	a = 0;
	while a < 2 {
		a = a + 1;
	}
	return a;
}
`,
			want: NewNumVal(2),
		},
	}
	for _, tt := range tts {
		ugo := NewUgolang()
		actual, err := ugo.Exec(tt.code)
		if (err != nil) != tt.wantError {
			t.Errorf("%s expect error is %v but got %s", tt.code, tt.wantError, err)
		}
		if err == nil && actual.Ne(tt.want) {
			t.Errorf("%s expect %v but got %v", tt.code, tt.want, actual)
		}
	}
}
