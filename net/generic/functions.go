package generic

import (
	"errors"
	"fmt"
	"github.com/hellgate75/go-deploy/log"
	"github.com/hellgate75/go-deploy/types/defaults"
	"github.com/hellgate75/go-deploy/types/module"
	"github.com/hellgate75/go-tcp-client/common"
	"os"
	"strings"
)

var Logger log.Logger = nil

func ConnectHandlerViaConfig(handler ConnectionHandler, host defaults.HostValue, netConfig *module.NetProtocolType, depConfig *module.DeployConfig) (NetworkClient, error) {
	var globalError error
	if string(module.RuntimeNetworkType.NetProtocol) == string(module.NET_PROTOCOL_SSH) {
		var missPassword bool = netConfig.Password == ""
		var missKeyFile bool = netConfig.KeyFile == ""
		var missKeyPassphrase bool = netConfig.Passphrase == ""
		var hostRef string = host.HostName
		if hostRef == "" {
			hostRef = host.IpAddress
		}
		addr := fmt.Sprintf("%s:%s", hostRef, host.Port)
		if !missPassword && missKeyFile {
			globalError = handler.ConnectWithPasswd(addr, netConfig.UserName, netConfig.Password)
		} else if !missPassword && !missKeyFile && !missKeyPassphrase {
			var keyFilePath string = netConfig.KeyFile
			if strings.Index(keyFilePath, ":") < 0 &&
				strings.Index(keyFilePath, "/") != 0 &&
				strings.Index(keyFilePath, "\\") != 0 {
				keyFilePath = depConfig.WorkDir + string(os.PathSeparator) + keyFilePath
			}
			globalError = handler.ConnectWithKeyAndPassphrase(addr, netConfig.UserName, keyFilePath, netConfig.Passphrase)
		} else if !missKeyFile {
			var keyFilePath string = netConfig.KeyFile
			if strings.Index(keyFilePath, ":") < 0 &&
				strings.Index(keyFilePath, "/") != 0 &&
				strings.Index(keyFilePath, "\\") != 0 {
				keyFilePath = depConfig.WorkDir + string(os.PathSeparator) + keyFilePath
			}
			globalError = handler.ConnectWithKey(addr, netConfig.UserName, keyFilePath)
		} else {
			Logger.Error("SSH: Unable to determine the Connection mode, between: password, key, key passphrase")
			return nil, errors.New("SSH: Unable to determine the connection mode")
		}
		if globalError != nil {
			return nil, globalError
		}
		return handler.GetClient(), nil

	} else if string(module.RuntimeNetworkType.NetProtocol) == string(module.NET_PROTOCOL_GO_DEPLOY_CLIENT) {
		var hostRef string = host.HostName
		if hostRef == "" {
			hostRef = host.IpAddress
		}
		var keyFilePath string = netConfig.KeyFile
		var certFilePath string = netConfig.Certificate
		if strings.Index(keyFilePath, ":") < 0 &&
			strings.Index(keyFilePath, "/") != 0 &&
			strings.Index(keyFilePath, "\\") != 0 {
			keyFilePath = depConfig.WorkDir + string(os.PathSeparator) + keyFilePath
		}
		if strings.Index(certFilePath, ":") < 0 &&
			strings.Index(certFilePath, "/") != 0 &&
			strings.Index(certFilePath, "\\") != 0 {
			certFilePath = depConfig.WorkDir + string(os.PathSeparator) + certFilePath
		}
		certificate := common.CertificateKeyPair{
			Key:  keyFilePath,
			Cert: certFilePath,
		}
		globalError = handler.ConnectWithCertificate(hostRef, host.Port, certificate)
		if globalError != nil {
			return nil, globalError
		}
		return handler.GetClient(), nil
	} else {
		Logger.Error("Unable to determine the Connection Handler type for: " + string(netConfig.NetProtocol))
		return nil, errors.New("Unable to determine the Connection Handler")
	}
	return nil, errors.New("Unable to connect to the network")
}
