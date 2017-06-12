package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	//exercise 1.2
	// for index, arg := range os.Args[1:] {
	// 	fmt.Println(index, arg)
	// }

	//exercise 1.3
	// s, seq := "", ""
	// for _, arg := range os.Args[1:] {
	// 	s += seq + arg
	// 	seq = " "
	// }
	// fmt.Println(s)

	//use string.Join
	fmt.Println(strings.Join(os.Args[1:], " "))

	for index := 0; index < 1000000;  {
		fmt.Println("just for test")
	}
}
