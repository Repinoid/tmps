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

	//var wg sync.WaitGroup
	//	wg.Add(1)
	in := generateInts(ctx, 5, stop)
	//	wg.Wait()

	//	wg.Add(1)
	//readCha(ctx, in, &wg)
	out := readCha(ctx, in, stop)
	_ = out
	//	wg.Wait()

	// out := duplicateChannels(ctx, in , 5, &wg)

	// go func() {
	// 	fmt.Println("len ", len(out))
	// 	for o := range out {
	// 		for c := range o {
	// 			fmt.Println("asdfg ", c)
	// 		}
	// 	}
	// 	fmt.Println("ALL <-cha, Unlock STOP ", <-stop)
	// }()

	fmt.Println("lock EXIT from main by write to STOP")
	//makeExit(ctx, stop)
	stop <- 7
	//	<-cha
	//	fmt.Println(<-cha)
	fmt.Println("all ok, exit")
	for res := range out {
		_ = res
		//fmt.Println("out ",res)
	}

}

func makeExit(ctx context.Context, stop chan int) {
	go func() {
		select {
		case <-ctx.Done():
			<-stop
			return
		}
	}()
}

func generateInts(ctx context.Context, n int, stop chan int) (cha chan int) {
	cha = make(chan int)
	go func() {
		//	defer wg.Done()
		defer func() {
			close(cha)
		}()

		for i := range n {
			select {
			case cha <- i:
				fmt.Println("Generate ", i)
			case <-ctx.Done():
				fmt.Println("general cancel ON  ", i)
				return
			}
		}
	}()
	return
}

func readCha(ctx context.Context, in chan int, stop chan int) (cha chan int) {
	cha = make(chan int)
	go func() {
		defer close(cha)

		//	defer wg.Done()
		for c := range in {
			//	time.Sleep(400*time.Millisecond)
			select {
			case <-ctx.Done():
				fmt.Println("read cancel ON  ", c)
				return
				//			case a := <-in:
			case cha <- c:
				//				cha <- a
				fmt.Println("READed  ", c)
			}
		}
		<-stop
	}()
	return
}

func duplicateChannels(ctx context.Context, in chan int, n int, wg *sync.WaitGroup) (out []chan int) {

	serie := 0
	fmt.Println("serie ", serie)
	go func() {

		for c := range in {
			fmt.Println("serie ", serie)
			for range n {
				if serie < n {
					cha := make(chan int)
					out = append(out, cha)
				}
				out[serie/n] <- c
				serie++
			}
		}
	}()
	//wg.Done()
	return
}
