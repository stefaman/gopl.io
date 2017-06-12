// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

package intset

import "fmt"

func Example_one() {
	//!+main
	var x, y IntSet
	x.Add(1)
	x.Add(144)
	x.Add(9)
	fmt.Println(x.String()) // "{1 9 144}"

	y.Add(9)
	y.Add(42)
	fmt.Println(y.String()) // "{9 42}"

	x.UnionWith(&y)
	fmt.Println(x.String()) // "{1 9 42 144}"

	fmt.Println(x.Has(9), x.Has(123)) // "true false"
	//!-main

	// Output:
	// {1 9 144}
	// {9 42}
	// {1 9 42 144}
	// true false
}

func Example_two() {
	var x IntSet
	x.Add(1)
	x.Add(144)
	x.Add(9)
	x.Add(42)

	//!+no
	fmt.Println(&x)         // "{1 9 42 144}"
	fmt.Println(x.String()) // "{1 9 42 144}"
	fmt.Println(x)          // "{[4398046511618 0 65536]}"
	//!-note

	// Output:
	// {1 9 42 144}
	// {1 9 42 144}
	// {[4398046511618 0 65536]}
}

//stefaman

func Example_three() {
	// var x IntSet
	x := new(IntSet)
	x.AddAll()

	fmt.Println(x.Len()) // "0"
	x.AddAll(0)
	fmt.Println(x.Len()) // "1"
	x.AddAll(1)
	fmt.Println(x.Len()) // "2"
	x.AddAll(2, 3, 4)
	fmt.Println(x) // {0 1 2 3 4}
	x.Remove(3)
	fmt.Println(x) //{0 1 2 4}

	t := x.Copy()
	x.Clear();
	x.AddAll(2,3,4,5) //
	fmt.Println(x) //{2 3 4 5}
	fmt.Println(t) //{0 1 2 4}
	fmt.Println(t.Elems()) //[0 1 2 4]





	// Output:
	// 0
	// 1
	// 2
	// {0 1 2 3 4}
	// {0 1 2 4}
	// {2 3 4 5}
	// {0 1 2 4}
	// [0 1 2 4]

}


func Example_four() {
	// var x IntSet
	a := make([]int, 160)

	r := new(IntSet)
	s := new(IntSet)
	t := new(IntSet)
	x := new(IntSet)


	for i := 0; i < len(a); i++ {
		a[i] = i
	}
	x.AddAll(a[8:64]...)
	x.AddAll(a[128:]...)

	t.AddAll(a[:80]...)
	r.AddAll(a[:]...)






	s = x.Copy()
	s.DifferenceWith(t);
	fmt.Println(s) //{3 66 67}

	s = x.Copy()
	s.IntersectWith(t);
	fmt.Println(s) //

	s = x.Copy()
	s.SymmetricDifference(t);
	fmt.Println(s) //{1 2 3 66 67}


	// Output:
	// {66 67}
	// {1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31 32 33 34 35 36 37 38 39 40 41 42 43 44 45 46 47 48 49 50 51 52 53 54 55 56 57 58 59 60 61 62 63 64 65}
	// {0 1 66 67}

}
