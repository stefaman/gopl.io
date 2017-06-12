package main

import(
	"fmt"
	"sort"
)


func IsPalindrome (data sort.Interface) bool {
	for i, j := 0, data.Len()-1; i < j; i, j = i+1, j-1 {
		if data.Less(i, j) || data.Less(j, i) {
			return false
		}
	}
	return true
}

func main()  {
	ints := sort.IntSlice([]int{1,2, 3, 1, 1})
	fmt.Printf("%v is palindrome? %v\n", ints, IsPalindrome(ints))
}
