package exercise5

import (
	"fmt"
	"net/http"
	"time"
	"log"
)

func WaitForServer (url string, timeout time.Duration) error {
	// fmt.Println("WaitForServer begin", time.Now())
	deadline := time.Now().Add(timeout)

	for retries := 0; time.Now().Before(deadline); retries++ {
		// fmt.Println("http.Head() begin", time.Now())
		_, err := http.Head(url)
		// fmt.Println("http.Head() return", time.Now())
		if err == nil {
			return nil
		}
		log.Printf("server not reponding(%v). retry...", err)
		// fmt.Println("sleep...", time.Now())
		time.Sleep(time.Second << uint(retries))
		// fmt.Println("wake up", time.Now())
	}

	// fmt.Println("time out", time.Now())
	return fmt.Errorf("server %s not reponding in %s", url, timeout)
}
