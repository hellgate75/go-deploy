package gotcp

import (
	"github.com/hellgate75/go-deploy/net/generic"
	"github.com/hellgate75/go-deploy/net/gotcp/client"
)

func NewGoTCPConnectionHandler() generic.ConnectionHandler {
	return client.NewGoTCPConnection()
}
