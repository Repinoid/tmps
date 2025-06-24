package main

import (
	"crypto/rand"
	"log"
	"sort"
	"time"
)
const Long = 33333

func main() {

	// создаём случайный массив
	mass := make([]byte, Long)
	n, err := rand.Read(mass)
	if err != nil || n != Long {
		log.Printf("rand.Read rerror")
		return
	}
	mass1 := make([]byte, Long)
	copy(mass1, mass)
	log.Printf("%v %v\n", &mass[0], &mass1[0])

	//	mLen := len(mass)

	t := time.Now()

	HeapSort(&mass)

	//	log.Printf("%v\nspend ms %v\n", mass, time.Since(t))
	log.Printf("HeapSort spend %v\n", time.Since(t))

	for i := 0; i < len(mass)-2; i++ {
		if mass[i] > mass[i+1] {
			log.Println("UWAGA !")
			break
		}
	}

	t = time.Now()
	sort.Slice(mass1, func(i, j int) bool {
		return mass1[i] < mass1[j]
	})
	log.Printf("Regular spend %v\n", time.Since(t))

	for i := 0; i < len(mass1)-2; i++ {
		if mass1[i] > mass1[i+1] {
			log.Println("UWAGA !")
			break
		}
	}

}

func getMaxOnTop(m *[]byte) {
	mass := *m
	mLen := len(mass)
	cycles := mLen / 2
	i := cycles - 1
	// if massive length is чотное - only 1 kid
	if mLen%2 == 0 {
		// если элемент меньше наследника - свап
		if mass[i] < mass[i*2+1] {
			mass[i], mass[i*2+1] = mass[i*2+1], mass[i]
		}
		// если нечётное - тогда 2 потомка, увеличиваем cycle на 1
		// т.к. далее  i := cycles - 2, а надо начать с  i := cycles - 1
	} else {
		cycles = cycles + 1
	}
	for i := cycles - 2; i >= 0; i-- {
		// выбираем правого потомка как бОльшего
		k := i*2 + 2
		// если он меньше левого - уменьшаем индекс на 1, получается левый
		if mass[i*2+2] < mass[i*2+1] {
			k = k - 1
		}
		// если элемент меньше наследника - свап
		if mass[i] < mass[k] {
			mass[i], mass[k] = mass[k], mass[i]
		}
	}
	// перекидываем максимум из верхушки в самый хвост
	mass[0], mass[mLen-1] = mass[mLen-1], mass[0]
}

func HeapSort(m *[]byte) {
	mass := *m
	for mLen := len(mass); mLen > 0; mLen-- {
		cycles := mLen / 2
		i := cycles - 1
		// if massive length is чотное - only 1 kid
		if mLen%2 == 0 {
			// если элемент меньше наследника - свап
			if mass[i] < mass[i*2+1] {
				mass[i], mass[i*2+1] = mass[i*2+1], mass[i]
			}
			// если нечётное - тогда 2 потомка, увеличиваем cycle на 1
			// т.к. далее  i := cycles - 2, а надо начать с  i := cycles - 1
		} else {
			cycles = cycles + 1
		}
		for i := cycles - 2; i >= 0; i-- {
			// выбираем правого потомка как бОльшего
			k := i*2 + 2
			// если он меньше левого - уменьшаем индекс на 1, получается левый
			if mass[i*2+2] < mass[i*2+1] {
				k = k - 1
			}
			// если элемент меньше наследника - свап
			if mass[i] < mass[k] {
				mass[i], mass[k] = mass[k], mass[i]
			}
		}
		// перекидываем максимум из верхушки в самый хвост
		mass[0], mass[mLen-1] = mass[mLen-1], mass[0]
	}
}
