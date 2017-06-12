package main

import(
	"fmt"
	"time"
)

func main()  {
	ping := make(chan int)
	pong := make(chan int)
	now := time.Now()

	go func(){
		pong<- 111
		n := 0
		for {
			pong<- <-ping
			n++
			if n == 100000 {
				d := time.Since(now)
				nPerSecond := float32(n)/float32(d) * 1e9
				fmt.Printf("taks %s ping-pong %d pass\n", d, n)
				fmt.Printf("can ping-pong %.0f pass per second\n", nPerSecond)
			}
		}
	}()
	go func(){
		for {
			ping<- <-pong
		}
	}()

time.Sleep(10*time.Second)


}
