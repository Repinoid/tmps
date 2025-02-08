package main

import (
	"github.com/theplant/luhn"
)

func Luhner(numb int) int {
	// if luhn.Valid(numb) {
	// 	return numb
	// }
	return 10*numb + luhn.CalculateLuhn(numb)
}
