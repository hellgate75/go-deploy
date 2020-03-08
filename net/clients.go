package net

import (
	"github.com/hellgate75/go-deploy/net/generic"
	"github.com/hellgate75/go-deploy/net/gotcp"
	"github.com/hellgate75/go-deploy/net/ssh"
)

type NewCoonectionHandlerFunc func(bool)(generic.ConnectionHandler)

func NewSshConnectionHandler(singleSession bool) generic.ConnectionHandler {
	return ssh.NewSshConnectionHandler(singleSession)
}

func NewGoTCPConnectionHandler(singleSession bool) generic.ConnectionHandler {
	return gotcp.NewGoTCPConnectionHandler(singleSession)
}
