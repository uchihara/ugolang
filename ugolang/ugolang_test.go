package ugolang

import (
	"testing"
)

func TestUgolang(t *testing.T) {
	type test struct {
		code string
		want int
	}
	tts := []test{
		{
			code: "1;",
			want: 1,
		},
		{
			code: "1;1+2;",
			want: 3,
		},
		{
			code: "a;",
			want: 0,
		},
		{
			code: "a=1;",
			want: 1,
		},
		{
			code: "a=1;a+2;",
			want: 3,
		},
		{
			code: "1+2;",
			want: 3,
		},
		{
			code: "1+2+3;",
			want: 6,
		},
		{
			code: "1+2*3;",
			want: 7,
		},
		{
			code: "(1+2)*3;",
			want: 9,
		},
		{
			code: "if 1 { 1; }",
			want: 1,
		},
		{
			code: "if 0 { 1; }",
			want: 0,
		},
		{
			code: "a=0;if a { 1; }",
			want: 0,
		},
		{
			code: "a=2;if a { 1; }",
			want: 1,
		},
		{
			code: "if 0 { 1; } else { 2; }",
			want: 2,
		},
		{
			code: "if 1 { 1; } else { 2; }",
			want: 1,
		},
		{
			code: "a=1; if a { b=1; } else { b=2; } a+b;",
			want: 2,
		},
		{
			code: "a=0; if a { b=1; } else { b=2; } a+b;",
			want: 2,
		},
		{
			code: "a=0;if a{b=1;}else{b=2;}a+b;",
			want: 2,
		},
		{
			code: "1 == 1;",
			want: 1,
		},
		{
			code: "1 != 1;",
			want: 0,
		},
		{
			code: "a=1; if a == 1 { 2; } else { 3; }",
			want: 2,
		},
		{
			code: "a=0; if a == 1 { 2; } else { 3; }",
			want: 3,
		},
		{
			code: "if 0<=0 { 1; } else { 0; }",
			want: 1,
		},
		{
			code: "if 0<0 { 1; } else { 0; }",
			want: 0,
		},
		{
			code: "if 0>=0 { 1; } else { 0; }",
			want: 1,
		},
		{
			code: "if 0>0 { 1; } else { 0; }",
			want: 0,
		},
		{
			code: "if 0<=1 { 1; } else { 0; }",
			want: 1,
		},
		{
			code: "if 0<1 { 1; } else { 0; }",
			want: 1,
		},
		{
			code: "if 0>=1 { 1; } else { 0; }",
			want: 0,
		},
		{
			code: "if 0>1 { 1; } else { 0; }",
			want: 0,
		},
		{
			code: "a=0; while a<2 { a=a+1; } a+1;",
			want: 3,
		},
		{
			code: "aaa=1;b_=2;c123=3;aaa+b_+c123;",
			want: 6,
		},
		{
			code: "123;",
			want: 123,
		},
		{
			code: "(123+456)*2;",
			want: 1158,
		},
	}
	for _, tt := range tts {
		ugo := NewUgolang()
		actual := ugo.Exec(tt.code)
		if actual != tt.want {
			t.Errorf("%s expect %d but got %d", tt.code, tt.want, actual)
		}
	}
}
