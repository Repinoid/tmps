package main

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	// исходный код, который будем разбирать

	src := `package main
import "fmt"

func main() {
	i := 77
	defer fmt.Println(i)
	i = i + 88
}`

	// дерево разбора AST ассоциируется с набором исходных файлов FileSet
	fset := token.NewFileSet()
	// парсер может работать с файлом
	// или исходным кодом, переданным в виде строки
	f, err := parser.ParseFile(fset, "", src, 0)
	//	f, err := parser.ParseFile(fset, "./splitByChar.go", nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	// печатаем дерево
	ast.Print(fset, f)
}
