package main

import(
	"fmt"
	"os"
	"net"
	"strings"
	"time"
	"log"
)

func main()  {
	// if
	clocks := make(map[string]string)
	for _, str := range os.Args[1:] {
		if i := strings.Index(str, "="); i == -1 {
			fmt.Fprintf(os.Stderr, "invalid parameter %s\n", str)
			continue
		} else {
			clocks[str[:i]] = str[i+1:]
		}
	}

	for name, clock := range clocks {
		fmt.Printf("%s\t%s\n", name, getTime(clock))
	}
	for name, clock := range clocks {
		go printTime(name, clock)
	}
	for name, clock := range clocks {
		go printTime(name, clock)
	}
	time.Sleep(10 * time.Second) //wait for goroutines
}

func getTime(server string) string {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		log.Print(err)
		return "error"
	}
	defer conn.Close()
	var s string
	fmt.Fscanf(conn, "%s", &s)
	return s
}

func printTime(name, server string) {
	fmt.Printf("%s\t%s\n", name, getTime(server))
	// conn, err := net.Dial("tcp", server)
	// if err != nil {
	// 	fmt.Print(err)
	// 	return
	// }
	// defer conn.Close()
	// var t string
	// fmt.Fscanf(conn, "%s", &t)
	// fmt.Printf("%s\t%s\n", name, t)
}
