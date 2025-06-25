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

func sortByGrt(m *[]byte) {

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

		var wg sync.WaitGroup
		wg.Add(4)
		// в начало каждой кварты перемещается максимум кварты

		go func() {
			moveMaxOnTop(&m1)
			wg.Done()
		}()
		go func() {
			moveMaxOnTop(&m2)
			wg.Done()
		}()
		go func() {
			moveMaxOnTop(&m3)
			wg.Done()
		}()
		go func() {
			moveMaxOnTop(&m4)
			wg.Done()
		}()

		wg.Wait()
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
