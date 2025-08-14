package main

import (
	"context"
	"fmt"
	"sync"
)

func main() {
	//	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//	cancel()
	stop := make(chan int)
	defer close(stop)

	// go func() {
	// 	fmt.Println("stop releasing")
	// 	<-stop
	// }()

	g := make(chan int)
	var wg sync.WaitGroup

	wg.Add(1)
	g1 := generateInts(ctx, 3, 10)

	go func() {
		for a := range g1 {
			g <- a
			fmt.Println("g1  ", a)
		}
		wg.Done()
	}()

	wg.Add(1)
	g2 := generateInts(ctx, 2, 20)
	go func() {
		for a := range g2 {
			g <- a
			fmt.Println("g2  ", a)
		}
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(g)
	}()

	out := readChan(ctx, g, stop)

	go func() {
		for a:= range out {
			fmt.Println("received.......", a)
		}
	}()

	fmt.Println("stopped point ")
	stop <- 666

	fmt.Println("exit ")

}

// readChan читает из  in отправляет в out
func readChan(ctx context.Context, in chan int, stop chan int) chan int {
	out := make(chan int)
	go func() {
		defer close(out)

		for c := range in {
			select {
			case <-ctx.Done():
				fmt.Println("read cancel ON  ", c)
				<-stop
				return
			case out <- c:
				fmt.Println("READed and sent to OUT ", c)
			}
		}
		<-stop
	}()
	return out
}

// создаём канал и засылам в chaGenerated n значений, цикл от 0 до n-1 + offset
func generateInts(ctx context.Context, n, offset int) chan int {
	chaGenerated := make(chan int)
	go func() {
		defer close(chaGenerated)
		for i := range n {
			select {
			case chaGenerated <- (i + offset):
				fmt.Println("Generate ", i+offset)
			case <-ctx.Done():
				fmt.Println("context cancel ON  ", i)
				return
			}
		}
	}()
	return chaGenerated
}
