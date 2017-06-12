// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 224.

// Reverb2 is a TCP server that simulates an echo.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
	"sync"
)

func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

//!+
func handleConn(c net.Conn) {
	defer c.Close()
	input := bufio.NewScanner(c)
	read := make(chan bool)
	go func(){
		read <-input.Scan()
	}()
	var wg sync.WaitGroup
	exit:
	for  {
		select{
			case ok := <- read:
				if !ok {
					break exit
				}
			case <-time.After(10*time.Second):
				fmt.Fprintln(c, "\ttime out!")
				break exit
			//return //bug, the left "echo" will lost
		}
		if err := input.Err(); err != nil {
			log.Print(err)
			continue
		}
		wg.Add(1)
		go func(){
			defer wg.Done()
			echo(c, input.Text(), 6*time.Second)
		}()
		go func(){ read <-input.Scan()}()
	}

	wg.Wait()
}

//!-

func main() {
	l, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn)
	}
}
