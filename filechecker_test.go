package main

import (
	"log"
	"os"
	"testing"
)

func BenchmarkVisitPictures(b *testing.B) {
	var path string
	path = os.Getenv("HOME")
	if path == "" {
		path = "."
	}
	log.Printf("Using %v as root\n", path)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make(map[string]*FileInfo)
		walk(path, result)
	}
}
