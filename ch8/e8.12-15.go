package main

import(
	"fmt"
	"bufio"
	"net"
	"log"
	"time"
	"bytes"
)

func main()  {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	go broadcast()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}

}
type client struct {
	name string
	ch chan string
}
var(
	entering = make(chan client)
	leaving = make(chan client)
	message = make(chan string)
)
func broadcast() {
	clients := make(map[client]bool)
	nClients := 0
	// tick := time.NewTicker(100 * time.Millisecond)
	announce := func(msg string) {
		delay := time.Duration(10 + 1000 / nClients) * time.Millisecond
		for cli := range clients {
			select {
			case cli.ch <- msg:
			case <-time.After(delay)://一个协程，代价是否太大？？
			// case <-tick: //most delay is 100ms, but least is 0,用来buffer channel 还是可能丢消息
			}
	}
}
	for {
		select {
		case msg := <- message:
			announce(msg)
		case cli := <- entering:
			if nClients == 0 {
				cli.ch <- "You are first one in chat room"
			} else{
				buf := new(bytes.Buffer)
				fmt.Fprintf(buf, "There had been %d clients in char room：", nClients)
				for c := range clients {
					fmt.Fprintf(buf, "%q ", c.name)
				}
				cli.ch <- buf.String()
			}
			announce(cli.name + " has arrived")//相比客户端处理程序处理，此次可以插队，延迟更小。
			clients[cli] = true
			nClients++
		case cli := <-leaving:
			delete(clients, cli)
			close(cli.ch)
			announce(cli.name + " has arrived")
		}
	}
}
func handleConn(conn net.Conn)  {
	fmt.Fprintf(conn, "%s Please enter your name: ", conn.RemoteAddr())
	var name string
	fmt.Fscanf(conn, "%s\n", &name)
	fmt.Fprintf(conn, "Hello %s. Begin entring\n", name)
	const delay = 300 * time.Second
	const chanBuf = 100
	ch := make(chan string, chanBuf)
	var cli = client{
		name,
		ch,
	}

	// message <-   name + " has arrived" //dealed by chat room broadcast，消息处理同步性好
	//消息被广播后，才处理加入；如果聊天室很繁忙，需要考虑延时；
	entering <- cli
	go clientWriter(cli, conn)

	input := bufio.NewScanner(conn)
	// read := make(chan bool)
	heartBeat := make(chan struct{})
	go func() {
		retry:
		select {
		case <- heartBeat:
			goto retry
		case <- time.After(delay):
		}
		fmt.Fprintln(conn, "no entering in %s, server disconnet", delay)
		conn.Close()
	}()
	
	for input.Scan() {
		heartBeat <- struct{}{}
		if err := input.Err(); err != nil {
			log.Printf("reading %q: %v", name, err)
			fmt.Fprintln(conn, err)
			fmt.Fprintln(conn, "Please enter again")
			continue
		}
		message <- name + ": " + input.Text()
	}

	leaving <- cli
	//已经离开，此消息还在排队；繁忙的聊天室注意延时
	// message <-  cli.name + " has left" //dealed by chat room broadcast，消息处理同步性好
	fmt.Fprintln(conn, "一路顺风")
	conn.Close()
}

func clientWriter(cli client, conn net.Conn) {
	for msg := range cli.ch {
		_, err := fmt.Fprintln(conn, msg)
		if err != nil {
			log.Printf("writing %q: %v", cli.name, err)
		}
	}
}
