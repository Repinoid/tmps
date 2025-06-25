package main

import (
	"crypto/rand"
	"sort"
	"testing"
)

// func (suite *TstHeapSort) ATest01() {

// }

func BenchmarkHeapSort(b *testing.B) {
	var cmps int64

	b.StopTimer()

	baseMass = make([]byte, Long)
	rand.Read(baseMass)

	mass := make([]byte, Long)
	copy(mass, baseMass)

	b.StartTimer()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cmps++ // увеличиваем счётчик
		HeapSort(&mass)
	}
	b.ReportMetric(float64(cmps), "compares/op")
}
func BenchmarkRegularSort(b *testing.B) {
	b.StopTimer()

	baseMass = make([]byte, Long)
	rand.Read(baseMass)

	mass := make([]byte, Long)
	copy(mass, baseMass)

	b.StartTimer()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		//slices.Sort
		sort.Slice(mass, func(i, j int) bool {
			return mass[i] < mass[j]
		})
	}
}

func ABenchmarkEmpty(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// пусто
	}
}

// go test -bench . -benchmem
// go tool pprof -http=":9090" -seconds=30 http://localhost:8080/debug/pprof/profile