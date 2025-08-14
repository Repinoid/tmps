package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//	cancel()
	stop := make(chan int)
	defer close(stop)

	g := make(chan int)
	var wg sync.WaitGroup
	var readers sync.WaitGroup

	g1 := generateInts(ctx, 30, 10)
	g2 := generateInts(ctx, 20, 20)

	fanInOut(g, &wg, g1, g2)

	go func() {
		wg.Wait()
		close(g)
	}()

	out := readChan(ctx, g, stop)

	receiver(out, &readers, 1)
	receiver(out, &readers, 2)
	receiver(out, &readers, 3)
	receiver(out, &readers, 4)

	fmt.Println("stopped point ")
	stop <- 666
	readers.Wait()

	fmt.Println("exit ")
}

func receiver(out chan int, readers *sync.WaitGroup, num int) {
	readers.Add(1)
	go func() {
		defer readers.Done()
		for a := range out {
			fmt.Println("receiver .......", a, num)
			time.Sleep(1 * time.Millisecond)
		}
	}()

}

// inout засылает в out вычитывая из каналов/канала ins
func fanInOut(out chan<- int, wg *sync.WaitGroup, ins ...<-chan int) {
	for _, in := range ins {
		wg.Add(1)
		go func(ch <-chan int) {
			defer wg.Done()
			for v := range ch {
				out <- v
			}
		}(in)
	}
}

// readChan читает из  in отправляет в out
func readChan(ctx context.Context, in, stop chan int) chan int {
	out := make(chan int)
	go func(in, out chan int) {
		defer close(out)
		for c := range in {
			select {
			case <-ctx.Done():
				fmt.Println("read cancel ON  ", c)
				<-stop
				return
			case out <- c:
			}
		}
		<-stop
	}(in, out)
	return out
}

// создаём канал и засылам в chaGenerated n значений, цикл от 0 до n-1 + offset
func generateInts(ctx context.Context, n, offset int) chan int {
	chaGenerated := make(chan int)
	go func() {
		defer close(chaGenerated)
		for i := 0; i < n; i++ {
			select {
			case chaGenerated <- (i + offset):
				// fmt.Println("Generate ", i+offset)
			case <-ctx.Done():
				fmt.Println("context cancel ON  ", i)
				return
			}
		}
	}()
	return chaGenerated
}
