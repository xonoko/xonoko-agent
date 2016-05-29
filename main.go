package main

import (
	"golang.org/x/net/websocket"
	"net/rpc"
	logrus "github.com/Sirupsen/logrus"
)

type Args struct {
	A int
	B int
}

func main() {
	rpc.Register(new(Agent))
	origin := "http://localhost/"
	url := "ws://carambar.lucas-galton.fr:8200/conn"
	ws, err := websocket.Dial(url, "", origin)

	if err != nil {
		logrus.WithFields(logrus.Fields{"Error": err}).Error("Cannot connect to websocket server")
	} else {
		rpc.ServeConn(ws)
	}
	//proxy()
}

