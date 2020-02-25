package ssh

import (
	"github.com/hellgate75/go-deploy/net/ssh/client"
)

func NewSshConnectionHandler() client.SSHConnectionHandler {
	return client.NewSSHConnection()
}
