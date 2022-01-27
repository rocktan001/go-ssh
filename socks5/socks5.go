package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/rocktan001/tcpkeepalive"
	"golang.org/x/crypto/ssh"
	"io"
	// "io/ioutil"
	"log"
	"net"
	"os"
	"reflect"
	"time"
	"unsafe"
)

var publickeybytes = []byte(
	`-----BEGIN OPENSSH PRIVATE KEY-----
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
-----END OPENSSH PRIVATE KEY-----`)

func socks5Proxy(conn net.Conn) {
	// defer conn.Close()

	var b [1024]byte

	n, err := conn.Read(b[:])
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("% x", b[:n])

	conn.Write([]byte{0x05, 0x00})

	n, err = conn.Read(b[:])
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("% x", b[:n])

	var addr string
	switch b[3] {
	case 0x01:
		sip := sockIP{}
		if err := binary.Read(bytes.NewReader(b[4:n]), binary.BigEndian, &sip); err != nil {
			log.Println("请求解析错误")
			return
		}
		addr = sip.toAddr()
	case 0x03:
		host := string(b[5 : n-2])
		var port uint16
		err = binary.Read(bytes.NewReader(b[n-2:n]), binary.BigEndian, &port)
		if err != nil {
			log.Println(err)
			return
		}
		addr = fmt.Sprintf("%s:%d", host, port)
	}
	// log.Println("addr: ", addr)
	server, err := client.Dial("tcp", addr)
	if err != nil {
		log.Println(err)
		return
	}
	conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	// go io.Copy(server, conn)
	// io.Copy(conn, server)
	//TCP keepalive
	_conn, ok := conn.(*net.TCPConn)
	if ok {
		// log.Println("tcpkeepalive")
		tcpkeepalive.SetKeepAlive(_conn, 15*time.Second, 3, 15*time.Second, 30*time.Second)
	}
	channel_err := make(chan int)
	go func(dst, src net.Conn, c chan int) {
		// log.Printf("[%s %s] <========>[%s %s]\n", server.LocalAddr(), server.RemoteAddr(), src.LocalAddr(), src.RemoteAddr())
		_, err := io.Copy(dst, src)
		if err != nil {
			log.Printf("io.Copy error: %s\n", err)
			c <- 1
		}
		c <- 4
	}(server, conn, channel_err)

	go func(dst, src net.Conn, c chan int) {

		_, err := io.Copy(dst, src)
		if err != nil {
			log.Printf("io.Copy error: %s\n", err)
			c <- 2
		}
		c <- 3
	}(conn, server, channel_err)

	val := <-channel_err

	defer func(val int) {
		log.Println("return =========val ", val)
		conn.Close()
		server.Close()
	}(val)
}

type sockIP struct {
	A, B, C, D byte
	PORT       uint16
}

func (ip sockIP) toAddr() string {
	return fmt.Sprintf("%d.%d.%d.%d:%d", ip.A, ip.B, ip.C, ip.D, ip.PORT)
}

func socks5ProxyStart() {
	log.SetFlags(log.Ltime | log.Lshortfile)

	server, err := net.Listen("tcp", ":62002")
	if err != nil {
		log.Panic(err)
	}
	defer server.Close()
	log.Println("开始接受连接")
	for {
		client, err := server.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("一个新连接")
		go socks5Proxy(client)
	}
}

var client *ssh.Client

func main() {
	log.SetOutput(os.Stdout)
	pKey, err := ssh.ParsePrivateKey(publickeybytes)
	if err != nil {
		log.Println(err)
		return
	}
	config := ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(pKey),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	client, err = ssh.Dial("tcp", "43.155.80.104:62001", &config)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("连接服务器成功")
	v0 := reflect.ValueOf(client). //ptr *ssh.Client
					Elem().   //struct ssh.Client
					Field(0). //interface ssh.Conn
					Elem().   //ptr *ssh.connection
					Elem().   //struct ssh.connection
					Field(1). //struct ssh.sshConn
					Field(0). //interface net.Conn
					Elem()    //ptr *net.TCPConn
	conn := (*net.TCPConn)(unsafe.Pointer(v0.Pointer()))
	fmt.Println(conn.RemoteAddr(), "<==>", conn.LocalAddr())
	// _conn, ok := conn.(*net.TCPConn)
	// if ok {
	tcpkeepalive.SetKeepAlive(conn, 15*time.Second, 3, 15*time.Second, 30*time.Second)
	defer client.Close()
	// client.Dial()
	go func() {
		client.Wait()
		log.Panic("return....")
	}()
	socks5ProxyStart()
	return
}
