package generic

import (
	"golang.org/x/crypto/ssh"
	"io"
)

// Manages the ssh remote scripting execution
type CommandsScript interface {

	// ExecuteWithOutput: Executes script(s) or command(s) using standard I/O
	ExecuteWithOutput() ([]byte, error)

	// ExecuteWithFullOutput: Executes script(s) or command(s) using standard and error I/O
	ExecuteWithFullOutput() ([]byte, error)

	// SetStdio: Sets the shell standard and error I/O streams
	SetStdio(stdout, stderr io.Writer) CommandsScript

	// NewCmd: Creates a client command
	NewCmd(cmd string) CommandsScript
}

// Configuration for the remote session terminal
type TerminalConfig struct {
	Term   string
	Height int
	Weight int
	Modes  ssh.TerminalModes
}

// Manages the interaction on client side with the Remote Shell
type RemoteShell interface {
	// Start: start a remote shell on client
	Start() error

	// SetStdio: Sets the shell standard and error I/O streams
	SetStdio(stdin io.Reader, stdout, stderr io.Writer) RemoteShell
}

// Manages the client server connectivity
type NetworkClient interface {

	// Close: Cloeses the reote connection
	Close() error

	// Terminal: Creates an interactive shell on client.
	Terminal(config *TerminalConfig) RemoteShell

	// NewCmd: Creates a client command
	NewCmd(cmd string) CommandsScript

	// Script: Creates a script client command list
	Script(script string) CommandsScript

	// ScriptFile: Creates a script client command list from file
	ScriptFile(fname string) CommandsScript

	// Shell: Creates a noninteractive shell on client.
	Shell() RemoteShell
}

// SSHConnectionHandler: Remote Connection handler and client connectivity maintainer
type ConnectionHandler interface {
	// GetClient: Retrieves current connected client or nil elsewise
	GetClient() NetworkClient

	// IsConnected: Gives state about connection success
	IsConnected() bool

	// Close: Closes the remote connection
	Close() error

	// ConnectWithPasswd: Connect the SSH server with given passwd authmethod.
	ConnectWithPasswd(addr string, user string, passwd string) error

	// ConnectWithKey: Connect the SSH server with given SSH server with key authmethod.
	ConnectWithKey(addr string, user string, keyfile string) error

	// ConnectWithKeyAndPassphrase: Connect the SSH server with given key and a passphrase to decrypt the private key
	ConnectWithKeyAndPassphrase(addr string, user, keyfile string, passphrase string) error

	// Connect: Connect the SSH server using given network, address and configuration
	Connect(network string, addr string, config *ssh.ClientConfig) error
}
