package main

import (
	"context"
	"fmt"
)

// add — добавляет 2 к каждому значению из inputCh и возвращает канал с результатами
func add(ctx context.Context, inputCh chan int) chan int {
	resultCh := make(chan int)

	go func() {
		defer close(resultCh)

		for value := range inputCh {
			result := value + 2

			select {
			case <-ctx.Done(): // если нужно завершить горутину
				return
			case resultCh <- result: // отправляем результат
			}
		}
	}()

	return resultCh
}

// multiply — умножает каждое значение на 3 и возвращает канал с результатами
func multiply(ctx context.Context, inputCh chan int) chan int {
	resultCh := make(chan int)

	go func() {
		defer close(resultCh)

		for value := range inputCh {
			result := value * 3

			select {
			case <-ctx.Done():
				return
			case resultCh <- result:
				_ = value
				//fmt.Println(result)
			}
		}
	}()

	return resultCh
}

// generator — отправляет данные в канал
func generator(ctx context.Context, numbers []int) chan int {
	outputCh := make(chan int)

	go func() {
		defer close(outputCh)

		for _, num := range numbers {
			select {
			case <-ctx.Done():
				return
			case outputCh <- num:

			}
		}
	}()

	return outputCh
}

func main2() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// данные, которые будем обрабатывать
	numbers := []int{1, 2, 3, 4, 5}

	// канал для остановки работы горутин
	// doneCh := make(chan struct{})
	// defer close(doneCh)

	// запускаем генератор, который отправляет числа
	inputCh := generator(ctx, numbers)

	// этапы конвейера: сначала add, потом multiply
	//	addCh := add(ctx, inputCh)
	//	resultCh := multiply(ctx, addCh)
	resultCh := multiply(ctx, inputCh)
	_ = resultCh
	//выводим результаты
	for res := range resultCh {
		_ = res
		fmt.Println(res)
	}
}
