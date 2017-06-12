package main

import(
	"log"
	"net"
	"io"
	"os"
)

func main()  {
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{net.ParseIP("127.0.0.1"), 8080, ""})
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func(){
		//when server closed write half of connection, io.Copy got EOF and return
		if _, err:= io.Copy(os.Stdout, conn); err != nil {
			log.Print(err)
		}
		log.Print("server disconnet\n")
		done <- struct{}{}
	}()
	 go func(){
		//when os.Stdin get EOF, io.Copy return, but server will not read EOF
		if _, err := io.Copy(conn, os.Stdin); err != nil {
			log.Print(err)
		}
		//closed write half of connection，server will read EOF
		conn.CloseWrite()
		log.Print("stop entring\n")
	}()

	<-done
}
