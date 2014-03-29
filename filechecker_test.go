package main

import(
	"testing"
)

func BenchmarkVisitPictures(b *testing.B) {
	for i := 0; i < b.N; i++ {
		result := make(map[string]*FileInfo)
		walk("C:\\Users\\Angus\\Documents", result)
	}
}
