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
	cha := generateInts(ctx, 5)
	//	}()
	go func() {
		for c := range cha {
			fmt.Println("-->> read cha ", c)
			if c == 1 {
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

func generateInts(ctx context.Context, n int) (cha chan int) {
	cha = make(chan int)
	go func() {
		defer close(cha)
		for i := range n {
			select {
			case cha <- i:
				fmt.Println("Generate ", i)
			case <-ctx.Done():
				fmt.Println("cancel ON  ", i)
				return
			}
		}
	}()
	return
}

func duplicateChannels(ctx context.Context, in chan int, n int) (out []chan int) {
	serie := 0
	for c := range in {
		for range n {
			if serie < n {
				cha := make(chan int)
				out = append(out, cha)
			}
			out[serie/n] <- c
			serie++
		}
	}
	return
}
