package main

import (
	"crypto/rand"
	"sort"
	"testing"
)

func (suite *TstHeapSort) ATest01() {

}

func BenchmarkHeapSort(b *testing.B) {
	b.StopTimer()

	baseMass = make([]byte, Long)
	rand.Read(baseMass)

	mass := make([]byte, Long)
	copy(mass, baseMass)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		HeapSort(&mass)
	}
}
func BenchmarkRegularSort(b *testing.B) {
	b.StopTimer()

	baseMass = make([]byte, Long)
	rand.Read(baseMass)

	mass := make([]byte, Long)
	copy(mass, baseMass)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		//slices.Sort
		sort.Slice(mass, func(i, j int) bool {
			return mass[i] < mass[j]
		})
	}
}

func BenchmarkEmpty(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// пусто
	}
}
