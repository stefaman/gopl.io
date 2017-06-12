package main

import(
	// "fmt"
	"sync"
	"testing"
)


var byteBits [256]byte
var bs [256]byte

func init(){
	for i := 0; i < 256; i++ {
		bs[i] = byteBits[i/2] + byte(i&1)
	}
}

func initialize() {
	for i := 0; i < 256; i++ {
		byteBits[i] = byteBits[i/2] + byte(i&1)
	}
}

var initOnce sync.Once
func PopCountLazy(n uint64) int {
	initOnce.Do(initialize)
	return int(
		byteBits[ n >> 56] +
		byteBits[ n << 8 >> 56] +
		byteBits[ n << 16 >> 56] +
		byteBits[ n << 24 >> 56] +
		byteBits[ n << 32 >> 56] +
		byteBits[ n << 40 >> 56] +
		byteBits[ n << 48 >> 56] +
		byteBits[ n << 56 >> 56])
}

func PopCountLazyV2(n uint64) int {
	initOnce.Do(initialize)
	return int(
		byteBits[ byte(n)] +
		byteBits[ byte(n >> 8) ] +
		byteBits[ byte(n >> 16) ] +
		byteBits[ byte(n >> 24) ] +
		byteBits[ byte(n >> 32) ] +
		byteBits[ byte(n >> 40) ] +
		byteBits[ byte(n >> 48) ] +
		byteBits[ byte(n >> 56) ])
}

func PopCount(n uint64) int {
	return int(
		bs[ n >> 56] +
		bs[ n << 8 >> 56] +
		bs[ n << 16 >> 56] +
		bs[ n << 24 >> 56] +
		bs[ n << 32 >> 56] +
		bs[ n << 40 >> 56] +
		bs[ n << 48 >> 56] +
		bs[ n << 56 >> 56])
}

func BenchmarkPopCountLazy(b *testing.B){
	for i := 0; i < b.N; i++ {
		PopCountLazy(255)
	}
}
func BenchmarkPopCountLazyV2(b *testing.B){
	for i := 0; i < b.N; i++ {
		PopCountLazyV2(255)
	}
}
func BenchmarkPopCount(b *testing.B){
	for i := 0; i < b.N; i++ {
		PopCount(255)
	}
}
