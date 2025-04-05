package main

import (
	"math/rand"
	"testing"
)

func BenchmarkGenerateGrid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = generateGrid(25, 75)
	}
}

func BenchmarkSolve(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomGenerator = rand.New(rand.NewSource(3))
		grid := generateGrid(25, 75)
		solve(grid, 75)
	}
}

// hyperfine "go run ."

// go run .
// http://localhost:6060/debug/pprof/
// http://localhost:6060/debug/pprof/ ← l’index
// http://localhost:6060/debug/pprof/profile ← pour profiler le CPU
// http://localhost:6060/debug/pprof/heap ← pour profiler la mémoire

// go test -bench=. -benchmem

