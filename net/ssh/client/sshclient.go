// sshclient implements an ssh client
package client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

type sshScriptType byte
type sshShellType byte

const (
	cmdLine sshScriptType = iota
	rawScript
	scriptFile

	interactiveShell sshShellType = iota
	nonInteractiveShell
)

// Manages the ssh remote scripting execution
type SSHScript interface {

	// ExecuteWithOutput: Executes script(s) or command(s) using standard I/O
	ExecuteWithOutput() ([]byte, error)

	// ExecuteWithFullOutput: Executes script(s) or command(s) using standard and error I/O
	ExecuteWithFullOutput() ([]byte, error)

	// SetStdio: Sets the shell standard and error I/O streams
	SetStdio(stdout, stderr io.Writer) SSHScript

	// NewCmd: Creates a client command
	NewCmd(cmd string) SSHScript
}

type sshScript struct {
	client     *ssh.Client
	_type      sshScriptType
	script     *bytes.Buffer
	scriptFile string
	err        error

	stdout io.Writer
	stderr io.Writer
}

// Execute
func (rs *sshScript) execute() error {
	if rs.err != nil {
		fmt.Println(rs.err)
		return rs.err
	}

	if rs._type == cmdLine {
		return rs.runCmds()
	} else if rs._type == rawScript {
		return rs.runScript()
	} else if rs._type == scriptFile {
		return rs.runScriptFile()
	} else {
		return errors.New("SSHScript.execute: Not supported sshScript type")
	}
}

func (rs *sshScript) ExecuteWithOutput() ([]byte, error) {
	if rs.stdout != nil {
		return nil, errors.New("SSHScript.ExecuteWithOutput: Stdout already set")
	}
	var out bytes.Buffer
	rs.stdout = &out
	err := rs.execute()
	return out.Bytes(), err
}

func (rs *sshScript) ExecuteWithFullOutput() ([]byte, error) {
	if rs.stdout != nil {
		return nil, errors.New("SSHScript.ExecuteWithFullOutput: Stdout already set")
	}
	if rs.stderr != nil {
		return nil, errors.New("SSHScript.ExecuteWithFullOutput: Stderr already set")
	}

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)
	rs.stdout = &stdout
	rs.stderr = &stderr
	err := rs.execute()
	if err != nil {
		return stderr.Bytes(), err
	}
	return stdout.Bytes(), err
}

func (rs *sshScript) NewCmd(cmd string) SSHScript {
	_, err := rs.script.WriteString(cmd + "\n")
	if err != nil {
		rs.err = err
	}
	return rs
}

func (rs *sshScript) SetStdio(stdout, stderr io.Writer) SSHScript {
	rs.stdout = stdout
	rs.stderr = stderr
	return rs
}

func (rs *sshScript) runCmd(cmd string) error {
	session, err := rs.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	session.Stdout = rs.stdout
	session.Stderr = rs.stderr

	if err := session.Run(cmd); err != nil {
		return err
	}
	return nil
}

func (rs *sshScript) runCmds() error {
	for {
		statment, err := rs.script.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if err := rs.runCmd(statment); err != nil {
			return err
		}
	}

	return nil
}

func (rs *sshScript) runScript() error {
	session, err := rs.client.NewSession()
	if err != nil {
		return err
	}

	session.Stdin = rs.script
	session.Stdout = rs.stdout
	session.Stderr = rs.stderr

	if err := session.Shell(); err != nil {
		return err
	}
	if err := session.Wait(); err != nil {
		return err
	}

	return nil
}

func (rs *sshScript) runScriptFile() error {
	var buffer bytes.Buffer
	file, err := os.Open(rs.scriptFile)
	if err != nil {
		return err
	}
	_, err = io.Copy(&buffer, file)
	if err != nil {
		return err
	}

	rs.script = &buffer
	return rs.runScript()
}

// Configuration for the remote session terminal
type TerminalConfig struct {
	Term   string
	Height int
	Weight int
	Modes  ssh.TerminalModes
}

// Manages the interaction on client side with the remote SSH Shell
type SSHShell interface {
	// Start: start a remote shell on client
	Start() error

	// SetStdio: Sets the shell standard and error I/O streams
	SetStdio(stdin io.Reader, stdout, stderr io.Writer) SSHShell
}

type sshShell struct {
	client         *ssh.Client
	requestPty     bool
	terminalConfig *TerminalConfig

	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func (rs *sshShell) SetStdio(stdin io.Reader, stdout, stderr io.Writer) SSHShell {
	rs.stdin = stdin
	rs.stdout = stdout
	rs.stderr = stderr
	return rs
}

func (rs *sshShell) Start() error {
	session, err := rs.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	if rs.stdin == nil {
		session.Stdin = os.Stdin
	} else {
		session.Stdin = rs.stdin
	}
	if rs.stdout == nil {
		session.Stdout = os.Stdout
	} else {
		session.Stdout = rs.stdout
	}
	if rs.stderr == nil {
		session.Stderr = os.Stderr
	} else {
		session.Stderr = rs.stderr
	}

	if rs.requestPty {
		tc := rs.terminalConfig
		if tc == nil {
			tc = &TerminalConfig{
				Term:   "xterm",
				Height: 40,
				Weight: 80,
			}
		}
		if err := session.RequestPty(tc.Term, tc.Height, tc.Weight, tc.Modes); err != nil {
			return err
		}
	}

	if err := session.Shell(); err != nil {
		return err
	}

	if err := session.Wait(); err != nil {
		return err
	}

	return nil
}

// Manages the client server connectivity
type SSHClient interface {

	// Close: Cloeses the reote connection
	Close() error

	// Terminal: Creates an interactive shell on client.
	Terminal(config *TerminalConfig) SSHShell

	// NewCmd: Creates a client command
	NewCmd(cmd string) SSHScript

	// Script: Creates a script client command list
	Script(script string) SSHScript

	// ScriptFile: Creates a script client command list from file
	ScriptFile(fname string) SSHScript

	// Shell: Creates a noninteractive shell on client.
	Shell() SSHShell
}

type sshClient struct {
	client *ssh.Client
}

func (c *sshClient) Close() error {
	return c.client.Close()
}

func (c *sshClient) Terminal(config *TerminalConfig) SSHShell {
	return &sshShell{
		client:         c.client,
		terminalConfig: config,
		requestPty:     true,
	}
}

func (c *sshClient) NewCmd(cmd string) SSHScript {
	return &sshScript{
		_type:  cmdLine,
		client: c.client,
		script: bytes.NewBufferString(cmd + "\n"),
	}
}

func (c *sshClient) Script(script string) SSHScript {
	return &sshScript{
		_type:  rawScript,
		client: c.client,
		script: bytes.NewBufferString(script + "\n"),
	}
}

func (c *sshClient) ScriptFile(fname string) SSHScript {
	return &sshScript{
		_type:      scriptFile,
		client:     c.client,
		scriptFile: fname,
	}
}

func (c *sshClient) Shell() SSHShell {
	return &sshShell{
		client:     c.client,
		requestPty: false,
	}
}

// SSHConnectionHandler: SSH Connection handler and client connectivity maintainer
type SSHConnectionHandler interface {
	// GetClient: Retrieves current connected client or nil elsewise
	GetClient() SSHClient

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

type sshConnection struct {
	_client SSHClient
}

func (conn *sshConnection) GetClient() SSHClient {
	return conn._client
}

func (conn *sshConnection) IsConnected() bool {
	return conn._client != nil
}

func (conn *sshConnection) Close() error {
	if !conn.IsConnected() {
		return errors.New("SSHConnectionHandler.Close: Not connected!!")
	}
	err := conn._client.Close()
	if err != nil {
		return errors.New("SSHConnectionHandler.Close: " + err.Error())
	}
	return nil
}

func (conn *sshConnection) ConnectWithPasswd(addr string, user string, passwd string) error {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(passwd),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}

	return conn.Connect("tcp", addr, config)
}

func (conn *sshConnection) ConnectWithKey(addr string, user string, keyfile string) error {
	key, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return errors.New("SSHConnectionHandler.ConnectWithKey: " + err.Error())
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return errors.New("SSHConnectionHandler.ConnectWithKey: " + err.Error())
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}

	return conn.Connect("tcp", addr, config)
}

func (conn *sshConnection) ConnectWithKeyAndPassphrase(addr string, user, keyfile string, passphrase string) error {
	key, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return errors.New("SSHConnectionHandler.ConnectWithKeyAndPassphrase: " + err.Error())
	}

	signer, err := ssh.ParsePrivateKeyWithPassphrase(key, []byte(passphrase))
	if err != nil {
		return errors.New("SSHConnectionHandler.ConnectWithKeyAndPassphrase: " + err.Error())
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}

	return conn.Connect("tcp", addr, config)
}

func (conn *sshConnection) Connect(network string, addr string, config *ssh.ClientConfig) error {
	client, err := ssh.Dial(network, addr, config)
	if err != nil {
		return errors.New("SSHConnectionHandler.Connect: " + err.Error())
	}
	conn._client = &sshClient{
		client: client,
	}
	return nil
}

// NewSSHConnection: Creates a new SSH connection handler
func NewSSHConnection() SSHConnectionHandler {
	return &sshConnection{
		_client: nil,
	}
}
