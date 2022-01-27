package main

import (
	"flag"
	"fmt"
	"github.com/rocktan001/tcpkeepalive"
	"io"
	"net"
	"os"
	"time"
	// "syscall"
	// "unsafe"
)

var host = flag.String("host", "", "host")
var port = flag.String("port", "2223", "port")

type IdleTimeoutConn struct {
	Conn net.Conn
}

func (self IdleTimeoutConn) Read(buf []byte) (int, error) {
	// self.Conn.SetReadDeadline(time.Now().Add(15 * time.Second))
	// log.Println("==========read=================")
	return self.Conn.Read(buf)
}

func (self IdleTimeoutConn) Write(buf []byte) (int, error) {
	// self.Conn.SetWriteDeadline(time.Now().Add(15 * time.Second))
	nw, ew := self.Conn.Write(buf)
	if ew != nil {
		fmt.Println("Error Write")
	}
	return nw, ew
}
func main() {
	flag.Parse()
	var l net.Listener
	var err error
	l, err = net.Listen("tcp4", *host+":"+*port)
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}

	defer l.Close()
	fmt.Println("Listening on " + *host + ":" + *port)
	for {
		conn, err := l.Accept()
		fmt.Printf("%T %v", conn, conn)
		_conn, ok := conn.(*net.TCPConn)
		if ok {
			fmt.Println("SetKeepAlive")
			// _conn.SetNoDelay(true)
			// c.SetKeepAlive(true)
			// c.SetKeepAlivePeriod(5 * time.Second)
			kaConn, _ := tcpkeepalive.EnableKeepAlive(_conn)
			kaConn.SetKeepAliveIdle(5 * time.Second)
			kaConn.SetKeepAliveCount(3)
			kaConn.SetKeepAliveInterval(5 * time.Second)
			kaConn.SetTcpUserTimeout(15 * time.Second)
		}

		// if err := conn.SetKeepAlive(true); err != nil {
		// 	return err
		// }
		if err != nil {
			fmt.Println("Error accepting: ", err)
			os.Exit(1)
		}
		//logs an incoming message
		fmt.Printf("Received message %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()
	_conn := IdleTimeoutConn{
		Conn: conn,
	}
	chDone := make(chan bool)
	go func() {
		_, err := io.Copy(_conn, _conn)
		// log.Println("copy r-> l ", len)
		if err != nil {
			fmt.Println(fmt.Sprintf("error while copy remote->local: %s", err))

		}
		chDone <- true
	}()
	go func() {
		for {
			_conn.Write([]byte("hello"))
			time.Sleep(3 * time.Second)
		}
	}()
	<-chDone
	fmt.Println("========================handleRequest over==========================")
}
