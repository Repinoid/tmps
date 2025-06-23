package main

import (
	"crypto/rand"
	"log"
)

func main() {

	// создаём случайный массив
	const Long = 33
	mass := make([]byte, Long)
	n, err := rand.Read(mass)
	if err != nil || n != Long {
		log.Printf("rand.Read rerror")
		return
	}

	mLen := len(mass)
	cycles := mLen / 2

	for i := cycles - 1; i >= 0; i-- {
		// if massive length is odd - only 1 kid
		if mLen%2 == 1 {
			// если элемент меньше наследника - свап
			if mass[i] < mass[i*2+1] {
				mass[i], mass[i*2+1] = mass[i*2+1], mass[i]
			}
			continue
		}
		// выбираем правого потомка как бОльшего
		k := i*2+2
		// если он меньше левого - уменьшаем индекс на 1, получается левый
		if mass[i*2+2] < mass[i*2+1] {
			k = k - 1
		}
		// если элемент меньше наследника - свап
		if mass[i] < mass[k] {
			mass[i], mass[k] = mass[k], mass[i]
		}
	}
	log.Printf("top element is %v\n", mass[0])

}
