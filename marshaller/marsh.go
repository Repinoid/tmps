package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type strorint int

type stru struct {
	Inta int      `json:"inta"`
	Hz   strorint `json:"hz"`
}

func (si *strorint) UnmarshalJSON(b []byte) (err error) {
	var s string
	json.Unmarshal(b, &s)
	ii, err := strconv.Atoi(s)
	*si = strorint(ii)
	return err
}

func main() {

	u := stru{}
	st := []byte(`{"inta":7, "hz":"44"}`)

	json.Unmarshal(st, &u)

	fmt.Printf("%[1]v %[1]T",u.Hz+1)

}
