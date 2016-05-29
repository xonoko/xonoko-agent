package main

import (
	"github.com/Sirupsen/logrus"
)

type Bool struct {
	Value bool
}

type Agent int


func (t *Agent) Multiply(args *Args, reply *Args) error {
	logrus.Info("RPC multiply call")
	reply.A = args.A * args.B * 2
	return nil
}
