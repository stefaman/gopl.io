package main

import(
	"log"
	"flag"
	"time"
	"net"
	"io"
)

func main()  {
	ip := flag.String("ip", "localhost", "ip address to listen")
	port := flag.String("port", "13", "port number to listen")
	flag.Parse()

	listener, err := net.Listen("tcp", *ip + ":" + *port)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(c net.Conn )  {
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		if err != nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
}
