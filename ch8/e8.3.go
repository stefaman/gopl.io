package main

import(
	"log"
	"net"
	"io"
	"os"
)

func main()  {
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{net.ParseIP("127.0.0.1"), 8000, ""})
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func(){
		if _, err:= io.Copy(os.Stdout, conn); err != nil {
			log.Print(err)
		}
		conn.Close()// server closed first
		log.Print("done")
		done <- struct{}{}
	}()
	if _, err := io.Copy(conn, os.Stdin); err != nil {
		log.Print(err)
	}
	//如果不再`<-done`之前关闭write half of connection，当client os.Stdin EOF后，由于io.Copy(conn, os.Stdin)退出，但是server 没有读到EOF，没有关闭链接，goroutine 中io.Copy(os.Stdout, conn)将阻塞，导致程序僵尸
	// conn.Close() // Closing the read half causes the backg round goroutine’s call to io.Copy to retur n a ‘‘read from closed connection’’ error
	//client closed write first
	conn.CloseWrite()
	// used (*net.TCPConn), Closing the write half of the connection causes the server to see an end-of-file condition
	<-done

}
