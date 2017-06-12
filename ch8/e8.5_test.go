package test

import(
	// "log"
	"testing"
	"math"
	"math/rand"
	"time"
)

func task(p int, f func())  {
	N := 10000
	done := make(chan struct{}, p)
	// done := make(chan struct{})
	for i := 0; i < p; i++ {
		go func(i int){
			// log.Printf("%d goroutine begin\n", i)
			for i := 0; i < N/p; i++ {
				f()
			}
			done <- struct{}{}
			// log.Printf("%d goroutine over\n", i)
		}(i)
	}
	for i := 0; i < p; i++ {
		<- done
	}
}

func CPUTask(p int)  {
	task(p, cpuHeavy)
}
func IOTask(p int)  {
	task(p, ioHeavy)
}
func RandNormTask(p int)  {
	task(p, randNorm)
}
func RandNewTask(p int)  {
	task(p, randNew)
}
func cpuHeavy() {
	// 下面的执行执行时间数并发数增加，待查？？
	// for i := 0; i < 1000; i++ {
	// 	math.Sinh(rand.NormFloat64())
	// 	math.Log2(rand.NormFloat64())
	// }

	for i := 0; i < 10000; i++ {
		for j := 0; i < 10000; i++ {
			math.Sinh(float64(i*j)+float64(i+j))
		}
	}
}

func ioHeavy(){
	math.Log10(rand.NormFloat64())
	delay := 100*time.Nanosecond
	time.Sleep( delay + time.Duration(float64(delay) * 0.4 * rand.NormFloat64()))

}

// rand的default source的多协程安全的，内部有sync.Mutex。
//多协程反而恶化性能
func randNorm() {
	for i := 0; i < 1000; i++ {
		math.Sinh(rand.NormFloat64())
	}
}

//每个协程使用独立的源
func randNew() {
	rander := rand.New(rand.NewSource(1))
	for i := 0; i < 1000; i++ {
		math.Sinh(rander.NormFloat64())
	}
}

func BenchmarkRandNorm1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandNormTask(1)
	}
}
func BenchmarkRandNorm4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandNormTask(4)
	}
}
func BenchmarkRandNorm8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandNormTask(8)
	}
}

func BenchmarkRandNew1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandNewTask(1)
	}
}
func BenchmarkRandNew4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandNewTask(4)
	}
}
func BenchmarkRandNew8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandNewTask(8)
	}
}


func BenchmarkIO1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IOTask(1)
	}
}
func BenchmarkIO4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IOTask(4)
	}
}
func BenchmarkIO8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IOTask(8)
	}
}

func BenchmarkCPU1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CPUTask(1)
	}
}

func BenchmarkCPU2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CPUTask(2)
	}
}
func BenchmarkCPU4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CPUTask(4)
	}
}
// func BenchmarkCPU6(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		CPUTask(6)
// 	}
// }
func BenchmarkCPU8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CPUTask(8)
	}
}

func BenchmarkCPU16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CPUTask(16)
	}
}
