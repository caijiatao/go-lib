package map_modify

import "testing"

func BenchmarkMutex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		writeToMapWithMutex()
	}
}

func BenchmarkChannel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		writeToMapWithChannel()
	}
}
