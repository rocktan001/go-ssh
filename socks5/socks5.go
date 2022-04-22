package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/rocktan001/tcpkeepalive"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
	"reflect"
	"strconv"
	"time"
	"unsafe"
)

var channel_exit = make(chan int)

func main() {
	for {

		sshclient, _ := SSHConnect("43.155.80.104", 62001)
		if sshclient != nil {
			go SSHLooper(sshclient)
			SSHWait(sshclient)
		}
		time.Sleep(time.Duration(10) * time.Second)
	}

}

func SSHConnect(host string, port int) (sshclient *ssh.Client, err error) {
	pKey, err := ssh.ParsePrivateKey(publickeybytes)
	if err != nil {
		log.Panic(err)
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
	ip := host + ":" + strconv.Itoa(port)
	sshclient, err = ssh.Dial("tcp", ip, &config)
	if err != nil {
		log.Println("ssh dial error=============>", err)
		return nil, err
	}
	log.Println("连接" + ip + "成功")
	v0 := reflect.ValueOf(sshclient). //ptr *ssh.Client
						Elem().   //struct ssh.Client
						Field(0). //interface ssh.Conn
						Elem().   //ptr *ssh.connection
						Elem().   //struct ssh.connection
						Field(1). //struct ssh.sshConn
						Field(0). //interface net.Conn
						Elem()    //ptr *net.TCPConn
	conn := (*net.TCPConn)(unsafe.Pointer(v0.Pointer()))
	tcpkeepalive.SetKeepAlive(conn, 15*time.Second, 3, 15*time.Second, 30*time.Second)
	return sshclient, nil
}

func SSHWait(sshclient *ssh.Client) {
	log.Println("SSHWait...start")
	sshclient.Wait()
	log.Println("SSHWait...over")
	channel_exit <- 0
}
func SSHLooper(sshclient *ssh.Client) {
	server, err := net.Listen("tcp", ":62002")
	if err != nil {
		log.Println("listen 62000 err ", err)
		return
	}

	// 当sshclient 断开后，要主动关闭lisener 。结束循环
	go func(server net.Listener) {
		<-channel_exit
		log.Println("SSHLooper ...OVER")
		server.Close()
	}(server)

	for {
		client, err := server.Accept()
		if err != nil {
			log.Println("Accept failed: %v", err)
			return
		}
		go process(client, sshclient)
	}

}
func process(client net.Conn, sshclient *ssh.Client) {
	if err := Socks5Auth(client); err != nil {
		fmt.Println("auth error:", err)
		client.Close()
		return
	}

	target, err := Socks5Connect(client, sshclient)
	if err != nil {
		fmt.Println("connect error:", err)
		client.Close()
		return
	}

	_conn, ok := client.(*net.TCPConn)
	if ok {
		tcpkeepalive.SetKeepAlive(_conn, 15*time.Second, 3, 15*time.Second, 30*time.Second)
	}
	Socks5Forward(client, target)
}

func Socks5Auth(client net.Conn) (err error) {
	buf := make([]byte, 256)

	// 读取 VER 和 NMETHODS
	n, err := io.ReadFull(client, buf[:2])
	if n != 2 {
		return errors.New("reading header: " + err.Error())
	}

	ver, nMethods := int(buf[0]), int(buf[1])
	if ver != 5 {
		return errors.New("invalid version")
	}

	// 读取 METHODS 列表
	n, err = io.ReadFull(client, buf[:nMethods])
	if n != nMethods {
		return errors.New("reading methods: " + err.Error())
	}

	//无需认证
	n, err = client.Write([]byte{0x05, 0x00})
	if n != 2 || err != nil {
		return errors.New("write rsp: " + err.Error())
	}

	return nil
}

func Socks5Connect(client net.Conn, sshclient *ssh.Client) (net.Conn, error) {
	buf := make([]byte, 256)

	n, err := io.ReadFull(client, buf[:4])
	if n != 4 {
		return nil, errors.New("read header: " + err.Error())
	}

	ver, cmd, _, atyp := buf[0], buf[1], buf[2], buf[3]
	if ver != 5 || cmd != 1 {
		return nil, errors.New("invalid ver/cmd")
	}

	addr := ""
	switch atyp {
	case 1:
		n, err = io.ReadFull(client, buf[:4])
		if n != 4 {
			return nil, errors.New("invalid IPv4: " + err.Error())
		}
		addr = fmt.Sprintf("%d.%d.%d.%d", buf[0], buf[1], buf[2], buf[3])

	case 3:
		n, err = io.ReadFull(client, buf[:1])
		if n != 1 {
			return nil, errors.New("invalid hostname: " + err.Error())
		}
		addrLen := int(buf[0])

		n, err = io.ReadFull(client, buf[:addrLen])
		if n != addrLen {
			return nil, errors.New("invalid hostname: " + err.Error())
		}
		addr = string(buf[:addrLen])

	case 4:
		return nil, errors.New("IPv6: no supported yet")

	default:
		return nil, errors.New("invalid atyp")
	}

	n, err = io.ReadFull(client, buf[:2])
	if n != 2 {
		return nil, errors.New("read port: " + err.Error())
	}
	port := binary.BigEndian.Uint16(buf[:2])

	destAddrPort := fmt.Sprintf("%s:%d", addr, port)
	dest, err := sshclient.Dial("tcp", destAddrPort)
	if err != nil {
		sshclient.Close()
		return nil, errors.New("dial dst: " + err.Error())
	}

	n, err = client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	if err != nil {
		dest.Close()
		return nil, errors.New("write rsp: " + err.Error())
	}

	return dest, nil
}

func Socks5Forward(client, target net.Conn) {
	channel_err := make(chan int)
	go func(dst, src net.Conn, c chan int) {
		// log.Printf("[%s %s] <========>[%s %s]\n", server.LocalAddr(), server.RemoteAddr(), src.LocalAddr(), src.RemoteAddr())
		_, err := io.Copy(dst, src)
		if err != nil {
			// log.Printf("io.Copy error: %s\n", err)
			c <- 1
		}
		c <- 4
	}(target, client, channel_err)

	go func(dst, src net.Conn, c chan int) {

		_, err := io.Copy(dst, src)
		if err != nil {
			// log.Printf("io.Copy error: %s\n", err)
			c <- 2
		}
		c <- 3
	}(client, target, channel_err)

	val := <-channel_err

	defer func(val int) {
		// log.Println("return =========val ", val)
		// log.Printf("[%s %s] <========>[%s %s]\n", server.LocalAddr(), server.RemoteAddr(), conn.LocalAddr(), conn.RemoteAddr())
		client.Close()
		target.Close()
	}(val)
}

var publickeybytes = []byte(
	`-----BEGIN OPENSSH PRIVATE KEY-----
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
-----END OPENSSH PRIVATE KEY-----`)
