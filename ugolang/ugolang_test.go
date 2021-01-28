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
  }
  for _, tt := range tts {
    actual := Exec(tt.code)
    if actual != tt.want {
      t.Errorf("%s expect %d but got %d", tt.code, tt.want, actual)
    }
  }
}
