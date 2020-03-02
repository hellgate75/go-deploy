package generic

import (
	"github.com/hellgate75/go-tcp-client/common"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
)

// Manages the remote file transfer
type FileTransfer interface {

	//MkDir: Create Remote folder
	MkDir(path string) error

	//MkDir: Create Remote folder
	MkDirAs(path string, mode os.FileMode) error

	//TransferFile: Tranfer single file
	TransferFileAs(path string, remotePath string, mode os.FileMode) error

	//TransferFolder: Tranfer folder recursively
	TransferFolderAs(path string, remotePath string, mode os.FileMode) error

	//TransferFile: Tranfer single file
	TransferFile(path string, remotePath string) error

	//TransferFolder: Tranfer folder recursively
	TransferFolder(path string, remotePath string) error

	// SetStdio: Sets the tranfer standard and error I/O streams
	SetStdio(stdout, stderr io.Writer) FileTransfer
}

// Manages the remote scripting execution
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

	// Close: Cloeses the reote connection
	Close() error

	// Start: start a remote shell on client
	Start() error

	// SetStdio: Sets the shell standard and error I/O streams
	SetStdio(stdin io.Reader, stdout, stderr io.Writer) RemoteShell
}

// Manages the client server connectivity
type NetworkClient interface {

	// Close: Cloeses the reote connection
	Close() error

	// Clone: Create clone of current Network Client, within all state
	Clone() NetworkClient

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

	// FileTranfer: Creates a file transfer manager session.
	FileTranfer() FileTransfer
}

// SSHConnectionHandler: Remote Connection handler and client connectivity maintainer
type ConnectionHandler interface {
	// GetClient: Retrieves current connected client or nil elsewise
	GetClient() NetworkClient

	// IsConnected: Gives state about connection success
	IsConnected() bool

	// Clone: Create clone of current Connection Handler, within all state
	Clone() ConnectionHandler

	// Close: Closes the remote connection
	Close() error

	// ConnectWithPasswd: Connect the SSH server with given passwd authmethod.
	ConnectWithPasswd(addr string, user string, passwd string) error

	// ConnectWithKey: Connect the SSH server with given SSH server with key authmethod.
	ConnectWithKey(addr string, user string, keyfile string) error

	// ConnectWithKeyAndPassphrase: Connect the SSH server with given key and a passphrase to decrypt the private key
	ConnectWithKeyAndPassphrase(addr string, user, keyfile string, passphrase string) error

	// Connect: Connect the SSH server using given network, address and configuration
	Connect(network, addr string, config *ssh.ClientConfig) error

	// Connect: Connect the PEM certificate and client key using given address and port
	ConnectWithCertificate(addr string, port string, certificate common.CertificateKeyPair) error
}
