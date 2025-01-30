package main

import (
	"fmt"
)

func main() {

	numbers := []int{7, 9, 1, 2, 3, 4, 5}

	doneCh := make(chan struct{})
	defer close(doneCh)

	addCh := generator(numbers)

	resultCh := multiply(addCh)

	// выводим результаты
	for res := range resultCh {
		fmt.Print(res, " ")
	}
}

func generator(numbers []int) chan int {
	outputCh := make(chan int)
	go func() {
		defer close(outputCh)
		for _, num := range numbers {
			outputCh <- num
		}
	}()
	return outputCh
}

func multiply(inputCh chan int) chan int {
	resultCh := make(chan int)

	go func() {
		defer close(resultCh)

		for value := range inputCh {

			if value > 3 {
				result := value
				resultCh <- result
			}
		}
	}()

	return resultCh
}
