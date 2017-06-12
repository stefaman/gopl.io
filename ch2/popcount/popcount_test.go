// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// package popcount_test
package popcount

import (
	"testing"

	// "gopl.io/ch2/popcount"
	"sync"
	"time"
	"fmt"
)

// -- Alternative implementations --

func BitCount(x uint64) int {
	// Hacker's Delight, Figure 5-2.
	x = x - ((x >> 1) & 0x5555555555555555)
	x = (x & 0x3333333333333333) + ((x >> 2) & 0x3333333333333333)
	x = (x + (x >> 4)) & 0x0f0f0f0f0f0f0f0f
	x = x + (x >> 8)
	x = x + (x >> 16)
	x = x + (x >> 32)
	return int(x & 0x7f)
}

func PopCountByClearing(x uint64) int {
	n := 0
	for x != 0 {
		x = x & (x - 1) // clear rightmost non-zero bit
		n++
	}
	return n
}

func PopCountByShifting(x uint64) int {
	n := 0
	for i := uint(0); i < 64; i++ {
		if x&(1<<i) != 0 {
			n++
		}
	}
	return n
}

func PopCountByRecurring(x uint64) uint64 {
	if x == 0 {
		return 0
	}
	return PopCountByRecurring(x>>1) + x & 1
}



func TestPopCountByRecurring(t *testing.T) {
	n := PopCountByRecurring(0x1234567890ABCDEF)
	if n != 32 {
		t.Errorf("PopCountByRecurring failed. Got %v, expected 32", n)
	}
}

// -- Benchmarks --

func BenchmarkPopCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PopCount(0x1234567890ABCDEF)
	}
}

func BenchmarkBitCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BitCount(0x1234567890ABCDEF)
	}
}

func BenchmarkPopCountByClearing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PopCountByClearing(0x1234567890ABCDEF)
	}
}

func BenchmarkPopCountByShifting(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PopCountByShifting(0x1234567890ABCDEF)
	}
}

func BenchmarkPopCountByRecurring(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PopCountByRecurring(0x1234567890ABCDEF)
	}
}

//e11.6 p323
func benchmark(b *testing.B, f func(uint64) int, x uint64) {
	for i := 0; i < b.N; i++ {
		f(x)
	}
}
func benchmarkAvr(f func(uint64) int) {
	// const N = ^uint64(0);
	const N = 1000000000
	const proc = 4
	const base = N/proc
	var wg sync.WaitGroup
	now := time.Now()
	cal := func(from, to uint64) {
		defer wg.Done()
		for x := from; x < to; x++ {
			f(x)
		}
		fmt.Printf("avr %s per op\n", time.Since(now)/N)
	}

	for i := uint64(0); i< proc; i++ {
		wg.Add(1)
		go cal(i*base, (i+1)*base)
	}
	wg.Wait()

}
//通过一个间接的函数调用，不利于对比
func BenchmarkClearing1(b *testing.B) {
	// benchmark(b, PopCountByClearing, 0x1234567890ABCDEF)
	for i := 0; i < b.N; i++ {
		PopCountByClearing(0x0)
	}
}
func BenchmarkPop1(b *testing.B) {
	// benchmark(b, PopCount, 0x1234567890ABCDEF)
	for i := 0; i < b.N; i++ {
		PopCount(0x0)
	}
}


func BenchmarkPopAvr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkAvr(PopCount)
	}
}
func BenchmarkClearingAvr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkAvr(PopCountByClearing)
	}
}
// Go 1.6, 2.67GHz Xeon
// $ go test -cpu=4 -bench=. gopl.io/ch2/popcount
// BenchmarkPopCount-4                  200000000         6.30 ns/op
// BenchmarkBitCount-4                  300000000         4.15 ns/op
// BenchmarkPopCountByClearing-4        30000000         45.2 ns/op
// BenchmarkPopCountByShifting-4        10000000        153 ns/op
//
// Go 1.6, 2.5GHz Intel Core i5
// $ go test -cpu=4 -bench=. gopl.io/ch2/popcount
// BenchmarkPopCount-4                  200000000         7.52 ns/op
// BenchmarkBitCount-4                  500000000         3.36 ns/op
// BenchmarkPopCountByClearing-4        50000000         34.3 ns/op
// BenchmarkPopCountByShifting-4        20000000        108 ns/op
//
// Go 1.7, 3.5GHz Xeon
// $ go test -cpu=4 -bench=. gopl.io/ch2/popcount
// BenchmarkPopCount-12                 2000000000        0.28 ns/op
// BenchmarkBitCount-12                 2000000000        0.27 ns/op
// BenchmarkPopCountByClearing-12       100000000        18.5 ns/op
// BenchmarkPopCountByShifting-12       20000000         70.1 ns/op
