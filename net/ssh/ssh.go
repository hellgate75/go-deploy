package ssh

import (
	"github.com/hellgate75/go-deploy/net/generic"
	"github.com/hellgate75/go-deploy/net/ssh/client"
)

func NewSshConnectionHandler(singleSession bool) generic.ConnectionHandler {
	if singleSession {
		return client.NewSingleSessionSSHConnection()
	}
	return client.NewSSHConnection()
}
