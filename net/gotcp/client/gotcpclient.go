// sshclient implements an ssh client
package client

import (
	"bytes"
	"errors"
	"fmt"
	depio "github.com/hellgate75/go-deploy/io"
	"github.com/hellgate75/go-deploy/net/generic"
	"github.com/hellgate75/go-tcp-server/client/worker"
	"github.com/hellgate75/go-tcp-server/common"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"os"
	"time"
)

type goTcpScriptType byte
type goTcpShellType byte

const (
	cmdLine goTcpScriptType = iota
	rawScript
	scriptFile

	interactiveShell goTcpShellType = iota
	nonInteractiveShell
)

type goTcpTranfer struct {
	client common.TCPClient
	stdout io.Writer
	stderr io.Writer
}

func (ts *goTcpTranfer) SetStdio(stdout, stderr io.Writer) generic.FileTransfer {
	ts.stdout = stdout
	ts.stderr = stderr
	return ts
}

func (ts *goTcpTranfer) MkDir(path string) error {
	return ts.MkDirAs(path, 0644)
}

func (ts *goTcpTranfer) MkDirAs(path string, mode os.FileMode) error {
	err := mkDir(path, ts.client, mode, ts.stdout, ts.stderr)
	return err
}

func (ts *goTcpTranfer) TransferFile(path string, remotePath string) error {
	return ts.TransferFileAs(path, remotePath, 0644)
}

func (ts *goTcpTranfer) TransferFileAs(path string, remotePath string, mode os.FileMode) error {
	stat, errS := os.Stat(path)
	if errS != nil {
		return errors.New("GoTcpTransfer.TransferFileAs::StatFile: " + errS.Error())
	}
	if stat.IsDir() {
		return ts.TransferFolderAs(path, remotePath, mode)
	}
	err := copyFile(path, remotePath, ts.client, mode, ts.stdout, ts.stderr)
	if err != nil {
		return err
	}
	return nil
}

func (ts *goTcpTranfer) TransferFolder(path string, remotePath string) error {
	return ts.TransferFolderAs(path, remotePath, 0644)
}

func (ts *goTcpTranfer) TransferFolderAs(path string, remotePath string, mode os.FileMode) error {
	stat, errS := os.Stat(path)
	if errS != nil {
		return errors.New("GoTcpTransfer.TransferFolder::StatFile: " + errS.Error())
	}
	if !stat.IsDir() {
		return ts.TransferFileAs(path, remotePath, mode)
	}
	err := executeFunc(path, remotePath, ts.client, mode, ts.stdout, ts.stderr)
	if err != nil {
		return err
	}
	return nil
}

func executeFunc(path string, remotePath string, client common.TCPClient, mode os.FileMode, stdout io.Writer, stderr io.Writer) error {
	stat, errS := os.Stat(path)
	if errS != nil {
		return errS
	}
	if stat.IsDir() {
		mkDir(remotePath, client, mode, stdout, stderr)
		files, err := ioutil.ReadDir(path)
		if err != nil {
			return err
		}
		for _, f := range files {
			var fName = path + depio.GetPathSeparator() + f.Name()
			var fRemoteName = remotePath + "/" + f.Name()
			err := executeFunc(fName, fRemoteName, client, f.Mode(), stdout, stderr)
			if err != nil {
				return err
			}
		}
	} else {
		err := copyFile(path, remotePath, client, mode, stdout, stderr)
		if err != nil {
			return err
		}
	}
	return nil
}

func mkDir(remotePath string, client common.TCPClient, mode os.FileMode, stdout io.Writer, stderr io.Writer) error {
	defer func() {
		client.SendText("exit")
		client.Close()
	}()
	err := client.Open(false)
	if err != nil {
		return err
	}
	err = client.ApplyCommand("transfer-file", "folder", remotePath, mode.String())
	if err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	resp, errR := client.ReadAnswer()
	if errR != nil {
		return errR
	}
	if len(resp) >= 2 {
		if resp[0:2] == "ko" {
			if len(resp) > 3 {
				if stderr != nil {
					stderr.Write([]byte(resp[3:]))
				}
				return errors.New("TCPScript.runCmd: " + resp[3:])
			} else {
				if stderr != nil {
					stderr.Write([]byte(resp))
				}
				return errors.New("TCPScript.runCmd: Undefined error from server")
			}
		} else if resp[0:2] == "ok" {
			if stdout != nil {
				stdout.Write([]byte(resp))
			}
		} else {
			if len(resp) > 3 {
				if stdout != nil {
					stdout.Write([]byte(resp[3:]))
				}
			} else {
				if stdout != nil {
					stdout.Write([]byte(resp))
				}
			}

		}
	} else {
		if stdout != nil {
			stdout.Write([]byte(resp))
		}
	}
	return nil
}

func copyFile(localPath string, remotePath string, client common.TCPClient, mode os.FileMode, stdout io.Writer, stderr io.Writer) error {
	defer func() {
		client.SendText("exit")
		client.Close()
	}()
	err := client.Open(false)
	if err != nil {
		return err
	}
	err = client.ApplyCommand("transfer-file", localPath, remotePath, mode.String())
	if err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	resp, errR := client.ReadAnswer()
	if errR != nil {
		return errR
	}
	if len(resp) >= 2 {
		if resp[0:2] == "ko" {
			if len(resp) > 3 {
				if stderr != nil {
					stderr.Write([]byte(resp[3:]))
				}
				return errors.New("TCPScript.runCmd: " + resp[3:])
			} else {
				if stderr != nil {
					stderr.Write([]byte(resp))
				}
				return errors.New("TCPScript.runCmd: Undefined error from server")
			}
		} else if resp[0:2] == "ok" {
			if stdout != nil {
				stdout.Write([]byte(resp))
			}
		} else {
			if len(resp) > 3 {
				if stdout != nil {
					stdout.Write([]byte(resp[3:]))
				}
			} else {
				if stdout != nil {
					stdout.Write([]byte(resp))
				}
			}

		}
	} else {
		if stdout != nil {
			stdout.Write([]byte(resp))
		}
	}
	return nil
}

type goTcpScript struct {
	client     common.TCPClient
	_type      goTcpScriptType
	script     *bytes.Buffer
	scriptFile string
	err        error

	stdout io.Writer
	stderr io.Writer
}

// Execute
func (rs *goTcpScript) execute() error {
	if rs.err != nil {
		return errors.New("GoTCPScript.execute: " + rs.err.Error())
	}
	if rs._type == cmdLine {
		return rs.runCmds()
	} else if rs._type == rawScript {
		return rs.runScript()
	} else if rs._type == scriptFile {
		return rs.runScriptFile()
	} else {
		return errors.New(fmt.Sprintf("GoTCPScript.execute: Not supported execution type: %v", rs._type))
	}
}

func (rs *goTcpScript) ExecuteWithOutput() ([]byte, error) {
	if rs.stdout != nil {
		return nil, errors.New("GoTCPScript.ExecuteWithOutput: Stdout already set")
	}
	var out bytes.Buffer
	rs.stdout = &out
	err := rs.execute()
	if err != nil {
		err = errors.New("GoTCPScript.ExecuteWithFullOutput: " + err.Error())
	}
	return out.Bytes(), err
}

func (rs *goTcpScript) ExecuteWithFullOutput() ([]byte, error) {
	if rs.stdout != nil {
		return nil, errors.New("GoTCPScript.ExecuteWithFullOutput: Stdout already set")
	}
	if rs.stderr != nil {
		return nil, errors.New("GoTCPScript.ExecuteWithFullOutput: Stderr already set")
	}

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)
	rs.stdout = &stdout
	rs.stderr = &stderr
	err := rs.execute()
	if err != nil {
		return stderr.Bytes(), errors.New("GoTCPScript.ExecuteWithFullOutput: " + err.Error())
	}
	return stdout.Bytes(), err
}

func (rs *goTcpScript) NewCmd(cmd string) generic.CommandsScript {
	_, err := rs.script.WriteString(cmd + "\n")
	if err != nil {
		rs.err = errors.New("GoTCPScript.NewCmd: " + err.Error())
	}
	return rs
}

func (rs *goTcpScript) SetStdio(stdout, stderr io.Writer) generic.CommandsScript {
	rs.stdout = stdout
	rs.stderr = stderr
	return rs
}

func (rs *goTcpScript) runCmd(cmd string) error {
	defer func() {
		rs.client.Close()
	}()
	err := rs.client.Open(false)
	if err != nil {
		return err
	}
	err = rs.client.ApplyCommand("shell", "false", cmd)
	if err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	resp, errR := rs.client.ReadAnswer()
	if errR != nil {
		return errR
	}
	if len(resp) >= 2 {
		if resp[0:2] == "ko" {
			if len(resp) > 3 {
				if rs.stderr != nil {
					rs.stderr.Write([]byte(resp[3:]))
				}
				return errors.New("GoTCPScript.runCmd: " + resp[3:])
			} else {
				if rs.stderr != nil {
					rs.stderr.Write([]byte(resp))
				}
				return errors.New("GoTCPScript.runCmd: Undefined error from server")
			}
		} else if resp[0:2] == "ok" {
			if rs.stdout != nil {
				rs.stdout.Write([]byte(resp))
			}
		} else {
			if len(resp) > 3 {
				if rs.stdout != nil {
					rs.stdout.Write([]byte(resp[3:]))
				}
			} else {
				if rs.stdout != nil {
					rs.stdout.Write([]byte(resp))
				}
			}

		}
	} else {
		if rs.stdout != nil {
			rs.stdout.Write([]byte(resp))
		}
	}
	return nil
}

func (rs *goTcpScript) runCmds() error {
	for {
		statment, err := rs.script.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.New("GoTCPScript.runCmds: " + err.Error())
		}

		if err := rs.runCmd(statment); err != nil {
			return errors.New("GoTCPScript.runCmds: " + err.Error())
		}
	}

	return nil
}

func (rs *goTcpScript) runScript() error {
	defer func() {
		rs.client.Close()
	}()
	err := rs.client.Open(false)
	if err != nil {
		return err
	}
	err = rs.client.ApplyCommand("shell", "false", string(rs.script.Bytes()))
	if err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	resp, errR := rs.client.ReadAnswer()
	if errR != nil {
		return errR
	}
	if len(resp) >= 2 {
		if resp[0:2] == "ko" {
			if len(resp) > 3 {
				if rs.stderr != nil {
					rs.stderr.Write([]byte(resp[3:]))
				}
				return errors.New("GoTCPScript.runScript: " + resp[3:])
			} else {
				if rs.stderr != nil {
					rs.stderr.Write([]byte(resp))
				}
				return errors.New("GoTCPScript.runScript: Undefined error from server")
			}
		} else if resp[0:2] == "ok" {
			if rs.stdout != nil {
				rs.stdout.Write([]byte(resp))
			}
		} else {
			if len(resp) > 3 {
				if rs.stdout != nil {
					rs.stdout.Write([]byte(resp[3:]))
				}
			} else {
				if rs.stdout != nil {
					rs.stdout.Write([]byte(resp))
				}
			}

		}
	} else {
		if rs.stdout != nil {
			rs.stdout.Write([]byte(resp))
		}
	}
	return nil
}

func (rs *goTcpScript) runScriptFile() error {
	defer func() {
		rs.client.Close()
	}()
	err := rs.client.Open(false)
	if err != nil {
		return err
	}
	err = rs.client.ApplyCommand("shell", "false", rs.scriptFile)
	if err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	resp, errR := rs.client.ReadAnswer()
	if errR != nil {
		return errR
	}
	if len(resp) >= 2 {
		if resp[0:2] == "ko" {
			if len(resp) > 3 {
				if rs.stderr != nil {
					rs.stderr.Write([]byte(resp[3:]))
				}
				return errors.New("GoTCPScript.runScript: " + resp[3:])
			} else {
				if rs.stderr != nil {
					rs.stderr.Write([]byte(resp))
				}
				return errors.New("GoTCPScript.runScript: Undefined error from server")
			}
		} else if resp[0:2] == "ok" {
			if rs.stdout != nil {
				rs.stdout.Write([]byte(resp))
			}
		} else {
			if len(resp) > 3 {
				if rs.stdout != nil {
					rs.stdout.Write([]byte(resp[3:]))
				}
			} else {
				if rs.stdout != nil {
					rs.stdout.Write([]byte(resp))
				}
			}

		}
	} else {
		if rs.stdout != nil {
			rs.stdout.Write([]byte(resp))
		}
	}
	return nil
}

type goTcpShell struct {
	client         common.TCPClient
	requestPty     bool
	terminalConfig *generic.TerminalConfig

	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func (rs *goTcpShell) SetStdio(stdin io.Reader, stdout, stderr io.Writer) generic.RemoteShell {
	rs.stdin = stdin
	rs.stdout = stdout
	rs.stderr = stderr
	return rs
}

func (rs *goTcpShell) Close() error {
	if rs.stdin == nil || rs.stdout == nil || rs.stderr == nil {
		return errors.New("GoTcpShell:Close() -> Shell not open for miss of stdout, stderr, stdin")
	}
	rs.stdout.Write([]byte("exit\n"))

	return nil
}

func (rs *goTcpShell) Start() error {
	if rs.stdin == nil || rs.stdout == nil || rs.stderr == nil {
		return errors.New("GoTcpShell:Start() -> Please provide stdout, stderr, stdin for the shell execution")
	}
	err := rs.client.ApplyCommand("shell", "true", "", rs.stdin, rs.stdout, rs.stderr)
	if err != nil {
		return errors.New("GoTcpShell:Start() -> Details: " + err.Error())
	}
	return nil
}

type goTcpClient struct {
	client common.TCPClient
}

func (c *goTcpClient) Close() error {
	return c.client.Close()
}

func (c *goTcpClient) Terminal(config *generic.TerminalConfig) generic.RemoteShell {
	return &goTcpShell{
		client:         c.client.Clone(),
		terminalConfig: config,
		requestPty:     true,
	}
}

func (c *goTcpClient) NewCmd(cmd string) generic.CommandsScript {
	return &goTcpScript{
		_type:  cmdLine,
		client: c.client.Clone(),
		script: bytes.NewBufferString(cmd + "\n"),
	}
}

func (c *goTcpClient) Script(script string) generic.CommandsScript {
	return &goTcpScript{
		_type:  rawScript,
		client: c.client.Clone(),
		script: bytes.NewBufferString(script + "\n"),
	}
}

func (c *goTcpClient) ScriptFile(fname string) generic.CommandsScript {
	return &goTcpScript{
		_type:      scriptFile,
		client:     c.client.Clone(),
		scriptFile: fname,
	}
}

func (c *goTcpClient) FileTranfer() generic.FileTransfer {
	return &goTcpTranfer{
		client: c.client.Clone(),
	}
}

func (c *goTcpClient) Shell() generic.RemoteShell {
	return &goTcpShell{
		client:     c.client.Clone(),
		requestPty: false,
	}
}

type goTcpConnection struct {
	_client generic.NetworkClient
}

func (conn *goTcpConnection) GetClient() generic.NetworkClient {
	return conn._client
}

func (conn *goTcpConnection) IsConnected() bool {
	return conn._client != nil
}

func (conn *goTcpConnection) Close() error {
	if !conn.IsConnected() {
		return errors.New("SSHConnectionHandler.Close: Not connected!!")
	}
	err := conn._client.Close()
	if err != nil {
		return errors.New("SSHConnectionHandler.Close: " + err.Error())
	}
	return nil
}

func (conn *goTcpConnection) ConnectWithPasswd(addr string, user string, passwd string) error {
	return errors.New("User/password connection not allowed to Go TCP Server")
}

func (conn *goTcpConnection) ConnectWithKey(addr string, user string, keyfile string) error {
	return errors.New("User/rsa key connection not allowed to Go TCP Server")
}

func (conn *goTcpConnection) ConnectWithKeyAndPassphrase(addr string, user, keyfile string, passphrase string) error {
	return errors.New("User/rsa key connection not allowed to Go TCP Server")
}

func (conn *goTcpConnection) Connect(network, addr string, config *ssh.ClientConfig) error {
	return errors.New("User/rsa key connection not allowed to Go TCP Server")
}

func (conn *goTcpConnection) ConnectWithCertificate(addr string, port string, certificate common.CertificateKeyPair) error {
	conn._client = &goTcpClient{
		client: worker.NewClient(certificate, addr, port),
	}
	return nil
}

// NewSSHConnection: Creates a new SSH connection handler
func NewGoTCPConnection() generic.ConnectionHandler {
	return &goTcpConnection{
		_client: nil,
	}
}
