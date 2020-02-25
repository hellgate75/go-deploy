package ssh

import (
	"github.com/hellgate75/go-deploy/net/generic"
	"github.com/hellgate75/go-deploy/net/ssh/client"
)

func NewSshConnectionHandler() generic.ConnectionHandler {
	return client.NewSSHConnection()
}
