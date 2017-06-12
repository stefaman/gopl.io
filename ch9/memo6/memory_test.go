package memory_test

import(
	"testing"
 "gopl.io/ch9/memo6"
	// "gopl.io/ch9/memotest"
)

func TestSequential(t *testing.T) {
	memo := memory.New(HTTPGetBody)
	Sequential(t, memo)
}

func TestConcurrent(t *testing.T) {
	memo := memory.New(HTTPGetBody)
	Concurrent(t, memo)
}
