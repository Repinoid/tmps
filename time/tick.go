package main

import (
	"fmt"
)

func main() {

	numbers := []int{7, 9, 1, 2, 3, 4, 5}

	doneCh := make(chan struct{})
	defer close(doneCh)

	addCh := generator(doneCh, numbers)

	resultCh := multiply(doneCh, addCh)

	// выводим результаты
	for res := range resultCh {
		fmt.Print(res, " ")
	}
}

func generator(doneCh chan struct{}, numbers []int) chan int {
	outputCh := make(chan int)
	go func() {
		defer close(outputCh)
		for _, num := range numbers {
			select {
//			case <-doneCh:
//				return
			case outputCh <- num:
			}
		}
	}()
	return outputCh
}

func multiply(doneCh chan struct{}, inputCh chan int) chan int {
	resultCh := make(chan int)

	go func() {
		defer close(resultCh)

		for value := range inputCh {

			if value > 3 {
				result := value
				select {
//				case <-doneCh:
//					return
				case resultCh <- result:
				}
			}
		}
	}()

	return resultCh
}
