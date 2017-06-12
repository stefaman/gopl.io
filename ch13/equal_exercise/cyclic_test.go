// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

package equal

import (
	"bytes"
	"fmt"
	"testing"
)

func TestCyclic(t *testing.T) {
	one := 1

	type CyclePtr *CyclePtr
	var cyclePtr CyclePtr
	cyclePtr = &cyclePtr

	type CycleSlice []CycleSlice
	var cycleSlice = make(CycleSlice, 1)
	cycleSlice[0] = cycleSlice

	type CycleMap map[string]CycleMap
	var cycleMap = CycleMap{}
	cycleMap[""] = cycleMap

	ch1 := make(chan int)

	type mystring string

	var iface1 interface{} = &one

	for _, test := range []struct {
		x    interface{}
		want bool
	}{
		// basic types
		{1, false},   // different values
		{1.0, false}, // different types
		{"foo", false},
		{mystring("foo"), false}, // different types
		// slices
		{[]string{"foo"}, false},
		{[]string(nil), false},
		{[]string{}, false},
		// slice cycles
		{cycleSlice, true},
		// maps
		{
			map[string][]int{"foo": {1, 2, 3}},
			false,
		},
		{
			map[string][]CyclePtr{"key": {cyclePtr}},
			true,
		},
		{
			map[string]CycleSlice{"key": cycleSlice},
			true,
		},
		// {//stack out
		// 	cycleMap,
		// 	true,
		// },
		// pointers
		{&one, false},
		{new(bytes.Buffer), false},
		// pointer cycles
		{cyclePtr, true},
		// functions
		{func() {}, false},
		// arrays
		{[...]int{1, 2, 3}, false},
		// channels
		{ch1, false},
		// interfaces
		{&iface1, false},
	} {
		if IsCyclic(test.x) != test.want {
			t.Errorf("IsCyclic(%v) = %t",
				test.x, !test.want)
		}
	}
}

func Example_IsCyclic() {
	//!+cycle
	// Circular linked lists a -> b -> a and c -> c.
	type link struct {
		value string
		tail  *link
	}
	a, b, c := &link{value: "a"}, &link{value: "b"}, &link{value: "c"}
	a.tail, b.tail, c.tail = b, a, c
	fmt.Println(IsCyclic(a)) // "true"
	fmt.Println(IsCyclic(b)) // "true"
	fmt.Println(IsCyclic(c)) // "true"
	//!-cycle

	// Output:
	// true
	// true
	// true

}
