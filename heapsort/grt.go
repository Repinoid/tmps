package main

import "sync"

func sortByQuart(m *[]byte) {

	mass := *m
	initLen := len(mass)
	// если длина массива 4 то уменьшать далее не надо
	for mLen := initLen; mLen >= 4; mLen-- {
		mass = (*m)[:mLen]

		// длина подмассивов
		qNum := mLen / 4

		// создаём кварты - подмассивы, последний скорее всего огрызок
		m1 := mass[:qNum]
		m2 := mass[qNum : qNum*2]
		m3 := mass[qNum*2 : qNum*3]
		m4 := mass[qNum*3:]

		// в начало каждой кварты перемещается максимум кварты
		moveMaxOnTop(&m1)
		moveMaxOnTop(&m2)
		moveMaxOnTop(&m3)
		moveMaxOnTop(&m4)

		// определяем максимум всего массива из 4х начальных элементов (которые максимумы кварт)
		// и его индекс
		maxOf4 := m1[0]
		ind := 0
		if m2[0] > maxOf4 {
			maxOf4 = m2[0]
			ind = qNum
		}
		if m3[0] > maxOf4 {
			maxOf4 = m3[0]
			ind = qNum * 2
		}
		if m4[0] > maxOf4 {
			//	maxOf4 = m4[0]
			ind = qNum * 3
		}
		// свап максимума  с концом массива
		mass[ind], mass[mLen-1] = mass[mLen-1], mass[ind]
	}
}

func sortByGrt(m *[]byte, topper func(m *[]byte), N int) {

	//N := 4

	mass := *m
	initLen := len(mass)
	// если длина массива 4 то уменьшать далее не надо
	for mLen := initLen; mLen >= N; mLen-- {
		mass = (*m)[:mLen]

		// длина подмассивов
		qNum := mLen / N

		var wg sync.WaitGroup
		wg.Add(N)

		for i := 0; i < N; i++ {
			var mLoc []byte
			if i != N-1 {
				// определяем границы подмассива,
				mLoc = mass[qNum*i : qNum*(i+1)]
			} else {
				// последний скорее всего огрызок
				mLoc = mass[qNum*i:]
			}
			// запускаем горутины
			go func() {
				// moveMaxOnTop - в начало каждого подмассива перемещается максимум подмассива
				topper(&mLoc)
				wg.Done()
			}()

		}
		wg.Wait()
		// определяем максимум всего массива из начальных элементов подмассивов
		// и его индекс
		maxOf4 := mass[0]
		ind := 0
		for i := 1; i < N; i++ {
			if mass[i*qNum] > maxOf4 {
				maxOf4 = mass[i*qNum]
				ind = i * qNum
			}
		}
		// свап максимума  с концом массива
		mass[ind], mass[mLen-1] = mass[mLen-1], mass[ind]
	}
}

// simpleTop нахождение максимума массива и перемещение его в начало
func simpleTop(m *[]byte) {
	mass := *m
	max := mass[0]
	ind := 0
	for i := 0; i < len(mass); i++ {
		if mass[i] > max {
			max = mass[i]
			ind = i
		}
	}
	mass[ind], mass[0] = mass[0], mass[ind]
}
