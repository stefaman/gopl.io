package main

import(
	"fmt"
	"os"
	"strconv"
)


func main()  {
	var nums []int
	for _, numStr := range os.Args[1:] {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parsing input: %v\n", err)
			continue
		}
		nums = append(nums, num) //ignore error
	}
	fmt.Printf("max num is %d\n", max(nums...))

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "use main <integer> ...\n")
		os.Exit(1)
	}
	nums = nil
	var first int
	for index, numStr := range os.Args[1:] {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parsing input: %v\n", err)
			continue
		}
		if index == 0 {
			first = num
			continue
		}
		nums = append(nums, num) //ignore error
	}
	fmt.Printf("max num is %d\n", maxV2(first, nums...))
}

func max(nums ...int) int  {
	if len(nums) == 0 {
		return 0
	}
	biggest := nums[0]
	for _, num := range nums {
		if biggest < num {
			biggest = num
		}
	}
	return biggest

}

func maxV2(first int, nums ...int) int  {
	biggest := first
	for _, num := range nums {
		if biggest < num {
			biggest = num
		}
	}
	return biggest

}
