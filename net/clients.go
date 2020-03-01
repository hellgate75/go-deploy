package net

import (
	"github.com/hellgate75/go-deploy/net/generic"
	"github.com/hellgate75/go-deploy/net/gotcp"
	"github.com/hellgate75/go-deploy/net/ssh"
)

func NewSshConnectionHandler() generic.ConnectionHandler {
	return ssh.NewSshConnectionHandler()
}

func NewGoTCPConnectionHandler() generic.ConnectionHandler {
	return gotcp.NewGoTCPConnectionHandler()
}
