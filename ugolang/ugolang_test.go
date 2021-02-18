package ugolang

import (
	"strings"
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
			code: "func main() int { 1; }",
			want: NewNumVal(1),
		},
		{
			code: "func main() int { 1;1+2; }",
			want: NewNumVal(3),
		},
		{
			code: "func main() int { var a int; a; }",
			want: NewNumVal(0),
		},
		{
			code: "func main() int { var a int; a=1; }",
			want: NewNumVal(1),
		},
		{
			code: "func main() int { var a int; a=1;a+2; }",
			want: NewNumVal(3),
		},
		{
			code: "func main() int { 1+2; }",
			want: NewNumVal(3),
		},
		{
			code: "func main() int { 1+2+3; }",
			want: NewNumVal(6),
		},
		{
			code: "func main() int { 1+2*3; }",
			want: NewNumVal(7),
		},
		{
			code: "func main() int { (1+2)*3; }",
			want: NewNumVal(9),
		},
		{
			code: "func main() int { if 1 { 1; } }",
			want: NewNumVal(1),
		},
		{
			code: "func main() int { if 0 { 1; } }",
			want: NewNumVal(0),
		},
		{
			code: "func main() int { var a int; a=0;if a { 1; } }",
			want: NewNumVal(0),
		},
		{
			code: "func main() int { var a int; a=2;if a { 1; } }",
			want: NewNumVal(1),
		},
		{
			code: "func main() int { if 0 { 1; } else { 2; } }",
			want: NewNumVal(2),
		},
		{
			code: "func main() int { if 1 { 1; } else { 2; } }",
			want: NewNumVal(1),
		},
		{
			code: "func main() int { var a int; var b int; a=1; if a { b=1; } else { b=2; } a+b; }",
			want: NewNumVal(2),
		},
		{
			code: "func main() int { var a int; var b int; a=0; if a { b=1; } else { b=2; } a+b; }",
			want: NewNumVal(2),
		},
		{
			code: "func main() int { var a int; var b int; a=0;if a{b=1;}else{b=2;}a+b; }",
			want: NewNumVal(2),
		},
		{
			code: "func main() int { 1 == 1; }",
			want: NewNumVal(1),
		},
		{
			code: "func main() int { 1 != 1; }",
			want: NewNumVal(0),
		},
		{
			code: "func main() int { var a int; a=1; if a == 1 { 2; } else { 3; } }",
			want: NewNumVal(2),
		},
		{
			code: "func main() int { var a int; a=0; if a == 1 { 2; } else { 3; } }",
			want: NewNumVal(3),
		},
		{
			code: "func main() int { if 0<=0 { 1; } else { 0; } }",
			want: NewNumVal(1),
		},
		{
			code: "func main() int { if 0<0 { 1; } else { 0; } }",
			want: NewNumVal(0),
		},
		{
			code: "func main() int { if 0>=0 { 1; } else { 0; } }",
			want: NewNumVal(1),
		},
		{
			code: "func main() int { if 0>0 { 1; } else { 0; } }",
			want: NewNumVal(0),
		},
		{
			code: "func main() int { if 0<=1 { 1; } else { 0; } }",
			want: NewNumVal(1),
		},
		{
			code: "func main() int { if 0<1 { 1; } else { 0; } }",
			want: NewNumVal(1),
		},
		{
			code: "func main() int { if 0>=1 { 1; } else { 0; } }",
			want: NewNumVal(0),
		},
		{
			code: "func main() int { if 0>1 { 1; } else { 0; } }",
			want: NewNumVal(0),
		},
		{
			code: "func main() int { var a int; a=0; while a<2 { a=a+1; } a+1; }",
			want: NewNumVal(3),
		},
		{
			code: "func main() int { var aaa int; var b_ int; var c123 int; aaa=1;b_=2;c123=3;aaa+b_+c123; }",
			want: NewNumVal(6),
		},
		{
			code: "func main() int { 123; }",
			want: NewNumVal(123),
		},
		{
			code: "func main() int { (123+456)*2; }",
			want: NewNumVal(1158),
		},
		{
			code: "func foo() int { 1; } func main() int { foo(); }",
			want: NewNumVal(1),
		},
		{
			code: "func foo() int { 1; } func main() int { foo() + 2; }",
			want: NewNumVal(3),
		},
		{
			code: "func foo() int { if 1 { return 1; } return 2; } func main() int { foo(); }",
			want: NewNumVal(1),
		},
		{
			code: "func foo() int { if 0 { return 1; } return 2; } func main() int { foo(); }",
			want: NewNumVal(2),
		},
		{
			code: "func foo() int { var a int; a=2; } func main() int { var a int; a=1; foo(); a; }",
			want: NewNumVal(1),
		},
		{
			code: "func foo() int { var a int; return a; } func main() int { var a int; a=1; foo(); }",
			want: NewNumVal(0),
		},
		{
			code: "func foo(a int) int { a; } func main() int { foo(1); }",
			want: NewNumVal(1),
		},
		{
			code: "func foo(a int, b int) int { a+b; } func main() int { foo(1, 2); }",
			want: NewNumVal(3),
		},
		{
			code: "func foo(a int) int { bar(a+1); } func bar(a int) int { a+1; } func main() int { foo(1); }",
			want: NewNumVal(3),
		},
		{
			code: "func main() int { 2-1; }",
			want: NewNumVal(1),
		},
		{
			code: "func main() int { 1-1; }",
			want: NewNumVal(0),
		},
		{
			code: "func main() int { 1-2; }",
			want: NewNumVal(-1),
		},
		{
			code: "func fib(n int) int { if n < 2 { n; } else { fib(n-2) + fib(n-1); } } func main() int { fib(0); }",
			want: NewNumVal(0),
		},
		{
			code: "func fib(n int) int { if n < 2 { n; } else { fib(n-2) + fib(n-1); } } func main() int { fib(1); }",
			want: NewNumVal(1),
		},
		{
			code: "func fib(n int) int { if n < 2 { n; } else { fib(n-2) + fib(n-1); } } func main() int { fib(2); }",
			want: NewNumVal(1),
		},
		{
			code: "func fib(n int) int { if n < 2 { n; } else { fib(n-2) + fib(n-1); } } func main() int { fib(3); }",
			want: NewNumVal(2),
		},
		{
			code: "func fib(n int) int { if n < 2 { n; } else { fib(n-2) + fib(n-1); } } func main() int { fib(4); }",
			want: NewNumVal(3),
		},
		{
			code: "func fib(n int) int { if n < 2 { n; } else { fib(n-2) + fib(n-1); } } func main() int { fib(5); }",
			want: NewNumVal(5),
		},
		{
			code: "func fib(n int) int { if n < 2 { n; } else { fib(n-2) + fib(n-1); } } func main() int { fib(6); }",
			want: NewNumVal(8),
		},
		{
			code: "var a int; a=1; func foo(a int) int { bar(a+1); } func bar(b int) int { a+b; } func main() int { var a int; a=2; foo(a); }",
			want: NewNumVal(4),
		},
		{
			code: "func main() int { var a int; a=0; while 1 { if a == 1 { break; } a=a+1; } a; }",
			want: NewNumVal(1),
		},
		{
			code: "func main() int { var a int; a=0; while 1 { if a == 2 { break; } a=a+1; continue; a=a+10; } a; }",
			want: NewNumVal(2),
		},
		{
			code: `
func main() int {
	var a int;
	a = 0;
	while a < 2 {
		a = a + 1;
	}
	return a;
}
`,
			want: NewNumVal(2),
		},
		{
			code: `
func mod(a int, b int) int {
	while a - b >= 0 {
		a = a - b;
	}
	return a;
}
func main() int {
	var s string;
	s = "";
	var i int;
	i = 1;
	while i <= 15 {
		var s2 string;
		s2 = "";
		if mod(i, 3) == 0 {
			s2 = s2 + "fizz";
		}
		if mod(i, 5) == 0 {
			s2 = s2 + "buzz";
		}
		if s2 == "" {
			s2 = s2 + sprintf("%d", i);
		}
		s = s + s2 + "\n";
		i = i + 1;
	}
	return s;
}
`,
			want: NewStrVal(`1
2
fizz
4
buzz
fizz
7
8
fizz
buzz
11
fizz
13
14
fizzbuzz
`),
		},
		{
			code: "func main() int { var a int = 1; a; }",
			want: NewNumVal(1),
		},
		{
			code: "var b int = 1; func main() int { var a int = 2; a+b; }",
			want: NewNumVal(3),
		},
		{
			code:      `func main() int { var a int = "1"; }`,
			wantError: true,
		},
		{
			code:      `func main() int { var a string = 1; }`,
			wantError: true,
		},
		{
			code:      `func main() int { var a int; a = "1"; }`,
			wantError: true,
		},
		{
			code:      `func main() int { var a string; a = 1; }`,
			wantError: true,
		},
		{
			code:      `var a int = "1"; func main() int { a; }`,
			wantError: true,
		},
		{
			code:      `func main() int { 1+"a"; }`,
			wantError: true,
		},
		{
			code:      `func main() int { var a string; 1+a; }`,
			wantError: true,
		},
		{
			code:      `func main() int { var a string = 1+"a"; }`,
			wantError: true,
		},
		{
			code:      `func main() int { foo(); }`,
			wantError: true,
		},
	}
	for _, tt := range tts {
		ugo := NewUgolang()
		if !strings.Contains(tt.code, "func mod") {
			continue
		}
		actual, err := ugo.Exec(tt.code)
		if (err != nil) != tt.wantError {
			t.Errorf("%s expect error is %v but got %s", tt.code, tt.wantError, err)
		}
		if err == nil && actual.Ne(tt.want) {
			t.Errorf("%s expect %v but got %v", tt.code, tt.want, actual)
		}
	}
}
