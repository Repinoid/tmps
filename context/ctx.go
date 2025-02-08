package main

import (
	"context"
	"log"
	"math/rand"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	context.AfterFunc()
	ch := make(chan int)
	rand.Seed(1000)

	go sContext(&ctx, ch, 1)
	go sContext(&ctx, ch, 2)

	go func() {
		for {
			ch <- 1 // остановка в ожидании статуса 429
			log.Println("OFF ch <- 1")
			cancel() // для срабатывания ctx.Done
		//	ctx, cancel = context.WithCancel(context.Background())
		}
	}()

	stopp := make(chan struct{})
	stopp <- struct{}{}

}

func sContext(ctx *context.Context, chi chan int, id int) {
	for {
		ra := accRural(*ctx) // эмуляция обращения к счетоводу
		if ra == 429 {       // если 429 статус - слишком много запросов
			log.Printf("id %d  ---- 429 !!!", id)
			<-chi // разблокируем канал
		}
		select {
		case <-(*ctx).Done():
			log.Println("ctx.Done", id)
		//	<-chi
		default:
			continue
		}
	}
}

func accRural(ctx context.Context) (ra int) {
	rand.Seed(time.Now().UnixNano())
	ra = 425 + rand.Intn(10)
	//	log.Println("Random  ", ra)
	time.Sleep(500 * time.Millisecond)
	return
}
