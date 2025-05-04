package main

import (
	"context"
	"fmt"
)

func main() {

	//	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//	cancel()
	stop := make(chan int)

	// go func() {
	// 	fmt.Println("stop releasing")
	// 	<-stop
	// }()

	g := generateInts(ctx, 5, stop)
	readChan(ctx, g, stop)

	fmt.Println("stopped point ")
	stop <- 666
	fmt.Println("exit ")

}

func readChan(ctx context.Context, in chan int, stop chan int) { // } (cha chan int) {
	//	cha = make(chan int)
	go func() {
		//		defer close(cha)

		//	defer wg.Done()
		for c := range in {
			//	time.Sleep(400*time.Millisecond)
			select {
			case <-ctx.Done():
				fmt.Println("read cancel ON  ", c)
				return
			case a := <-in:
				fmt.Println("READed  ", a)
			}
		}
		<-stop
	}()
	//	return
}

func generateInts(ctx context.Context, n int, stop chan int) (chaGenerated chan int) {
	chaGenerated = make(chan int)
	go func() {
		defer close(chaGenerated)
		for i := range n {
			select {
			case chaGenerated <- i:
				fmt.Println("Generate ", i)
			case <-ctx.Done():
				fmt.Println("context cancel ON  ", i)
				return
			}
		}
	}()
	return
}
