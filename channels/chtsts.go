package main

import (
	"context"
	"fmt"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
//	cancel()

	stop := make(chan int)
	//	go func() {
	cha := generatoras(ctx, 5)
	//	}()
	go func() {
		for c := range cha {
			fmt.Println("-->> read cha ", c)
			if c == 0 {
				cancel()
			}
		}
		fmt.Println("ALL <-cha, Unlock STOP ", <-stop)
		//	fmt.Println(<-stop)

	}()

	fmt.Println("lock EXIT from main by write to STOP")
	stop <- 7
	//	<-cha
	//	fmt.Println(<-cha)
	fmt.Println("all ok, exit")

}

func generatoras(ctx context.Context, n int) (cha chan int) {
	cha = make(chan int)
	go func() {
		defer close(cha)
		for i := range n {
			select {
			case <-ctx.Done():
				fmt.Println("cancel ON  ", i)
				return
			default:
				fmt.Println("Generate ", i)
				cha <- i
			}
		}
	}()
	return
}
