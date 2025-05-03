package main

import (
	"context"
	"fmt"
	"testing"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func Test_generatoras(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cha := generatoras(ctx, 5)
	for c := range cha {
		fmt.Println("-->> read cha ", c)
	}
}
