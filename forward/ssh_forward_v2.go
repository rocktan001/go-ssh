/*

Go-Language implementation of an SSH Reverse Tunnel, the equivalent of below SSH command:

   ssh -R 8080:127.0.0.1:8080 operatore@146.148.22.123

which opens a tunnel between the two endpoints and permit to exchange information on this direction:

   server:8080 -----> client:8080
``
   once authenticated a process on the SSH server can interact with the service answering to port 8080 of the client
   without any NAT rule via firewall

Copyright 2017, Davide Dal Farra
MIT License, http://www.opensource.org/licenses/mit-license.php

*/

package main

import (
    "fmt"
    "io"
    "os"
    // "io/ioutil"
    "log"
    "net"
    "reflect"
    // "unsafe"
    "github.com/rocktan001/tcpkeepalive"
    "golang.org/x/crypto/ssh"
    "time"
    "unsafe"
)

var publickeybytes = []byte(
    `-----BEGIN OPENSSH PRIVATE KEY-----
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
-----END OPENSSH PRIVATE KEY-----`)

type Endpoint struct {
    Host string
    Port int
}

func (endpoint *Endpoint) String() string {
    return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

// From https://sosedoff.com/2015/05/25/ssh-port-forwarding-with-go.html
// Handle local client connections and tunnel data to the remote server
// Will use io.Copy - http://golang.org/pkg/io/#Copy

type IdleTimeoutConn struct {
    Conn net.Conn
}

func (self IdleTimeoutConn) Read(buf []byte) (int, error) {
    // self.Conn.SetDeadline(time.Now().Add(30 * time.Second))
    // log.Println("==========read=================")
    return self.Conn.Read(buf)
}

func (self IdleTimeoutConn) Write(buf []byte) (int, error) {
    // self.Conn.SetDeadline(time.Now().Add(30 * time.Second))
    // log.Println("==========Write=================")
    return self.Conn.Write(buf)
}

func handleClient(client net.Conn, remote net.Conn) {
    // if c, ok := client.(*net.TCPConn); ok {
    //     log.Println("===============================OK=======================================")
    //     c.SetKeepAlive(true)
    //     c.SetKeepAlivePeriod(5 * time.Second)
    // }
    clientConn := IdleTimeoutConn{
        Conn: client,
    }
    remoteConn := IdleTimeoutConn{
        Conn: remote,
    }
    log.Println("========================haneleClient start=========================")
    log.Printf("%T %v", client, client)
    log.Println(client.LocalAddr(), client.RemoteAddr(), "<========>", remote.LocalAddr(), remote.RemoteAddr())
    log.Println("========================haneleClient start==========================")
    defer client.Close()
    defer remote.Close()

    chDone := make(chan bool)
    // Start remote -> local data transfer
    go func() {
        len, err := io.Copy(clientConn, remoteConn)
        log.Println("copy r-> l ", len)
        if err != nil {
            log.Println(fmt.Sprintf("error while copy remote->local: %s", err))

        }
        chDone <- true
    }()

    // Start local -> remote data transfer
    go func() {
        len, err := io.Copy(remoteConn, clientConn)
        log.Println("copy l-> r ", len)
        if err != nil {
            log.Println(fmt.Sprintf("error while copy local->remote: %s", err))
        }
        chDone <- true
    }()

    <-chDone
    log.Println("========================haneleClient over=========================")
    log.Println(client.LocalAddr(), client.RemoteAddr(), "<========>", remote.LocalAddr(), remote.RemoteAddr())
    log.Println("========================haneleClient over==========================")
}

func publicKeyFile(file string) ssh.AuthMethod {
    // buffer, err := ioutil.ReadFile(file)
    // if err != nil {
    //     log.Fatalln(fmt.Sprintf("Cannot read SSH public key file %s", file))
    //     return nil
    // }

    key, err := ssh.ParsePrivateKey(publickeybytes)
    if err != nil {
        log.Fatalln(fmt.Sprintf("Cannot parse SSH public key file %s", file))
        return nil
    }
    return ssh.PublicKeys(key)
}

// local service to be forwarded
var localEndpoint = Endpoint{
    Host: "192.168.199.233",
    Port: 5900,
}

// remote SSH server
var serverEndpoint = Endpoint{
    Host: "www.rocktan001.com",
    // Host: "rock-001",
    Port: 62001,
}

// remote forwarding port (on remote SSH server network)
var remoteEndpoint = Endpoint{
    Host: "0.0.0.0",
    Port: 62002,
}

var sshConfig = &ssh.ClientConfig{
    // SSH connection username
    User: "root",
    Auth: []ssh.AuthMethod{
        // put here your private key path
        publicKeyFile(".ssh/id_rsa"),
    },
    HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    Timeout:         15 * time.Second, // max time to establish connection
}

func __init() (net.Listener, error) {
    serverConn, err := ssh.Dial("tcp", serverEndpoint.String(), sshConfig)
    if err != nil {
        log.Fatalln(fmt.Printf("Dial INTO remote server error: %s", err))
    }

    v0 := reflect.ValueOf(serverConn). //ptr *ssh.Client
                        Elem().   //struct ssh.Client
                        Field(0). //interface ssh.Conn
                        Elem().   //ptr *ssh.connection
                        Elem().   //struct ssh.connection
                        Field(1). //struct ssh.sshConn
                        Field(0). //interface net.Conn
                        Elem()    //ptr *net.TCPConn

    fmt.Println("==========================================")
    conn := (*net.TCPConn)(unsafe.Pointer(v0.Pointer()))
    fmt.Println(conn.RemoteAddr(), "<==>", conn.LocalAddr())
    tcpkeepalive.SetKeepAlive(conn, 15*time.Second, 3, 15*time.Second, 30*time.Second)
    fmt.Println("==========================================")

    log.Println("=====================connect www.rocktan001.com ok =============================>\n")
    // Listen on remote server port
    listener, err := serverConn.Listen("tcp", remoteEndpoint.String())

    if err != nil {
        log.Fatalln(fmt.Printf("Listen open port ON remote server error: %s", err))
        return nil, err
    }
    return listener, nil
}

func main() {

    // refer to https://godoc.org/golang.org/x/crypto/ssh for other authentication types

    // Connect to SSH remote server using serverEndpoint
    // log.Println(remoteEndpoint.String(), "=============================>\n")
    log.SetOutput(os.Stdout)
    listener, _ := __init()

    defer listener.Close()

    for {

        client, err := listener.Accept()
        if err != nil {
            log.Fatalln("accept error ", err)
        }

        // 把本地连接socket移到远程连接accept 后，必须长时间无连接，本地连接断开情况。
        // Open a (local) connection to localEndpoint whose content will be forwarded so serverEndpoint
        local, err := net.Dial("tcp", localEndpoint.String())
        // handle incoming connections on reverse forwarded tunnel
        // local.SetDeadline(15 * time.Second)
        if err != nil {
            log.Fatalln(fmt.Printf("Dial INTO local service error: %s", err))
        }
        log.Println("serverConn accept--> ", client.LocalAddr(), client.RemoteAddr())
        go handleClient(client, local)
    }

}
