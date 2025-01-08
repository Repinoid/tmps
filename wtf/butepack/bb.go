package main

import (
	"fmt"
	"strconv"

	"go.uber.org/zap"
)

var sugar zap.SugaredLogger

func main() {

	d := map[string]int{}

	fmt.Println(wis(5))
	fmt.Println(wis(5.5))
	fmt.Println(wisa("str"))
	fmt.Println(wis(d))

}

func wisa(wa interface{}) string {
	switch fmt.Sprintf("%T", wa) {
	case "string":
		return "string " + wa.(string)
	default:
		return "hz"
	}
}

func wis(wa interface{}) string {
	switch wa.(type) {
	case string:
		return "string " + wa.(string)
	case float64:
		return "float " + retflo(wa)
	case int64, int32, int:
		return "int"
	default:
		return "hz"
	}
}
func retflo(flo interface{}) string {
	return strconv.FormatFloat(flo.(float64), 'f', -1, 64)
}

// assigning the result of this type assertion to a variable (switch wa := wa.(type)) could eliminate type assertions in switch cases (S1034)
