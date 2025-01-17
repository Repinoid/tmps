package main

import "fmt"

type Inter interface {
	getty()
	putty()
}

type mem struct {
	a int
}
type baz struct {
	a int
}

func (t mem) getty() {
	fmt.Println("mem ", t.a)
}
func (t mem) putty() {
	fmt.Println("mem ", t.a)
}
func (t baz) getty() {
	fmt.Println("baz ", t.a)
}

func main() {
	var v Inter

	v = mem{}
	v.getty()
	v = baz{}
	v.getty()

}
