package main

import(
	"fmt"
	"os"
	"time"
	"strconv"
)

type talk chan int

func connect(from <-chan int, to chan<- int) {
	go func(){
		for x := range from {
			to <- x
		}
		close(to)
	}()
}

func main()  {
	n, _ := strconv.Atoi(os.Args[1])
	init := make(chan int)
	from := init
	var to chan int
	for i:=0;i<n;i++ {
		to = make(chan int)
		connect(from, to)
		from = to
	}

	N, _ := strconv.Atoi(os.Args[2])
	now := time.Now()
	go func(){
		for i:=0;i<N;i++ {
			init <- i
		}
		close(init)
	}()

	<- to
	fmt.Printf("taks %s from pass across %d stages\n", time.Since(now), n)
	// var ret int
	for range to {
	}
	fmt.Printf("taks %s from  pass %d values across %d stages\n", time.Since(now), N, n)


}
