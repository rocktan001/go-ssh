package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

var host = flag.String("host", "www.rocktan001.com", "host")
var port = flag.String("port", "62002", "port")

// var host = flag.String("host", "localhost", "host")
// var port = flag.String("port", "1197", "port")

func main() {
	flag.Parse()
	conn, err := net.Dial("tcp", *host+":"+*port)
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("Connecting to " + *host + ":" + *port)
	done := make(chan string)
	// go handleWrite(conn, done)
	go handleLoopWrite(conn, done)
	// go handleRead(conn, done)
	go handleLoopRead(conn, done)
	fmt.Println(<-done)
	fmt.Println(<-done)
}
func handleWrite(conn net.Conn, done chan string) {
	for i := 10; i > 0; i-- {
		_, e := conn.Write([]byte("hello " + strconv.Itoa(i) + "\r\n"))
		if e != nil {
			fmt.Println("Error to send message because of ", e.Error())
			break
		}
	}
	done <- "Sent"
}
func handleRead(conn net.Conn, done chan string) {
	buf := make([]byte, 1024)
	reqLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error to read message because of ", err)
		return
	}
	fmt.Println(string(buf[:reqLen-1]))
	done <- "Read"
}
func handleLoopRead(conn net.Conn, done chan string) {
	buf := make([]byte, 1024)
	for {
		reqLen, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error to read message because of ", err)
			break
		}
		fmt.Println(string(buf[:reqLen-1]))
	}
	done <- "Read"
}
func handleLoopWrite(conn net.Conn, done chan string) {
	for {

		_, e := conn.Write([]byte("hello " + time.Now().Format("2006-01-02 15:04:05") + "\r\n"))
		if e != nil {
			fmt.Println("Error to send message because of ", e.Error())
			break
		}
		time.Sleep(3000 * time.Second)
	}
	done <- "Sent"
}
