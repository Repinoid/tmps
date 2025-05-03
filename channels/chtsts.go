package main

import "fmt"

func main() {
	cha := make(chan int)
	stop := make(chan int)
	go func() {
		generatoras(cha, 5)
	}()
	go func() {
		for c := range cha {
			fmt.Println("-->> read cha ", c)
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

func generatoras(cha chan int, n int) chan int {
	//cha = make(chan int)
	defer func() {
		close(cha)
		fmt.Println("cha closed by defer")
	}()

	for i := range n {
		fmt.Println("Generate ", i)
		cha <- i
	}
	return cha
}
