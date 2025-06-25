package main

import (
	"crypto/rand"
	"log"
	"net/http"
	_ "net/http/pprof" // подключаем пакет pprof
	"sort"
	"time"
)

const Long = 100_000

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

	t := time.Now()
	HeapSort(&mass)
	log.Printf("HeapSort spend %v\n", time.Since(t))

	for i := 0; i < len(mass)-2; i++ {
		if mass[i] > mass[i+1] {
			log.Println("UWAGA heapsort !")
			break
		}
	}

	mass2 := make([]byte, Long)
	copy(mass2, mass)
	t = time.Now()
	sortByMaxFunc(&mass2)
	log.Printf("sortByMaxFunc spend %v\n", time.Since(t))

	for i := 0; i < len(mass2)-2; i++ {
		if mass2[i] > mass2[i+1] {
			log.Println("UWAGA sortByMaxFunc !")
			break
		}
	}

	mass3 := make([]byte, Long)
	copy(mass3, mass)
	t = time.Now()
	sortByQuart(&mass3)
	log.Printf("sortByMaxQuart spend %v\n", time.Since(t))

	for i := 0; i < len(mass3)-2; i++ {
		if mass3[i] > mass3[i+1] {
			log.Println("UWAGA sortByMaxFunc !")
			break
		}
	}

	massa := make([]byte, Long)
	copy(massa, mass)
	t = time.Now()
	sortByGrt(&massa, simpleTop, 5)
	//	sortByGrt(&massa, moveMaxOnTop, 5)
	log.Printf("sort By Gouroutines spend %v\n", time.Since(t))

	for i := 0; i < len(massa)-2; i++ {
		if massa[i] > massa[i+1] {
			log.Println("UWAGA sortByMaxFunc !")
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
			log.Println("UWAGA sort.Slice!")
			break
		}
	}

	http.ListenAndServe(":8080", nil) // запускаем сервер

}

func sortByMaxFunc(mass *[]byte) {
	//mass := *m
	for mLen := len(*mass); mLen > 1; mLen-- {
		m := (*mass)[:mLen]
		moveMaxOnTop(&m)
		m[0], m[mLen-1] = m[mLen-1], m[0]
	}
}

// getMaxOnTop перемещает максимальный элемент в начало массива
func moveMaxOnTop(m *[]byte) {
	mass := *m
	mLen := len(mass)
	cycles := mLen / 2
	// if cycles == 0 {
	// 	return
	// }
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
			if mass[k-1] > mass[k] {
				k = k - 1
			}
			// если элемент меньше наследника - свап
			//			if mass[i] < mass[k] {
			if mass[k] > mass[i] {
				mass[i], mass[k] = mass[k], mass[i]
			}
		}
		// перекидываем максимум из верхушки в самый хвост
		mass[0], mass[mLen-1] = mass[mLen-1], mass[0]
	}
}
