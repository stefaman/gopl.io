//stefaman 20170526

//e11.2 p307
//e11.7 p323
package intset_test

import(
	"testing"
	"math/rand"
	"time"
	// "gopl.io/ch6/intset"
	// "gopl.io/ch6/intsetUnit"
	"gopl.io/ch11/intset"
)


func TestAddAndHas(t *testing.T) {
	var set intset.IntSet
	for i := 0; i < 1e4; i++{
		if set.Has(i) {
			t.Errorf("%v has %d", set, i)
		}
		set.Add(i)
		if !set.Has(i) {
			t.Errorf("%v has not %d", set, i)
		}
	}
}

func TestAddRand(t *testing.T) {
	var set intset.IntSet
	seed := time.Now().UTC().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	for i := 0; i < 1e4; i++{
		// x := rng.Int()//bug, 太大的整数将会失败，无法分配内存
		x := int(rng.Int31n(1e4))
		set.Add(x)
		if !set.Has(x) {
			t.Errorf("for seed %d, index %d has not %d", seed, i, x)
		}
	}
}

var tests = []struct{
	vals []int
	str  string
}{
	{[]int{1,2,3}, "{1 2 3}"},
}
func TestString(t *testing.T) {
	for _, test := range tests {
		var set intset.IntSet
		for _, v := range test.vals {
			set.Add(v)
		}
		if set.String() != test.str {
			t.Errorf("get %s, want %s", set.String(), test.str)
		}
	}
}

type Set interface{
	Add(int);
	Has(int) bool;
}
type MapSet map[int]bool
func NewSet() MapSet{
	return make(map[int]bool)
}
func (m MapSet) Add(x int) {
	m[x] = true
}
func (m MapSet) Has(x int) bool {
	return m[x]
}

//模拟大整数下IntSet的性能和内存消耗
const base = 1e9

//随机数生成器由外部提供，避免每次调用形成新的生成器
func addS(b *testing.B, set Set, num int) {
	for i := 0; i < b.N; i++ {
		for i := 0; i < num; i++{
			set.Add(i+base)
		}
	}
}
func addRand(b *testing.B, set Set, num int32) {
	seed := time.Now().UTC().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	s := make([]int, num)
	for i := range s {
		s[i] = base + int(rng.Int31n(num))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, x := range s {
			set.Add(x)
		}
	}
}
func addX(b *testing.B, set Set) {
	seed := time.Now().UTC().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	x := base+int(rng.Int31n(1e6))//不要放在循环内
	for i := 0; i < b.N; i++ {
		set.Add(x)
		if !set.Has(x) {
			b.Fatalf("for seed %d, set %v has not %d", seed, set, x)
		}
	}
}

func BenchmarkMapAddX(b *testing.B) {
	set := NewSet()
	addX(b, set)
}

func BenchmarkIntSetAddX(b *testing.B) {
	var set intset.IntSet
	addX(b, &set)
}
func BenchmarkMapAddRand1e4(b *testing.B) {
	s := NewSet()
	addRand(b, s, 1e4)
}
func BenchmarkMapAddRand1e5(b *testing.B) {
	s := NewSet()
	addRand(b, s, 1e5)
}
func BenchmarkMapAddRand1e6(b *testing.B) {
	s := NewSet()
	addRand(b, s, 1e6)
}

func BenchmarkIntSetAddRand1e4(b *testing.B) {
	var s intset.IntSet
	addRand(b, &s, 1e4)
}
func BenchmarkIntSetAddRand1e5(b *testing.B) {
	var s intset.IntSet
	addRand(b, &s, 1e5)
}
func BenchmarkIntSetAddRand1e6(b *testing.B) {
	var s intset.IntSet
	addRand(b, &s, 1e6)
}

func BenchmarkMapAddS1e4(b *testing.B) {
	s := NewSet()
	addS(b, s, 1e4)
}
func BenchmarkMapAddS1e5(b *testing.B) {
	s := NewSet()
	addS(b, s, 1e5)
}
func BenchmarkMapAddS1e6(b *testing.B) {
	s := NewSet()
	addS(b, s, 1e6)
}

func BenchmarkIntSetAddS1e4(b *testing.B) {
	var s intset.IntSet
	addS(b, &s, 1e4)
}
func BenchmarkIntSetAddS1e5(b *testing.B) {
	var s intset.IntSet
	addS(b, &s, 1e5)
}
func BenchmarkIntSetAddS1e6(b *testing.B) {
	var s intset.IntSet
	addS(b, &s, 1e6)
}
