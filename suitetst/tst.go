package main

import (
	"encoding/json"
	"fmt"
)

type OrderStatus struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func main() {

	b := OrderStatus{Order: "orda", Status: "norm"}

	bjson, _ := json.Marshal(b)

	var c OrderStatus

	err := json.Unmarshal(bjson, &c)

	fmt.Printf("%v %v %T\n", c.Accrual+2, err, c.Accrual)

}
