package main

import (
	"io"
	"net"
	"bufio"
	log "github.com/Sirupsen/logrus"
	"strings"
)

/*
proxy Starts a connection to the proxy server as well as a connection to localserver
It performs auth on the proxy server
 */
func proxy() {
	serverConn := connect("localhost:8088")
	buf := bufio.NewReader(serverConn)
	//TODO link credentials received from rpc call
	proxyAuth(buf, serverConn, "dummy")

	localConn, _ := net.Dial("tcp", "localhost:22")
	go copyConn(localConn, serverConn)
	copyConn(serverConn, localConn)
}

func waitForString(buf *bufio.Reader) string {
	line, err := buf.ReadString('\n')
	if err != nil {
		log.WithFields(log.Fields{"Error": err, "Line": line}).Fatal("AgentConnection unexpected end")
	}
	return strings.TrimSuffix(line, "\n")
}

func copyConn(writer net.Conn, reader net.Conn) {
	_, err:= io.Copy(writer, reader)
	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Error("io Copy error")
	}
}

func connect(address string) net.Conn {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}
