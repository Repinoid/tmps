package main

import "fmt"

type Inter interface {
	getty(a *int, b *string)
//	putty()
}

type mem struct {
	a int
}
type baz struct {
	b int
}

func (t mem) getty(a *int,  b *string) {
	fmt.Println("mem ", *a, b)
}
// func (t mem) putty() {
// 	fmt.Println("mem ", t.a)
// }
// func (t baz) putty() {
// 	fmt.Println("mem ", t.b)
// }

func (t baz) getty(a *int,  b *string) {
	fmt.Println("baz ", *b)
}

func main() {
	var v Inter

	i := 3
	str := "qwerty"

	v = mem{}
	v.getty(&i, nil)
	v = baz{}
	v.getty(nil, &str)

}
